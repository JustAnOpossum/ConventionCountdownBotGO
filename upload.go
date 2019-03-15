package main

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"os"
)

var zipFilePath string
var creditName string
var creditURL string

func askQuestions() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("What is the Data Dir: ")
	scanner.Scan()
	dataDir = scanner.Text()

	fmt.Print("ZIP File Path to Upload: ")
	scanner.Scan()
	zipFilePath = scanner.Text()

	fmt.Print("Who to Credit for Pictures: ")
	scanner.Scan()
	creditName = scanner.Text()

	fmt.Print("URL for Credit: ")
	scanner.Scan()
	creditURL = scanner.Text()
}

func uploadZip() error {
	zipFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return errors.New("Cannot Open Zip file")
	}
	defer zipFile.Close()
	for _, file := range zipFile.File {
		processZipFile(file)
	}
	return nil
}

func processZipFile(file *zip.File) {
	fileName := file.Name
}
