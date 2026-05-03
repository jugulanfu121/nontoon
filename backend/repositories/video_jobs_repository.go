package repositories

import (
	"context"

	"github.com/arifazola/nontoon/internal/db"
)

type VideoJobsRepository struct {
	Queries *db.Queries
}

func (repo *VideoJobsRepository) AddVideoJobs(context context.Context, id, uploadId, filename string, index int) error {
	return repo.Queries.AddVideoJob(context, db.AddVideoJobParams{
		ID: id,
		UploadId: uploadId,
		Index: int32(index),
		Filename: filename,
	})
}

func (repo *VideoJobsRepository) GetLatestUploadedChunk(ctx context.Context, uploadId string) (db.GetLatestUploadedChunkRow, error) {
	latestChunk, err := repo.Queries.GetLatestUploadedChunk(ctx, uploadId)

	var videoJob db.GetLatestUploadedChunkRow
	if err != nil {
		return videoJob, err
	}

	return latestChunk, nil
}

func (repo *VideoJobsRepository) GetVideoJobByFilename(ctx context.Context, filename string) (db.GetLatestUploadedChunkByFilenameRow, error){
	videoJob, err := repo.Queries.GetLatestUploadedChunkByFilename(ctx, filename)

	return videoJob, err
}