package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/zhang2092/mediahls/internal/db"
)

type pageData struct {
	Authorize
	Videos []db.Video
}

func (server *Server) home(w http.ResponseWriter, r *http.Request) {
	pd := pageData{}
	auth, err := server.withCookie(r)
	if err == nil {
		pd.Authorize = *auth
	}

	ctx := r.Context()
	videos, err := server.store.ListVideos(ctx, db.ListVideosParams{
		Limit:  100,
		Offset: 0,
	})
	if err == nil {
		for _, item := range videos {
			if len(item.Description) > 65 {
				temp := strings.TrimSpace(item.Description[0:65]) + "..."
				item.Description = temp
				log.Println(item.Description)
			}
			pd.Videos = append(pd.Videos, item)
		}
	}
	renderHome(w, pd)
}

func renderHome(w http.ResponseWriter, data any) {
	renderLayout(w, data, "web/templates/home.html.tmpl")
}
