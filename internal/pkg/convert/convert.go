package convert

import (
	"os/exec"
	"strings"
)

func ConvertHLS(savePath, filePath string) error {
	binary, err := exec.LookPath("ffmpeg")
	if err != nil {
		return err
	}

	// ffmpeg -i web/statics/git.mp4 -profile:v baseline -level 3.0 -s 1920x1080 -start_number 0 -hls_time 10 -hls_list_size 0 -hls_segment_filename %d.ts -f hls web/statics/git.m3u8
	command := "-i " + filePath + " -profile:v baseline -level 3.0 -s 1920x1080 -start_number 0 -hls_time 10 -hls_list_size 0 -f hls " + savePath + "index.m3u8"
	args := strings.Split(command, " ")
	cmd := exec.Command(binary, args...)
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	return nil
}
