package interfaces

type VideoProcessor interface {
	CreateHlsFile(videoSrc, outputDir, filename, uploadId string) error
}