package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func transcodeToHLS(fileName string) {
	os.Mkdir(filepath.Join("freeschool", "transcoded", fileName), 0700)

	_, err := exec.Command(
		"ffmpeg",
		"-i", "freeschool/media/"+fileName,
		"-profile:v", "baseline",
		"-level", "3.0", "-s", "854x480",
		"-start_number", "0",
		"-hls_time", "10", "-hls_list_size", "0",
		"-f", "hls", "freeschool/transcoded/"+fileName+"/index.m3u8").Output()

	if err != nil {
		fmt.Println(err.Error())
	}
}
