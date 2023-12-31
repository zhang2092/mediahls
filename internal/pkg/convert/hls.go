package convert

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/zhang2092/mediahls/internal/pkg/fileutil"
	"github.com/zhang2092/mediahls/internal/pkg/logger"
	"github.com/zhang2092/mediahls/internal/pkg/rand"
)

func ConvertHLS(savePath, filePath string) error {
	exist, _ := fileutil.PathExists(savePath)
	if !exist {
		err := os.MkdirAll(savePath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	binary, err := exec.LookPath("ffmpeg")
	if err != nil {
		logger.Logger.Errorf("exec look path ffmpeg: %v", err)
		return err
	}

	// ffmpeg -i web/statics/git.mp4 -profile:v baseline -level 3.0 -s 1920x1080 -start_number 0 -hls_time 10 -hls_list_size 0 -hls_segment_filename %d.ts -f hls web/statics/git.m3u8
	command := "-i " + filePath + " -profile:v baseline -level 3.0 -s 1920x1080 -start_number 0 -hls_time 10 -hls_list_size 0 -f hls " + savePath + "index.m3u8"
	// log.Println(command)
	args := strings.Split(command, " ")
	cmd := exec.Command(binary, args...)
	_, err = cmd.Output()
	if err != nil {
		logger.Logger.Errorf("ffmpeg cmd output: %v", err)
		return err
	}

	// 替换 m3u8 文件信息
	return replaceM3u8(savePath)

	// ffmpeg -i upload/20231129/o6e6qKaMdk0VC1Ys2SHnr.mp4 -profile:v baseline -level 3.0 -s 1920x1080 -start_number 0 -f hls -hls_time 10 -segment_list web/statics/js/index.m3u8  web/statics/mmm/index.m3u8
	// ffmpeg -i upload/20231129/o6e6qKaMdk0VC1Ys2SHnr.mp4 -c copy -map 0 -f segment -segment_list web/statics/js/index.m3u8 -segment_time 20 web/statics/js/index_%3d.ts
}

type newOldName struct {
	old string
	new string
}

func replaceM3u8(savePath string) error {
	f, err := os.OpenFile(savePath+"/index.m3u8", os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	var result []newOldName
	var content string
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ls := string(line)
		if strings.HasSuffix(ls, ".ts") {
			ext := path.Ext(ls)
			filename := strings.TrimSuffix(ls, ext)
			temp := rand.RandomString(5)
			result = append(result, newOldName{
				old: ls,
				new: temp + ext,
			})
			ls = strings.ReplaceAll(ls, filename, temp)
		}
		content += ls + "\n"
	}

	fw, err := os.OpenFile(savePath+"/index.m3u8", os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer fw.Close()

	writer := bufio.NewWriter(fw)
	writer.WriteString(content)
	if err := writer.Flush(); err != nil {
		return err
	}

	for _, item := range result {
		if err := os.Rename(savePath+item.old, savePath+item.new); err != nil {
			return err
		}
	}

	return nil
}
