package services_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/arifazola/nontoon/internal/db"
	"github.com/arifazola/nontoon/repositories"
	"github.com/arifazola/nontoon/services"
	"github.com/stretchr/testify/assert"
)

func TestSaveVideoSuccess(t *testing.T) {
	mockStorage := repositories.MockStorage{}

	service := services.VideoService {
		FileStorage: &mockStorage,
	}

	fileContent := []byte("hello world")
    reader := bytes.NewReader(fileContent)

    _, err := service.SaveVideo(reader, "test.txt", int64(len(fileContent)))
	
	assert.Nil(t, err, "Error must be nil")
}

func TestSaveVideoError(t *testing.T) {
	mockStorage := repositories.MockStorage{
		Err: errors.New("disk full"),
	}

	service := services.VideoService {
		FileStorage: &mockStorage,
	}

	fileContent := []byte("hello world")
    reader := bytes.NewReader(fileContent)

    _, err := service.SaveVideo(reader, "test.txt", int64(len(fileContent)))
	
	assert.NotNil(t, err, "Error must exist")
}

func TestSaveChunkSuccess(t *testing.T){
	mockStorage := repositories.MockStorage{}
	mockVideoJobs := repositories.MockVideoJobsRepository{}

	service := services.VideoService {
		VideoJobs: &mockVideoJobs,
		FileStorage: &mockStorage,
	}

	fileContent := []byte("hello world")
    reader := bytes.NewReader(fileContent)

    err := service.SaveChunk("mockID", "test_name", 0, reader, t.Context())
	
	assert.Nil(t, err, "Error must be nil")
}

func TestSaveChunkError(t *testing.T){
	mockStorage := repositories.MockStorage{
		Err: errors.New("disk full"),
	}

	mockVideoJobs := repositories.MockVideoJobsRepository{}

	service := services.VideoService {
		FileStorage: &mockStorage,
		VideoJobs: &mockVideoJobs,
	}

	fileContent := []byte("hello world")
    reader := bytes.NewReader(fileContent)

    err := service.SaveChunk("mockID", "test_name", 0, reader, t.Context())
	
	assert.NotNil(t, err, "Error must exist")
}

func TestSaveChunkWriteDbFail(t *testing.T){
	mockStorage := repositories.MockStorage{}
	mockVideoJobs := repositories.MockVideoJobsRepository{
		Err: errors.New("Failed To Write To DB"),
	}

	service := services.VideoService {
		VideoJobs: &mockVideoJobs,
		FileStorage: &mockStorage,
	}

	fileContent := []byte("hello world")
    reader := bytes.NewReader(fileContent)

    err := service.SaveChunk("mockID", "test_name", 0, reader, t.Context())
	
	assert.NotNil(t, err, "Error must be nil")
}

func TestMergeChunks_Success(t *testing.T) {
	tempDir := t.TempDir()

	uploadId := "test-upload"
	filename := "final.txt"
	totalChunks := 3

	basePath := filepath.Join(tempDir)
	chunkDir := filepath.Join(basePath, uploadId)

	err := os.MkdirAll(chunkDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	chunks := []string{"hello ", "world", "!"}

	for i, content := range chunks {
		chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d.part", i))
		err := os.WriteFile(chunkPath, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	localStorage := repositories.LocalStorage{
		BasePath: "./files",
	}

	videoProcessor := repositories.VideoProcessorRepository{}

	service := &services.VideoService{
		FileStorage: &localStorage,
		VideoProcessor: &videoProcessor,
	}

	finalPath, err := service.MergeChunks(uploadId, filename, totalChunks, basePath)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(finalPath)
	if err != nil {
		t.Fatal(err)
	}

	expected := "hello world!"

	assert.Equal(t, expected, string(data))
}

// func TestMergeChunks_MissingChunk(t *testing.T) {
// 	tempDir := t.TempDir()

// 	uploadId := "test-upload"
// 	filename := "final.txt"
// 	totalChunks := 2

// 	basePath := tempDir
// 	chunkDir := filepath.Join(basePath, uploadId)

// 	os.MkdirAll(chunkDir, os.ModePerm)

// 	os.WriteFile(filepath.Join(chunkDir, "0.part"), []byte("hello"), 0644)

// 	localStorage := repositories.LocalStorage{
// 		BasePath: "./files",
// 	}

// 	videoProcessor := repositories.VideoProcessorRepository{}

// 	service := &services.VideoService{
// 		FileStorage: &localStorage,
// 		VideoProcessor: &videoProcessor,
// 	}

// 	_, err := service.MergeChunks(uploadId, filename, totalChunks, basePath)

// 	if err != nil {
// 		fmt.Println("Should pass")
// 	}

// 	assert.NotNil(t, err, "Error should exist")
// }

func TestGetLatestUploadChunk_Success(t *testing.T){
	videoJobs := repositories.MockVideoJobsRepository{}
	service := &services.VideoService{
		VideoJobs: &videoJobs,
	}

	res, err := service.GetLatestUploadedChunk(t.Context(), "abc")

	assert.Nil(t, err, "Error should be nil")

	var videoJob db.VideoJob

	videoJob.ID = "zzzz"
	videoJob.UploadId = "abc"
	videoJob.Index = int32(0)

	assert.Equal(t, videoJob, res)
}

func TestGetLatestUploadChunk_Error(t *testing.T){
	videoJobs := repositories.MockVideoJobsRepository{
		Err: errors.New("Failed to get data from db"),
	}
	service := &services.VideoService{
		VideoJobs: &videoJobs,
	}

	_, err := service.GetLatestUploadedChunk(t.Context(), "abc")

	assert.NotNil(t, err, "Error should exist")
}

