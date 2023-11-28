package handlers

import "net/http"

func (server *Server) play(w http.ResponseWriter, r *http.Request) {

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
