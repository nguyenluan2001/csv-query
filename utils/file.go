package utils

import (
	"log"
	"os"
	"os/exec"
)

func NewFile(filepath string) (os.FileInfo, error) {
	fileInfo, err := os.Stat(filepath)

	//File existed
	if err == nil {
		return fileInfo, err
	}

	// File not existed
	file, err := os.Create(filepath)
	if err != nil {
		panic("Create file failed")
	}

	fileInfo, err = file.Stat()
	return fileInfo, err
}

func EmptyFile(filepath string) {
	_, err := os.Stat(filepath)

	if err != nil {
		log.Println("File not found. Please try again.")
		return
	}

	cmd := exec.Command("cat", "/dev/null", ">", filepath)
	_, err = cmd.Output()

	if err != nil {
		log.Println("Empty file failed. Please try again")
		return
	}
}
