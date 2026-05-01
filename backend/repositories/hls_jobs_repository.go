package repositories

import (
	"context"

	"github.com/arifazola/nontoon/internal/db"
	"github.com/google/uuid"
)

type HlsJobsRepository struct {
	Queries *db.Queries
}

func (repo *HlsJobsRepository) AddJob(context context.Context, uploadId string) error {
	params := db.AddHlsJobParams {
		ID: uuid.New().String(),
		UploadId: uploadId,
		Status: false,
	}
	return repo.Queries.AddHlsJob(context, params)
}

func (repo *HlsJobsRepository) UpdateJob(context context.Context, uploadId string, status bool) error {
	params := db.UpdateHlsJobStatusParams {
		Status: status,
		UploadId: uploadId,
	}
	return repo.Queries.UpdateHlsJobStatus(context, params)
}