package interfaces

import "context"

type HlsJobs interface {
	AddJob(context context.Context, uploadId string) error
	UpdateJob(context context.Context, uploadId string, status bool) error
}