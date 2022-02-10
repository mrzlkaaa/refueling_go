package download

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const (
	saveDir = "output"
)

var cwd, _ = os.Getwd()

type DownloadService interface {
	SavePDC(int, int) string
}

type StorageService interface {
	SavePDC(int) (string, *[]byte)
}

type downloadService struct {
	storage StorageService
}

func NewService(storage StorageService) DownloadService {
	return &downloadService{storage: storage}
}

func (s *downloadService) SavePDC(id int, actId int) string {
	homePath, _ := filepath.Split(cwd)
	name, pdc := s.storage.SavePDC(actId)
	dirPath := filepath.Join(homePath, saveDir, strconv.Itoa(id))
	err := os.MkdirAll(dirPath, 0644)
	if err != nil {
		log.Println(err)
	}
	filePath := filepath.Join(dirPath, name)
	err = os.WriteFile(filePath, *pdc, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(filePath)
	return filePath
}
