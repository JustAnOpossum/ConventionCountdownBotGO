package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/generaltso/vibrant"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func createImg() error {
	img := findImg()
	imgColors, err := getImgColors(img)
	if err != nil {
		return errors.Wrap(err, "Getting Color Pallete")
	}
	drawImgText(img, imgColors)
	return nil
}

func drawImgText(img photo, colors *vibrant.Swatch) {
	loadedImg, _ := os.Open(path.Join(imgDir, img.Photo))
	defer loadedImg.Close()

	imgWand := imagick.NewMagickWand()
	textWand := imagick.NewDrawingWand()
	textColor := imagick.NewPixelWand()
	defer imgWand.Destroy()
	defer textWand.Destroy()
	defer textColor.Destroy()
	imgWand.ReadImageFile(loadedImg)

	width := float64(imgWand.GetImageWidth())
	height := float64(imgWand.GetImageHeight())
	fmt.Println(width/2, height/2)
	textColor.SetColor("red")
	textWand.SetFontSize(150)
	textWand.SetFillColor(textColor)
	textWand.SetGravity(config.ImgSend.GravityMode)
	textWand.Annotation(0, 0, "109")

	imgWand.DrawImage(textWand)
	err := imgWand.WriteImage("test.jpg")
	fmt.Println(err)
}

func findImg() photo {
	items := db.distinct("photos", bson.M{"used": false}, "name")
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	var nameToFind photo
	nameToSearch := items[random.Intn(len(items))]
	db.findOne("photos", bson.M{"used": false, "name": nameToSearch}, &nameToFind)

	return nameToFind
}

func getImgColors(img photo) (*vibrant.Swatch, error) {
	loadedImg, err := os.Open(path.Join(imgDir, img.Photo))
	if err != nil {
		return nil, errors.Wrap(err, "Loading Img File Pallete")
	}
	defer loadedImg.Close()
	decodedImg, _, err := image.Decode(loadedImg)
	if err != nil {
		return nil, errors.Wrap(err, "Error Decoding Img")
	}
	pallete, err := vibrant.NewPaletteFromImage(decodedImg)
	if err != nil {
		return nil, errors.Wrap(err, "Error Creating Color Pallete")
	}
	swatches := pallete.ExtractAwesome()
	if swatches["Vibrant"] == nil {
		for _, swatch := range swatches {
			return swatch, nil
		}
	}
	return swatches["Vibrant"], nil
}
