package repositories

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	BasePath string
}

func (localStorage *LocalStorage) Save(file io.ReadSeeker, filename string) (string, error) {
	safeName := filepath.Base(filename)

	path := filepath.Join(localStorage.BasePath, safeName)

	dst, err := os.Create(path)

	if err != nil {
		log.Println("error upload video: ", err)
		return "", err
	}

	defer dst.Close()

	_, copyErr := io.Copy(dst, file)

	if copyErr != nil {
		log.Println("error copy video: ", err)
		return "", err
	}

	return path, nil
}

func (localStorage *LocalStorage) SaveChunk(uploadID string, index int, file io.ReadSeeker) error{
	path := filepath.Join(localStorage.BasePath, uploadID)

	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		log.Println("error saving chunks: ", err)
		return err
	}

	chunkPath := filepath.Join(path, fmt.Sprintf("%d.part", index))

	chunkFile, createErr := os.Create(chunkPath)

	if createErr != nil {
		log.Println("error create chunks: ", createErr)
		return createErr
	}

	defer chunkFile.Close()

	_, copyErr := io.Copy(chunkFile, file)

	if copyErr != nil {
		log.Println("Error copy chunks: ", copyErr)
		return copyErr
	}

	return nil
}

func (localStorage *LocalStorage) MergeChunks(uploadId, filename string, totalChunks int, basepath string) (string, error){
	err := os.MkdirAll(basepath, os.ModePerm)

	if err != nil {
		log.Println("error create dir: ", err)
	}

	finalPath := filepath.Join(basepath, filename)

	finalFile, err := os.Create(finalPath)

	if err != nil {
		log.Println("error creating final file: ", err)
		return "", err
	}


	for i := 0; i < totalChunks; i ++ {
		chunkPath := filepath.Join(basepath, uploadId, fmt.Sprintf("%d.part", i))

		chunkFile, err := os.Open(chunkPath)

		if err != nil {
			log.Println("error opening chunk file: ", err)
			return "", err
		}

		_, copyChunkErr := io.Copy(finalFile, chunkFile)

		chunkFile.Close()

		if copyChunkErr != nil {
			log.Println("Error copying chunks: ", copyChunkErr)
			return "", copyChunkErr
		}
	}

	finalFile.Close()

	return finalPath, nil
}

func (localStorage *LocalStorage) DeleteFile(path string) error {
	return os.Remove(path)
}