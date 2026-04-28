package services

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/arifazola/nontoon/constants"
	"github.com/arifazola/nontoon/interfaces"
	"github.com/arifazola/nontoon/models"
)

type VideoService struct{
	FileStorage interfaces.FileStorage
	FinalPath string
}

func GetAllVideos() []models.Video {
	var videos = []models.Video{
		{ID: "1", Title: "Blue Train", Artist: "John Coltrane"},
		{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan"},
		{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan"},
	}

	return videos
}

func (videoService *VideoService) SaveVideo(file io.ReadSeeker, filename string, size int64) (string, error) {

	saveVideo, err := videoService.FileStorage.Save(file, filename)

	if err != nil {
		return "", err
	}

	return saveVideo, nil
}

func (videoService *VideoService) SaveChunk(uploadID string, index int, file io.ReadSeeker) error {
	videoService.FileStorage.SaveChunk(uploadID, index, file)

	return nil
}

func (videoService *VideoService) MergeChunks(uploadId, filename string, totalChunks int) (string, error){
	err := os.MkdirAll(constants.BASE_PATH, os.ModePerm)

	if err != nil {
		log.Println("error create dir: ", err)
	}

	finalPath := filepath.Join(constants.BASE_PATH, filename)


	finalFile, err := os.Create(finalPath)

	if err != nil {
		log.Println("error creating final file: ", err)
	}

	defer finalFile.Close()

	for i := 0; i < totalChunks; i ++ {
		chunkPath := filepath.Join(constants.BASE_PATH, uploadId, fmt.Sprintf("%d.part", i))

		chunkFile, err := os.Open(chunkPath)

		if err != nil {
			log.Println("error opening chunk file: ", err)
		}

		_, copyChunkErr := io.Copy(finalFile, chunkFile)

		if copyChunkErr != nil {
			log.Println("Error copying chunks: ", copyChunkErr)
		}

		defer chunkFile.Close()
	}

	return "", nil
}