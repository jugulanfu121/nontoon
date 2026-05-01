package interfaces

import "io"

type FileStorage interface {
	Save(file io.ReadSeeker, filename string) (string, error)
	SaveChunk(uploadID string, index int, file io.ReadSeeker) error
	MergeChunks(uploadId, filename string, totalChunks int, basepath string) (string, error)
}