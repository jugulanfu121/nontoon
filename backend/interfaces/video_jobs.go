package interfaces

import (
	"context"

	"github.com/arifazola/nontoon/internal/db"
)

type VideoJobs interface {
	AddVideoJobs(context context.Context, id, uploadId, filename string, index int) error
	GetLatestUploadedChunk(ctx context.Context, uploadId string) (db.GetLatestUploadedChunkRow, error)
	GetVideoJobByFilename(ctx context.Context, filename string) (db.GetLatestUploadedChunkByFilenameRow, error)
}