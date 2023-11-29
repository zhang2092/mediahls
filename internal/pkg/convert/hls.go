package convert

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/zhang2092/mediahls/internal/pkg/fileutil"
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
		log.Println("1: ", err)
		return err
	}

	// ffmpeg -i web/statics/git.mp4 -profile:v baseline -level 3.0 -s 1920x1080 -start_number 0 -hls_time 10 -hls_list_size 0 -hls_segment_filename %d.ts -f hls web/statics/git.m3u8
	command := "-i " + filePath + " -profile:v baseline -level 3.0 -s 1920x1080 -start_number 0 -hls_time 10 -hls_list_size 0 -f hls " + savePath + "index.m3u8"
	log.Println(command)
	args := strings.Split(command, " ")
	cmd := exec.Command(binary, args...)
	_, err = cmd.Output()
	if err != nil {
		log.Println("2: ", err)
		return err
	}

	return nil
	// ffmpeg -i upload/20231129/o6e6qKaMdk0VC1Ys2SHnr.mp4 -profile:v baseline -level 3.0 -s 1920x1080 -start_number 0 -f hls -hls_time 10 -segment_list web/statics/js/index.m3u8  web/statics/mmm/index.m3u8
	// ffmpeg -i upload/20231129/o6e6qKaMdk0VC1Ys2SHnr.mp4 -c copy -map 0 -f segment -segment_list web/statics/js/index.m3u8 -segment_time 20 web/statics/js/index_%3d.ts
}
