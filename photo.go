package main

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/generaltso/vibrant"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type finalImg struct {
	ImgReader  *bytes.Reader
	FilePath   string
	FileName   string
	CreditName string
	CreditURL  string
	DaysLeft   int
}

func createImg() (finalImg, error) {
	imgToReturn := finalImg{
		DaysLeft: getDays(config.Date),
	}

	img := findImg()
	imgToReturn.CreditName = img.Name
	imgToReturn.CreditURL = img.URL

	imgColors, err := getImgColors(img)
	if err != nil {
		return finalImg{}, errors.Wrap(err, "Getting Color Pallete")
	}

	imgWand := imagick.NewMagickWand()
	defer imgWand.Destroy()
	drawImgText(imgWand, img, imgColors)
	imgWand.SetImageFormat("JPEG")
	fileName := (strconv.Itoa(getDays(config.Date))) + ".jpg"
	filePath := path.Join(dataDir, "countdown/"+fileName)
	imgToReturn.ImgReader = bytes.NewReader(imgWand.GetImageBlob())
	imgWand.WriteImage(filePath)
	imgToReturn.FilePath = filePath

	txtToAppend := []byte("\nfile 'countdown/" + fileName + "'\nduration 0.5")
	slideshowFile, err := os.OpenFile(path.Join(dataDir, "slideshow.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return finalImg{}, errors.Wrap(err, "Slideshow File")
	}
	if _, err := slideshowFile.Write(txtToAppend); err != nil {
		return finalImg{}, errors.Wrap(err, "Writing to Sliideshow File")
	}
	if err := slideshowFile.Close(); err != nil {
		return finalImg{}, errors.Wrap(err, "Closing Slideshow File")
	}

	photos.update(bson.M{"photo": img.Photo}, bson.M{"$set": bson.M{"used": true}})
	return imgToReturn, nil
}

func drawImgText(imgWand *imagick.MagickWand, img photo, colors *vibrant.Swatch) {
	loadedImg, _ := os.Open(path.Join(imgDir, img.Photo))
	defer loadedImg.Close()

	textWand := imagick.NewDrawingWand()
	textColor := imagick.NewPixelWand()
	defer textWand.Destroy()
	defer textColor.Destroy()
	imgWand.ReadImageFile(loadedImg)

	textColor.SetColor(colors.Color.RGBHex())
	textWand.SetFont(path.Join(dataDir, config.ImgSend.Font))
	textWand.SetFontSize(config.ImgSend.FontSize)
	textWand.SetFillColor(textColor)
	textWand.SetGravity(imagick.GRAVITY_SOUTH_WEST)
	textWand.Annotation(0, 0, strconv.Itoa(getDays(config.Date)))

	imgWand.DrawImage(textWand)
}

func findImg() photo {
	items := photos.distinct(bson.M{"used": false}, "name")
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	var nameToFind []photo
	nameToSearch := items[random.Intn(len(items))]
	photos.findAll(bson.M{"used": false, "name": nameToSearch}, &nameToFind)

	return nameToFind[random.Intn(len(nameToFind))]
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
