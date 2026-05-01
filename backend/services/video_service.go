package services

import (
	"context"
	"io"
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
	VideoProcessor interfaces.VideoProcessor
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
	
	finalPath, err := videoService.FileStorage.MergeChunks(uploadId, filename, totalChunks, basepath)

	if err != nil {
		return  "", err
	}

	go func ()  {
		time.Sleep(10 * time.Second)
		videoSrc := constants.ASSETS_PATH + "/" + filename
		output := constants.ASSETS_PATH + "/def" + filename + ".m3u8"
		videoService.VideoProcessor.CreateHlsFile(videoSrc, output, filename)
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