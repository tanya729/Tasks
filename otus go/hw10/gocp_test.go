package main

import (
	"log"
	"os"
	"testing"
)

func createTestFile() {
	size := int64(10 * 1024 * 1024)
	fd, err := os.Create("testFile")
	if err != nil {
		log.Fatal("Failed to create output")
	}
	_, err = fd.Seek(size-1, 0)
	if err != nil {
		log.Fatal("Failed to seek")
	}
	_, err = fd.Write([]byte{0})
	if err != nil {
		log.Fatal("Write failed")
	}
	err = fd.Close()
	if err != nil {
		log.Fatal("Failed to close file")
	}
}

func removeTestFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		log.Fatalf("Failed to remove file: %s", err.Error())
	}
}

func getFileInfo(fileName string) os.FileInfo {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatalf("Something wrong: %s", err.Error())
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Something wrong: %s", err.Error())
	}
	return fileInfo
}

func TestCopyEqual(t *testing.T) {
	createTestFile()
	defer removeTestFile("testFile")
	err := Copy("testFile", "testOutputFile", -1, 0)
	if err != nil {
		log.Fatalf("Something wrong: %s", err.Error())
	}
	defer removeTestFile("testOutputFile")
	fileInputInfo := getFileInfo("testFile")
	fileOutputInfo := getFileInfo("testOutputFile")
	if fileInputInfo.Size() != fileOutputInfo.Size() {
		t.Error("Size not equal.")
	}
}

func TestCopyEmpty(t *testing.T) {
	createTestFile()
	defer removeTestFile("testFile")
	err := Copy("testFile", "testOutputFile", 0, 0)
	if err != nil {
		log.Fatalf("Something wrong: %s", err.Error())
	}
	defer removeTestFile("testOutputFile")
	fileOutputInfo := getFileInfo("testOutputFile")
	if 0 != fileOutputInfo.Size() {
		t.Error("Expected zero size file.")
	}
}

func TestCopyLimit(t *testing.T) {
	limit := int64(5 * 1024 * 1024)
	createTestFile()
	defer removeTestFile("testFile")
	err := Copy("testFile", "testOutputFile", limit, 0)
	if err != nil {
		log.Fatalf("Something wrong: %s", err.Error())
	}
	defer removeTestFile("testOutputFile")
	fileOutputInfo := getFileInfo("testOutputFile")
	if limit != fileOutputInfo.Size() {
		t.Error("Expected 5mb size file.")
	}
}

func TestCopyOffset(t *testing.T) {
	offset := int64(6 * 1024 * 1024)
	createTestFile()
	defer removeTestFile("testFile")
	err := Copy("testFile", "testOutputFile", -1, offset)
	if err != nil {
		log.Fatalf("Something wrong: %s", err.Error())
	}
	defer removeTestFile("testOutputFile")
	fileOutputInfo := getFileInfo("testOutputFile")
	if offset <= fileOutputInfo.Size() {
		t.Error("Expected size <= 6mb.")
	}
}
