package repositories

import (
	"context"
	"errors"

	"github.com/arifazola/nontoon/internal/db"
)

type MockVideoJobsRepository struct{
	Err error
}

func (m *MockVideoJobsRepository) AddVideoJobs(context context.Context, id, uploadId, filename string, index int) error {
	if m.Err != nil {
		return errors.New("Failed To Store Jobs")
	}

	return nil
}

func (m *MockVideoJobsRepository) GetLatestUploadedChunk(ctx context.Context, uploadId string) (db.GetLatestUploadedChunkRow, error) {
	var videoJob db.GetLatestUploadedChunkRow
	if m.Err != nil {
		return videoJob, errors.New("Err")
	}

	videoJob.ID = "zzzz"
	videoJob.UploadId = "abc"
	index := int32(0)
	videoJob.Index = index

	return videoJob, nil

}

func (m *MockVideoJobsRepository) GetVideoJobByFilename(ctx context.Context, filename string) (db.GetLatestUploadedChunkByFilenameRow, error){
	var videoJob db.GetLatestUploadedChunkByFilenameRow
	if m.Err != nil {
		return videoJob, errors.New("Err")
	}

	videoJob.ID = "zzzz"
	videoJob.UploadId = "abc"
	index := int32(0)
	videoJob.Index = index

	return videoJob, nil
}