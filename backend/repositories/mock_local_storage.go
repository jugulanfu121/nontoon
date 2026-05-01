package repositories

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type MockStorage struct {
    SavedFiles map[string][]byte
    Err        error
}

func (m *MockStorage) Save(file io.ReadSeeker, filename string) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}

	b, _ := io.ReadAll(file)

	if m.SavedFiles == nil {
		m.SavedFiles = make(map[string][]byte)
	}

	m.SavedFiles[filename] = b

	return "/mock/" + filename, nil
}

func (m *MockStorage) SaveChunk(uploadID string, index int, file io.ReadSeeker) error {
	if m.Err != nil {
		return m.Err
	}
	
	path := filepath.Join("/mockfile", uploadID)

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

func (m *MockStorage) MergeChunks(uploadId, filename string, totalChunks int, basepath string) (string, error){
	return "", errors.New("")
}