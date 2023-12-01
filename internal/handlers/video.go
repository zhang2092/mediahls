package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/zhang2092/mediahls/internal/db"
	"github.com/zhang2092/mediahls/internal/pkg/convert"
	"github.com/zhang2092/mediahls/internal/pkg/fileutil"
	"github.com/zhang2092/mediahls/internal/pkg/logger"
)

type playData struct {
	Authorize
	Url   string
	Video db.Video
}

func (server *Server) play(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	xid := vars["xid"]
	video, _ := server.store.GetVideo(r.Context(), xid)
	data := playData{
		Video: video,
	}
	auth, err := server.withCookie(r)
	if err == nil {
		data.Authorize = *auth
	}
	renderLayout(w, data, "web/templates/video/play.html.tmpl")
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
	mediaFile := fmt.Sprintf("%s/%s", mediaBase, m3u8Name)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "application/x-mpegURL")

}

func serveHlsTs(w http.ResponseWriter, r *http.Request, mediaBase, segName string) {
	mediaFile := fmt.Sprintf("%s/%s", mediaBase, segName)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "video/MP2T")
}

type meVideoData struct {
	Authorize
	Videos []db.Video
}

func (server *Server) videosView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := meVideoData{
		Authorize: withUser(ctx),
	}

	vars := mux.Vars(r)
	page, err := strconv.Atoi(vars["page"])
	if err != nil {
		page = 1
	}
	videos, err := server.store.ListVideosByUser(ctx, db.ListVideosByUserParams{
		UserID: data.Authorize.ID,
		Limit:  16,
		Offset: int32((page - 1) * 16),
	})
	if err == nil {
		for _, item := range videos {
			if len(item.Description) > 65 {
				temp := strings.TrimSpace(item.Description[0:65]) + "..."
				item.Description = temp
			}
			data.Videos = append(data.Videos, item)
		}
	}

	renderLayout(w, data, "web/templates/me/videos.html.tmpl")
}

func (server *Server) createVideoView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	xid := vars["xid"]
	vm := videoCreateResp{
		Authorize: withUser(r.Context()),
	}
	if len(xid) > 0 {
		if v, err := server.store.GetVideo(r.Context(), xid); err == nil {
			vm.ID = v.ID
			vm.Title = v.Title
			vm.Images = v.Images
			vm.Description = v.Description
			vm.OriginLink = v.OriginLink
			vm.Status = int(v.Status)
		}
	}
	renderCreateVideo(w, vm)
}

func renderCreateVideo(w http.ResponseWriter, data any) {
	renderLayout(w, data, "web/templates/video/edit.html.tmpl")
}

type videoCreateResp struct {
	Authorize
	Summary        string
	ID             string
	IDErr          string
	Title          string
	TitleErr       string
	Images         string
	ImagesErr      string
	Description    string
	DescriptionErr string
	OriginLink     string
	OriginLinkErr  string
	Status         int
	StatusErr      string
}

func viladatorCreateVedio(r *http.Request) (*videoCreateResp, bool) {
	ok := true
	status, _ := strconv.Atoi(r.PostFormValue("status"))
	errs := &videoCreateResp{
		Authorize:   withUser(r.Context()),
		ID:          r.PostFormValue("id"),
		Title:       r.PostFormValue("title"),
		Images:      r.PostFormValue("images"),
		Description: r.PostFormValue("description"),
		OriginLink:  r.PostFormValue("origin_link"),
		Status:      status,
	}

	if len(errs.Title) == 0 {
		errs.TitleErr = "请填写正确的标题"
		ok = false
	}

	exist, _ := fileutil.PathExists(strings.TrimPrefix(errs.Images, "/"))
	if !exist {
		errs.ImagesErr = "请先上传图片"
		ok = false
	}

	if len(errs.Description) == 0 {
		errs.DescriptionErr = "请填写描述"
		ok = false
	}

	exist, _ = fileutil.PathExists(strings.TrimPrefix(errs.OriginLink, "/"))
	if !exist {
		errs.OriginLinkErr = "请先上传视频"
		ok = false
	}

	return errs, ok
}

func (server *Server) createVideo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		renderCreateVideo(w, videoCreateResp{Summary: "请求网络错误, 请刷新重试"})
		return
	}

	vm, ok := viladatorCreateVedio(r)
	if !ok {
		renderCreateVideo(w, vm)
		return
	}

	curTime := time.Now()
	ctx := r.Context()
	u := withUser(ctx)
	if len(vm.ID) == 0 {
		_, err := server.store.CreateVideo(ctx, db.CreateVideoParams{
			ID:          genId(),
			Title:       vm.Title,
			Description: vm.Description,
			Images:      vm.Images,
			OriginLink:  vm.OriginLink,
			PlayLink:    "",
			UserID:      u.ID,
			CreateBy:    u.Name,
		})
		if err != nil {
			vm.Summary = "添加视频失败"
			renderCreateVideo(w, vm)
			return
		}
	} else {
		_, err := server.store.UpdateVideo(ctx, db.UpdateVideoParams{
			ID:          vm.ID,
			Title:       vm.Title,
			Description: vm.Description,
			Images:      vm.Images,
			Status:      int32(vm.Status),
			UpdateAt:    curTime,
			UpdateBy:    u.Name,
		})
		if err != nil {
			vm.Summary = "更新视频失败"
			renderCreateVideo(w, vm)
			return
		}
	}

	http.Redirect(w, r, "/me/videos", http.StatusFound)
}

type transferData struct {
	Authorize
	Video db.Video
}

func (server *Server) transferView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	xid := vars["xid"]
	v, _ := server.store.GetVideo(r.Context(), xid)
	data := transferData{
		Video: v,
	}
	u, err := server.withCookie(r)
	if err == nil {
		data.Authorize = *u
	}
	renderLayout(w, data, "web/templates/video/transfer.html.tmpl")
}

func (server *Server) transfer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	xid := vars["xid"]
	ctx := r.Context()
	v, err := server.store.GetVideo(ctx, xid)
	if err != nil {
		http.Error(w, "视频信息错误", http.StatusInternalServerError)
		return
	}

	u := withUser(ctx)
	v, err = server.store.UpdateVideoStatus(ctx, db.UpdateVideoStatusParams{
		ID:       v.ID,
		Status:   1,
		UpdateAt: time.Now(),
		UpdateBy: u.Name,
	})
	if err != nil {
		http.Error(w, "视频转码错误", http.StatusInternalServerError)
		return
	}

	go func(v db.Video, name string) {
		ctx := context.Background()
		err := convert.ConvertHLS("media/"+v.ID+"/", strings.TrimPrefix(v.OriginLink, "/"))
		if err != nil {
			logger.Logger.Errorf("Convert HLS [%s]-[%s]: %v", v.ID, v.OriginLink, err)
			_, _ = server.store.UpdateVideoStatus(ctx, db.UpdateVideoStatusParams{
				ID:       v.ID,
				Status:   2,
				UpdateAt: time.Now(),
				UpdateBy: name,
			})
			return
		}

		// 转码成功
		if _, err = server.store.SetVideoPlay(ctx, db.SetVideoPlayParams{
			ID:       v.ID,
			Status:   200,
			PlayLink: "/media/" + v.ID + "/stream/",
			UpdateAt: time.Now(),
			UpdateBy: name,
		}); err != nil {
			logger.Logger.Errorf("Set Video Play [%s]-[%s]: %v", v.ID, v.OriginLink, err)
			return
		}

		logger.Logger.Infof("[%s]-[%s] 转码完成", v.ID, v.OriginLink)
	}(v, u.Name)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("视频正在转码中, 请稍后刷新页面"))
}
