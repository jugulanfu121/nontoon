package interfaces

import "context"

type VideoJobs interface {
	AddVideoJobs(context context.Context, id, uploadId string, index int) error
}