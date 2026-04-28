package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/arifazola/nontoon/internal/db"
)

type VideoJobsRepository struct {
	Queries *db.Queries
}

func (repo *VideoJobsRepository) AddVideoJobs(context context.Context, id, uploadId string, index int) error {
	var index32 = sql.NullInt32{
		Int32: int32(index),
		Valid: true,
	}

	fmt.Println("Index Video Jobs", index32)
	fmt.Println("Index Video Jobs", index)

	return repo.Queries.AddVideoJob(context, db.AddVideoJobParams{
		ID: id,
		UploadId: uploadId,
		Index: index32,
	})
}