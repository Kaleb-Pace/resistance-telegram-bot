package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// StichPicturesTogether stitches pictures together
func StichPicturesTogether(frames, filename string, framerate int) error {
	// ffmpeg -y -r 10 -i /F_%03d.png -c:v libx264 -vf fps=25 -pix_fmt yuv420p movie.mp4
	cmd := exec.Command("ffmpeg", "-y", "-r", fmt.Sprintf("%d", framerate), "-i", frames, "-c:v", "libx264", "-vf", fmt.Sprintf("fps=%d", framerate), "-pix_fmt", "yuv420p", filename)
	return cmd.Run()
}

func VideoFramerate(filename string) (int, error) {
	cmd := exec.Command("ffprobe", "-v", "0", "-of", "csv=p=0", "-select_streams", "v:0", "-show_entries", "stream=r_frame_rate", filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return -1, errors.New(string(output) + ": " + err.Error())
	}
	splitResults := strings.Split(string(output), "/")

	top, err := strconv.ParseFloat(splitResults[0], 64)
	if err != nil {
		return -1, err
	}
	bottom, err := strconv.ParseFloat(strings.TrimSpace(splitResults[1]), 64)
	if err != nil {
		return -1, err
	}

	if bottom <= 0 {
		return -1, errors.New("Parsed 0 out of denominator")
	}
	return int(top / bottom), nil
}

// UnstitchImages gets all images out of a mp4 and places them in the designated output
func UnstitchImages(videoPath, unstichOutput string) error {
	// ffmpeg -i file_7739.mp4 thumb%04d.png -hide_banner
	cmd := exec.Command("ffmpeg", "-i", videoPath, fmt.Sprintf("%s/%%04d.png", unstichOutput), "-hide_banner")
	return cmd.Run()
}
