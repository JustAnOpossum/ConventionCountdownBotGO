package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/nfnt/resize"

	"github.com/pkg/errors"
)

var zipFilePath string
var creditName string
var creditURL string
var clearOrNot string

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

	fmt.Print("Clear database and folders (y/n): ")
	scanner.Scan()
	clearOrNot = scanner.Text()
}

func uploadZip() error {
	//Clears database and folder if user wants to reset
	if clearOrNot == "y" {
		photos.removeAll()
		os.RemoveAll(imgDir)
	}
	//If the folder does not exist, then create the image holding area
	if _, err := os.Stat(imgDir); os.IsNotExist(err) {
		os.Mkdir(imgDir, 0644)
	}

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
	if file.UncompressedSize64 > 5000000 {
		fmt.Fprintln(out, "Called Reszie Img")
		err = resizeImg(&readFile, bytes.NewReader(readFile))
		if err != nil {
			return errors.Wrap(err, "Resizing Img")
		}
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
	photos.insert(itemToInsert)

	return nil
}

func resizeImg(outputImg *[]byte, inputReader *bytes.Reader) error {
	tempImg, _, err := image.Decode(inputReader)
	if err != nil {
		return errors.Wrap(err, "Decoding Img")
	}
	resizedImg := resize.Resize(1000, 0, tempImg, resize.Lanczos3)
	tempBuffer := new(bytes.Buffer)
	err = jpeg.Encode(tempBuffer, resizedImg, nil)
	if err != nil {
		return errors.Wrap(err, "Encoding Jpeg")
	}
	tempBytes := tempBuffer.Bytes()
	*outputImg = tempBytes
	return nil
}
