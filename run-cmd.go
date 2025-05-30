package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func runFfmpeg() (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	const cmdName = "ffmpeg"
	const cmdWorkingDirectory = "/home/pagans/Desktop/shrimps"
	videoFileName := filepath.Join(cmdWorkingDirectory, "merged_rotated_fitted_1080p.mp4")
	audioFileName := filepath.Join(cmdWorkingDirectory, "lofi-free-music-yt.aac")

	cmdArgs := []string{
		"-re",
		"-stream_loop", "-1",
		"-i", videoFileName,
		"-thread_queue_size", "512",
		"-stream_loop", "-1",
		"-i", audioFileName,
		"-max_delay", "5000000",
		"-shortest",
		"-c:v", "copy",
		"-c:a", "aac",
		"-b:a", "128k",
		"-f", "flv",
		os.Getenv("YOUTUBE_STREAMING_URL"),
	}

	streamCmd = exec.Command(cmdName, cmdArgs...)

	stdout, _ := streamCmd.StdoutPipe()
	stderr, _ := streamCmd.StderrPipe()

	if err := streamCmd.Start(); err != nil {
		return nil, nil, nil, fmt.Errorf("error starting stream: %v", err)
	}

	return streamCmd, stdout, stderr, nil
}
