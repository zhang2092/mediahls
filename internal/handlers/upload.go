package handlers

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	nanoId "github.com/matoous/go-nanoid"
	"github.com/zhang2092/mediahls/internal/pkg/fileutil"
)

func (server *Server) uploadVideo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("%v", err)
		}
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
	_, err = w.Write([]byte("/" + filePath))
	if err != nil {
		log.Printf("%v", err)
	}
}

func (server *Server) uploadImage(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	_, fh, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("%v", err)
		}
		return
	}

	f, err := fh.Open()
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("读取图片失败"))
		if err != nil {
			log.Printf("%v", err)
		}
		return
	}
	reader := bufio.NewReader(f)
	filePath, err := fileutil.UploadImage(reader)
	if errors.Is(err, fileutil.ErrUnsupportedFileFormat) {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		_, err = w.Write([]byte(fileutil.ErrUnsupportedFileFormat.Error()))
		if err != nil {
			log.Printf("%v", err)
		}
		return
	}

	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("%v", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(filePath))
	if err != nil {
		log.Printf("%v", err)
	}
}
