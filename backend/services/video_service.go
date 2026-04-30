package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/arifazola/nontoon/constants"
	"github.com/arifazola/nontoon/interfaces"
	"github.com/arifazola/nontoon/internal/db"
	"github.com/arifazola/nontoon/models"
	"github.com/google/uuid"
)

type VideoService struct{
	FileStorage interfaces.FileStorage
	FinalPath string
	VideoJobs interfaces.VideoJobs
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

func (videoService *VideoService) SaveChunk(uploadID string, index int, file io.ReadSeeker, ctx context.Context) error {
	err := videoService.FileStorage.SaveChunk(uploadID, index, file)

	if err != nil {
		return err
	}

	videoJobsId := uuid.New().String()

	addVideoJobsErr := videoService.VideoJobs.AddVideoJobs(ctx, videoJobsId, uploadID, index)

	if addVideoJobsErr != nil {
		return addVideoJobsErr
	}

	return nil
}

func (videoService *VideoService) MergeChunks(uploadId, filename string, totalChunks int, basepath string) (string, error){
	err := os.MkdirAll(basepath, os.ModePerm)

	if err != nil {
		log.Println("error create dir: ", err)
	}

	finalPath := filepath.Join(basepath, filename)


	finalFile, err := os.Create(finalPath)

	if err != nil {
		log.Println("error creating final file: ", err)
		return "", err
	}


	for i := 0; i < totalChunks; i ++ {
		chunkPath := filepath.Join(basepath, uploadId, fmt.Sprintf("%d.part", i))

		chunkFile, err := os.Open(chunkPath)

		if err != nil {
			log.Println("error opening chunk file: ", err)
			return "", err
		}

		_, copyChunkErr := io.Copy(finalFile, chunkFile)

		chunkFile.Close()

		if copyChunkErr != nil {
			log.Println("Error copying chunks: ", copyChunkErr)
			return "", copyChunkErr
		}
	}

	finalFile.Close()

	orgFile := basepath + "/" + filename
	destinationFile := constants.ASSETS_PATH + "/" + filename
	movFileRrr := os.Rename(orgFile, destinationFile)

	if movFileRrr != nil {
		log.Println("move file error", movFileRrr)
		return "", movFileRrr
	}

	go func ()  {
		time.Sleep(10 * time.Second)
		videoSrc := constants.ASSETS_PATH + "/" + filename
		output := constants.ASSETS_PATH + "/def" + filename + ".m3u8"
		CreateHlsFile(videoSrc, output, filename)
	}()

	return finalPath, nil
}

func (service *VideoService) GetLatestUploadedChunk(ctx context.Context, uploadId string) (db.VideoJob, error){
	latestChunk, err := service.VideoJobs.GetLatestUploadedChunk(ctx, uploadId)

	if err != nil {
		var videoJob db.VideoJob
		return videoJob, err
	}

	return latestChunk, nil
}