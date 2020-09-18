package camera

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"time"
)

func TakePhoto() ([]byte, error) {
	cmd := exec.Command("raspistill", "-t", "1", "-o", "-")

	var bts bytes.Buffer
	var btsErr bytes.Buffer

	cmd.Stderr = &btsErr
	cmd.Stdout = &bts

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("can't take photo: %s", btsErr.String())
	}

	return bts.Bytes(), nil
}

func TakeVideo(duration time.Duration) ([]byte, error) {
	const filenameH264 = "/tmp/vid.h264"
	const filenameMPEG4 = "/tmp/vid.mp4"
	timing := strconv.FormatInt(duration.Milliseconds(), 10)

	cmd := exec.Command("raspivid", "-t", timing, "-b", "1000000", "-fps", "25", "-h", "600", "-w", "800", "-o", filenameH264)

	var btsErr bytes.Buffer

	cmd.Stderr = &btsErr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("can't take video: %s", btsErr.String())
	}

	cmd = exec.Command("ffmpeg", "-y", "-i", filenameH264, filenameMPEG4)

	cmd.Stderr = &btsErr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("can't format video: %s", btsErr.String())
	}

	bts, err := ioutil.ReadFile(filenameMPEG4)
	if err != nil {
		return nil, fmt.Errorf("can't read videofile: %w", err)
	}

	return bts, nil
}
