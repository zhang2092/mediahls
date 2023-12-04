package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/zhang2092/mediahls/internal/db"
)

// obj

// homePageData 首页数据
type homePageData struct {
	Authorize
	Videos []db.Video
}

// view

// home 首页
func (server *Server) homeView(w http.ResponseWriter, r *http.Request) {
	data := homePageData{}
	auth, err := server.withCookie(r)
	if err == nil {
		data.Authorize = *auth
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
			data.Videos = append(data.Videos, item)
		}
	}

	renderLayout(w, r, data, "web/templates/home.html.tmpl")
}
