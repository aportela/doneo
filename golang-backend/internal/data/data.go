package data

import (
	"log"
	"os"
	"path/filepath"
)

func GetDataPath() string {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}
	return filepath.Join(currentPath, "data")
}

func CreateDataPathIfRequired() error {
	if err := os.MkdirAll(GetDataPath(), os.ModePerm); err != nil {
		return err
	}
	return nil
}
