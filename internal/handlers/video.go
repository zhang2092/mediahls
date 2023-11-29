package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type playData struct {
	AuthID   string
	AuthName string
	Url      string
}

func (server *Server) play(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	xid := vars["xid"]
	data := playData{
		Url: "/media/" + xid + "/stream/",
	}
	auth, err := server.withCookie(r)
	if err == nil {
		data.AuthID = auth.AuthID
		data.AuthName = auth.AuthName
	}
	render(w, data, "web/templates/video/play.html.tmpl", "web/templates/base/header.html.tmpl", "web/templates/base/footer.html.tmpl")
}

/*
// 直接播放mp4
video, err := os.Open("web/statics/git.mp4")
if err != nil {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("failed to open file"))
	return
}
defer video.Close()

http.ServeContent(w, r, "git", time.Now(), video)
*/

func (server *Server) stream(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	mId := vars["xid"]
	fmt.Println(mId)

	segName, ok := vars["segName"]

	if !ok {
		mediaBase := getMediaBase(mId)
		m3u8Name := "index.m3u8"
		serveHlsM3u8(response, request, mediaBase, m3u8Name)
	} else {
		mediaBase := getMediaBase(mId)
		serveHlsTs(response, request, mediaBase, segName)
	}

}
func getMediaBase(mId string) string {
	mediaRoot := "media"
	return fmt.Sprintf("%s/%s", mediaRoot, mId)
}

func serveHlsM3u8(w http.ResponseWriter, r *http.Request, mediaBase, m3u8Name string) {
	fmt.Println("serveHlsM3u8...")

	mediaFile := fmt.Sprintf("%s/%s", mediaBase, m3u8Name)
	fmt.Println(mediaFile)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "application/x-mpegURL")

}

func serveHlsTs(w http.ResponseWriter, r *http.Request, mediaBase, segName string) {
	fmt.Println("serveHlsTs...")

	mediaFile := fmt.Sprintf("%s/%s", mediaBase, segName)
	fmt.Println(mediaFile)

	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "video/MP2T")
}
