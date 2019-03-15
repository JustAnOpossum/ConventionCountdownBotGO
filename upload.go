package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/pkg/errors"
)

var zipFilePath string
var creditName string
var creditURL string

func askQuestions() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Bot Data Directory: ")
	scanner.Scan()
	dataDir = scanner.Text()
	imgDir = dataDir + "/img"

	fmt.Print("Path to ZIP File: ")
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
	for i, file := range zipFile.File {
		err = processZipFile(file)
		if err != nil {
			return errors.Wrap(err, "Processing Zip Img")
		}
		fmt.Println("Added Photo #" + strconv.Itoa(i))
	}
	return nil
}

func processZipFile(file *zip.File) error {
	openFile, err := file.Open()
	if err != nil {
		return errors.Wrap(err, "Open Photo From Zip")
	}
	readFile, err := ioutil.ReadAll(openFile)
	if err != nil {
		return errors.Wrap(err, "Reading Zip File")
	}
	fileType := http.DetectContentType(readFile)
	if fileType != "image/png" && fileType != "image/jpeg" {
		return nil
	}
	err = ioutil.WriteFile(path.Join(imgDir, file.Name), readFile, 0664)
	if err != nil {
		return errors.Wrap(err, "Writing File")
	}

	itemToInsert := photo{
		Photo: file.Name,
		Used:  false,
		Name:  creditName,
		URL:   creditURL,
	}
	db.insert("photos", itemToInsert)

	return nil
}
