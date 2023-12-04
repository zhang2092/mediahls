package handlers

import (
	"context"
	"encoding/json"
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

// obj

// videoPageData 播放页面数据
type videoPageData struct {
	Authorize
	Video db.Video
}

// videosPageData 视频列表数据
type videosPageData struct {
	Authorize
	Videos []db.Video
}

// videoEditPageData 视频编辑数据
type videoEditPageData struct {
	Authorize
	Summary        string
	ID             string
	IDMsg          string
	Title          string
	TitleMsg       string
	Images         string
	ImagesMsg      string
	Description    string
	DescriptionMsg string
	OriginLink     string
	OriginLinkMsg  string
	Status         int
	StatusMsg      string
}

// videoDeleteRequest 视频删除请求参数
type videoDeleteRequest struct {
	ID string `json:"id"`
}

// view

// videoView 视频播放页面
func (server *Server) videoView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	xid := vars["xid"]
	data := videoPageData{}
	video, err := server.store.GetVideo(r.Context(), xid)
	if err == nil {
		data.Video = video
	}
	auth, err := server.withCookie(r)
	if err == nil {
		data.Authorize = *auth
	}
	server.renderLayout(w, r, data, "video/play.html.tmpl")
}

// videosView 视频列表页面
func (server *Server) videosView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := videosPageData{
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

	server.renderLayout(w, r, data, "video/videos.html.tmpl")
}

// editVideoView 视频编辑页面
func (server *Server) editVideoView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	xid := vars["xid"]
	vm := videoEditPageData{
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
	server.renderEditVideo(w, r, vm)
}

// data

// stream 视频HLS播放处理
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

// editVideo 视频编辑
func (server *Server) editVideo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		server.renderEditVideo(w, r, videoEditPageData{Summary: "请求网络错误, 请刷新重试"})
		return
	}

	vm, ok := viladatorEditVedio(r)
	if !ok {
		server.renderEditVideo(w, r, vm)
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
			server.renderEditVideo(w, r, vm)
			return
		}
	} else {
		v, err := server.store.GetVideo(ctx, vm.ID)
		if err != nil {
			vm.Summary = "视频数据错误"
			server.renderEditVideo(w, r, vm)
			return
		}

		var sta int32 = int32(vm.Status)
		if sta != -1 && sta != 200 {
			sta = v.Status
		}

		_, err = server.store.UpdateVideo(ctx, db.UpdateVideoParams{
			ID:          vm.ID,
			Title:       vm.Title,
			Description: vm.Description,
			Images:      vm.Images,
			Status:      sta,
			UpdateAt:    curTime,
			UpdateBy:    u.Name,
		})
		if err != nil {
			vm.Summary = "更新视频失败"
			server.renderEditVideo(w, r, vm)
			return
		}
	}

	http.Redirect(w, r, "/me/videos", http.StatusFound)
}

// deleteVideo 视频删除
func (server *Server) deleteVideo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req videoDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondErr(w, "参数错误", nil)
		return
	}

	err := server.store.DeleteVideo(r.Context(), req.ID)
	if err != nil {
		RespondErr(w, "删除失败", nil)
		return
	}

	Respond(w, "删除成功", nil, http.StatusOK)
}

// transfer 视频转码
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

// method

// renderEditVideo 渲染视频编辑页面
func (server *Server) renderEditVideo(w http.ResponseWriter, r *http.Request, data any) {
	server.renderLayout(w, r, data, "video/edit.html.tmpl")
}

// viladatorEditVedio 检验视频编辑数据
func viladatorEditVedio(r *http.Request) (videoEditPageData, bool) {
	ok := true
	status, _ := strconv.Atoi(r.PostFormValue("status"))
	resp := videoEditPageData{
		Authorize:   withUser(r.Context()),
		ID:          r.PostFormValue("id"),
		Title:       r.PostFormValue("title"),
		Images:      r.PostFormValue("images"),
		Description: r.PostFormValue("description"),
		OriginLink:  r.PostFormValue("origin_link"),
		Status:      status,
	}

	if len(resp.Title) == 0 {
		resp.TitleMsg = "请填写正确的标题"
		ok = false
	}

	exist, _ := fileutil.PathExists(strings.TrimPrefix(resp.Images, "/"))
	if !exist {
		resp.ImagesMsg = "请先上传图片"
		ok = false
	}

	if len(resp.Description) == 0 {
		resp.DescriptionMsg = "请填写描述"
		ok = false
	}

	exist, _ = fileutil.PathExists(strings.TrimPrefix(resp.OriginLink, "/"))
	if !exist {
		resp.OriginLinkMsg = "请先上传视频"
		ok = false
	}

	return resp, ok
}

// getMediaBase 获取视频m3u8文件路径
func getMediaBase(mId string) string {
	mediaRoot := "media"
	return fmt.Sprintf("%s/%s", mediaRoot, mId)
}

// serveHlsM3u8 返回m3u8文件
func serveHlsM3u8(w http.ResponseWriter, r *http.Request, mediaBase, m3u8Name string) {
	mediaFile := fmt.Sprintf("%s/%s", mediaBase, m3u8Name)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "application/x-mpegURL")
}

// serveHlsTs 返回ts文件
func serveHlsTs(w http.ResponseWriter, r *http.Request, mediaBase, segName string) {
	mediaFile := fmt.Sprintf("%s/%s", mediaBase, segName)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "video/MP2T")
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
