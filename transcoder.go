package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func transcodeToHLS(fileName string) {
	os.Mkdir(filepath.Join("transcoded", fileName), 0700)

	_, err := exec.Command(
		"ffmpeg",
		"-i", "media/"+fileName,
		"-profile:v", "baseline",
		"-level", "3.0", "-s", "640x360",
		"-start_number", "0",
		"-hls_time", "10", "-hls_list_size", "0",
		"-f", "hls", "transcoded/"+fileName+"/index.m3u8").Output()

	if err != nil {
		fmt.Println(err.Error())
	}
}
