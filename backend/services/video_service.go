package services

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

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
	HlsJob interfaces.HlsJobs
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

func (videoService *VideoService) CompleteUpload(uploadId, filename string, totalChunks int, basepath string, context context.Context) (string, error){
	
	finalPath, err := videoService.MergeChunks(uploadId, filename, totalChunks, basepath)

	if err != nil {
		return  "", err
	}

	errMoveToAssetFolder := videoService.MoveFileToAssetsFolder(basepath, filename)

	if errMoveToAssetFolder != nil {
		return "", errMoveToAssetFolder
	}

	go func ()  {
		errAddHlsJob := videoService.HlsJob.AddJob(context, uploadId)

		if errAddHlsJob != nil {
			log.Println("error adding hls job: ", errAddHlsJob)
		}

		errGenerateHls := videoService.GenerateHls(filename, uploadId)

		if errGenerateHls != nil {
			log.Println("error generating hls: ", errGenerateHls)
		}

		errUpdateHlsJob := videoService.HlsJob.UpdateJob(context, uploadId, true)

		if errUpdateHlsJob != nil {
			log.Println("error adding hls job: ", errUpdateHlsJob)
		}

		videoService.DeleteTempFile(uploadId, filename)
	}()

	return finalPath, nil
}

func (service *VideoService) MergeChunks(uploadId, filename string, totalChunks int, basepath string) (string, error){
	finalPath, err := service.FileStorage.MergeChunks(uploadId, filename, totalChunks, basepath)

	if err != nil {
		return  "", err
	}

	return finalPath, nil
}

func (service *VideoService) GenerateHls(filename, uploadId string) error {
	videoSrc := constants.ASSETS_PATH + "/" + filename
	output := filepath.Join(constants.ASSETS_PATH, uploadId, filename + ".m3u8")
	errHls := service.VideoProcessor.CreateHlsFile(videoSrc, output, filename, uploadId)

	if errHls != nil {
		log.Println("error hls result:", errHls)
		return errHls
	}

	return nil
}

func (service *VideoService) MoveFileToAssetsFolder(basepath, filename string) error{
	//Move the merged file to assets folder before being processed by ffmpeg
	orgFile := basepath + "/" + filename
	destinationFile := constants.ASSETS_PATH + "/" + filename
	movFileRrr := os.Rename(orgFile, destinationFile)

	if movFileRrr != nil {
		log.Println("move file error", movFileRrr)
		return movFileRrr
	}

	return nil
}

func (service *VideoService) DeleteTempFile(uploadId, filename string){
	chunkPath := filepath.Join(constants.BASE_PATH, uploadId)
	deleteChunkFileError := service.FileStorage.DeleteFile(chunkPath)

	if deleteChunkFileError != nil {
		log.Println("error delete chunk path", deleteChunkFileError)
	}

	tempAssetFile := filepath.Join(constants.ASSETS_PATH, filename)
	deleteTempAssetErr := service.FileStorage.DeleteFile(tempAssetFile)

	if deleteTempAssetErr != nil {
		log.Println("error delete chunk path", deleteChunkFileError)
	}
}

func (service *VideoService) GetLatestUploadedChunk(ctx context.Context, uploadId string) (db.VideoJob, error){
	latestChunk, err := service.VideoJobs.GetLatestUploadedChunk(ctx, uploadId)

	if err != nil {
		var videoJob db.VideoJob
		return videoJob, err
	}

	return latestChunk, nil
}