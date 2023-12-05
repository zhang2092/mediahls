package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/zhang2092/mediahls/internal/db"
)

// view

// home 首页
func (server *Server) homeView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var result []db.Video
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
			result = append(result, item)
		}
	}

	server.renderLayout(w, r, result, "home.html.tmpl")
}
