package interfaces

type VideoProcessor interface {
	CreateHlsFile(videoSrc, outputDir, filename string) error
}