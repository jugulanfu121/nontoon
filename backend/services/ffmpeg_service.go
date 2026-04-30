package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/arifazola/nontoon/constants"
)


func CreateHlsFile(videoSrc, outputDir, filename string) {
	// out, _ := exec.Command("ffmpeg", "-i", videoSrc, "-codec:", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", outputDir).CombinedOutput()
	path := filepath.Join(constants.ASSETS_PATH, "def")
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		log.Println("Error creating hls folder", err)
	}
	
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
		"-hls_segment_filename", "assets/def/v%v/file_%03d.ts",
		"assets/def/v%v/prog_index.m3u8",
	).CombinedOutput()

	if err != nil {
		fmt.Println("Error Hls", err)
		return
	}
	
	fmt.Println(string(out))
}