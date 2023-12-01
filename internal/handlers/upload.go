package handlers

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	nanoId "github.com/matoous/go-nanoid"
	"github.com/zhang2092/mediahls/internal/pkg/fileutil"
	"github.com/zhang2092/mediahls/internal/pkg/logger"
)

// data

// uploadVideo 上传视频
func (server *Server) uploadVideo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		logger.Logger.Errorf("upload video receive: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// filetype := http.DetectContentType(buff)
	// if filetype != "image/jpeg" && filetype != "image/png" {
	// 	http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
	// 	return
	// }

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	curTime := time.Now()
	dir := path.Join("upload", "files", curTime.Format("2006"), curTime.Format("01"), curTime.Format("02"))
	exist, _ := fileutil.PathExists(dir)
	if !exist {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	filename, err := nanoId.Nanoid()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filePath := path.Join("", dir, filename+filepath.Ext(fileHeader.Filename))
	f, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("/" + filePath))
}

// uploadImage 上传图片
func (server *Server) uploadImage(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	_, fh, err := r.FormFile("file")
	if err != nil {
		logger.Logger.Errorf("upload image receive: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	f, err := fh.Open()
	if err != nil {
		logger.Logger.Errorf("upload image read: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("读取图片失败"))
		return
	}
	reader := bufio.NewReader(f)
	filePath, err := fileutil.UploadImage(reader)
	if errors.Is(err, fileutil.ErrUnsupportedFileFormat) {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fileutil.ErrUnsupportedFileFormat.Error()))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(filePath))
}
