package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/zhang2092/mediahls/internal/db"
	"github.com/zhang2092/mediahls/internal/pkg/config"
	"github.com/zhang2092/mediahls/internal/pkg/logger"
	"github.com/zhang2092/mediahls/internal/pkg/token"
)

type Server struct {
	conf   *config.Config
	router *mux.Router

	store      db.Store
	tokenMaker token.Maker
}

func NewServer(conf *config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(conf.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		conf:       conf,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := mux.NewRouter()
	router.PathPrefix("/statics/").Handler(http.StripPrefix("/statics/", http.FileServer(http.Dir("web/statics"))))

	router.HandleFunc("/register", server.registerView).Methods(http.MethodGet)
	router.HandleFunc("/register", server.register).Methods(http.MethodPost)
	router.HandleFunc("/login", server.loginView).Methods(http.MethodGet)
	router.HandleFunc("/login", server.login).Methods(http.MethodPost)

	router.HandleFunc("/", server.home).Methods(http.MethodGet)

	router.HandleFunc("/play/{xid}", server.play).Methods(http.MethodGet)

	router.HandleFunc("/upload", server.upload).Methods(http.MethodGet, http.MethodPost)

	server.router = router
}

func (server *Server) Start(db *sql.DB) {
	srv := &http.Server{
		Addr:    server.conf.ServerAddress,
		Handler: server.router,
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
