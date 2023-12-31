package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/csrf"
	hds "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/zhang2092/mediahls/internal/db"
	"github.com/zhang2092/mediahls/internal/pkg/config"
	"github.com/zhang2092/mediahls/internal/pkg/logger"
	"github.com/zhang2092/mediahls/internal/pkg/token"
	"github.com/zhang2092/mediahls/internal/worker"
)

type Server struct {
	templateFS fs.FS
	staticFS   fs.FS

	conf         *config.Config
	router       *mux.Router
	secureCookie *securecookie.SecureCookie

	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(templateFS fs.FS, staticFS fs.FS, conf *config.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(conf.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	hashKey := securecookie.GenerateRandomKey(32)
	blockKey := securecookie.GenerateRandomKey(32)
	secureCookie := securecookie.New(hashKey, blockKey)
	// secureCookie.MaxAge(7200)

	server := &Server{
		templateFS:      templateFS,
		staticFS:        staticFS,
		conf:            conf,
		secureCookie:    secureCookie,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.PathPrefix("/statics/").Handler(http.StripPrefix("/statics/", http.FileServer(http.FS(server.staticFS))))
	router.PathPrefix("/upload/imgs").Handler(http.StripPrefix("/upload/imgs/", http.FileServer(http.Dir("./upload/imgs"))))

	csrfMiddleware := csrf.Protect(
		[]byte(securecookie.GenerateRandomKey(32)),
		csrf.Secure(false),
		csrf.HttpOnly(true),
		csrf.FieldName("csrf_token"),
		csrf.CookieName("authorize_csrf"),
	)
	router.Use(csrfMiddleware)
	router.Use(server.setUser)

	router.Handle("/register", hds.MethodHandler{
		http.MethodGet:  http.HandlerFunc(server.registerView),
		http.MethodPost: http.HandlerFunc(server.register),
	})
	router.Handle("/login", hds.MethodHandler{
		http.MethodGet:  http.HandlerFunc(server.loginView),
		http.MethodPost: http.HandlerFunc(server.login),
	})
	router.HandleFunc("/logout", server.logout).Methods(http.MethodGet)

	router.HandleFunc("/", server.homeView).Methods(http.MethodGet)

	router.HandleFunc("/play/{xid}", server.videoView).Methods(http.MethodGet)
	router.HandleFunc("/media/{xid}/stream/", server.stream).Methods(http.MethodGet)
	router.HandleFunc("/media/{xid}/stream/{segName:[a-z0-9]+.ts}", server.stream).Methods(http.MethodGet)

	subRouter := router.PathPrefix("/").Subrouter()
	subRouter.Use(server.authorize)

	subRouter.HandleFunc("/me/videos", server.videosView).Methods(http.MethodGet)
	subRouter.HandleFunc("/me/videos/p{page}", server.videosView).Methods(http.MethodGet)
	subRouter.HandleFunc("/me/videos/update", server.editVideoView).Methods(http.MethodGet)
	subRouter.HandleFunc("/me/videos/update/{xid}", server.editVideoView).Methods(http.MethodGet)
	subRouter.HandleFunc("/me/videos/update", server.editVideo).Methods(http.MethodPost)
	subRouter.HandleFunc("/me/videos/delete", server.deleteVideo).Methods(http.MethodPost)

	subRouter.HandleFunc("/upload_image", server.uploadImage).Methods(http.MethodPost)
	subRouter.HandleFunc("/upload_file", server.uploadVideo).Methods(http.MethodPost)

	subRouter.HandleFunc("/transfer/{xid}", server.transfer).Methods(http.MethodPost)

	server.router = router
}

func (server *Server) Start(db *sql.DB) {
	srv := &http.Server{
		Addr:    server.conf.ServerAddress,
		Handler: hds.CompressHandler(server.router),
	}

	go func() {
		log.Printf("server start on: %s\n", server.conf.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Close(); err != nil {
		log.Fatal("Server db to shutdown:", err)
	}
	if err := logger.Logger.Sync(); err != nil {
		log.Fatal("Server log sync:", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
