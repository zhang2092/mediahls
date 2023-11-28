package handlers

import (
	"html/template"
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

func (server *Server) upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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

		dir := path.Join("upload", time.Now().Format("20060102"))
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
		_, err = w.Write([]byte(filePath))
		if err != nil {
			log.Printf("%v", err)
		}

		return
	}

	t, err := template.ParseFiles("web/templates/upload.html.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
