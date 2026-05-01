package repositories

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/arifazola/nontoon/constants"
)

type VideoProcessorRepository struct{}

func (v *VideoProcessorRepository) CreateHlsFile(videoSrc, outputDir, filename, uploadId string) error {
	path := filepath.Join(constants.ASSETS_PATH, uploadId)
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		log.Println("Error creating hls folder", err)
	}

	segmentFile := filepath.Join("assets", uploadId, "v%v", "file_%03d.ts")
	hlsFile := filepath.Join("assets", uploadId, "v%v", "prog_index.m3u8")

	out, err := exec.Command(
		"ffmpeg",
		"-i", videoSrc,
		"-filter_complex", "[0:v]split=2[v1][v2];[v1]scale=1280:720[v1out];[v2]scale=640:360[v2out]",
		"-map", "[v1out]",
		"-map", "0:a",
		"-map", "[v2out]",
		"-map", "0:a",
		"-c:v", "libx264",
		"-c:a", "aac",
		"-f", "hls",
		"-var_stream_map", "v:0,a:0 v:1,a:1",
		"-master_pl_name", "master.m3u8",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-hls_segment_filename", segmentFile,
		hlsFile,
	).CombinedOutput()

	if err != nil {
		fmt.Println("Error Hls", err)
		return err
	}

	fmt.Println(string(out))

	return nil
}