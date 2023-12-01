package fileutil

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
	"time"

	nanoId "github.com/matoous/go-nanoid"
)

var (
	// 图片最大数量
	MaxImageSize             = 10 << 20
	ErrUnsupportedFileFormat = errors.New("图片格式不支持")
	ErrFileStorePath         = errors.New("文件存储路径异常")
	ErrFileGenerateName      = errors.New("文件名称生成异常")
	ErrFileSaveFailed        = errors.New("文件保存失败")
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Mkdir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func UploadImage(r io.Reader) (string, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return "", fmt.Errorf("读取图片错误")
	}
	if errors.Is(err, image.ErrFormat) {
		return "", ErrUnsupportedFileFormat
	}
	if format != "png" && format != "jpeg" && format != "gif" {
		return "", ErrUnsupportedFileFormat
	}

	// 存放目录
	curTime := time.Now()
	dir := path.Join("upload", "imgs", curTime.Format("2006"), curTime.Format("01"), curTime.Format("02"))
	exist, _ := PathExists(dir)
	if !exist {
		err := Mkdir(dir)
		if err != nil {
			return "", ErrFileStorePath
		}
	}

	filename, err := nanoId.Nanoid()
	if err != nil {
		return "", ErrFileGenerateName
	}

	filePath := path.Join("", dir, filename+"."+format)
	f, err := os.Create(filePath)
	if err != nil {
		return "", ErrFileStorePath
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	var e error
	switch format {
	case "png":
		e = png.Encode(f, img)
	case "jpeg":
		e = jpeg.Encode(f, img, nil)
	case "gif":
		e = gif.Encode(f, img, nil)
	}
	if e != nil {
		return "", ErrFileSaveFailed
	}

	return "/" + filePath, nil
}
