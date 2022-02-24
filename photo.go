package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/fogleman/gg"
	"github.com/generaltso/vibrant"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

//Returned image to be parsed by twitter and the telegram upload functions
type finalImg struct {
	ImgReader  *bytes.Reader
	FilePath   string
	FileName   string
	CreditName string
	CreditURL  string
	DaysLeft   int
}

//Creates
func createImg() finalImg {
	//Sets the amount of days left for the final return image
	imgToReturn := finalImg{
		DaysLeft: getDays(config.Date),
	}

	//Finds a sutible image and appends the credit
	img := findImg()
	imgToReturn.CreditName = img.Name
	imgToReturn.CreditURL = img.URL
	decodedImg, err := loadImg(img)
	if err != nil {
		log.Println("Decoding image in createImg() " + err.Error())
		os.Exit(0)
	}

	//Gets the colors of the image
	imgColors, err := getImgColors(decodedImg)
	if err != nil {
		log.Println("Getting Color Pallete " + err.Error())
		os.Exit(0)
	}

	//Creates image context
	imgCtx := gg.NewContextForImage(decodedImg)
	fileName := (strconv.Itoa(getDays(config.Date))) + ".jpg"
	filePath := path.Join(dataDir, "countdown/"+fileName)
	err = imgCtx.LoadFontFace(path.Join(dataDir, config.ImgSend.Font), config.ImgSend.FontSize)
	if err != nil {
		log.Println("Setting font face in createImg() " + err.Error())
		os.Exit(0)
	}
	//Converts the vibrant color to RGB and feeds that into gg
	R, G, B := imgColors.Color.RGB()
	newColor := color.RGBA{R: uint8(R), G: uint8(G), B: uint8(B)}
	imgCtx.SetColor(newColor)
	imgCtx.DrawStringAnchored(strconv.Itoa(getDays(config.Date)), float64(decodedImg.Bounds().Dx()/2), float64(decodedImg.Bounds().Dy()/2), 0.5, 0.5)
	err = imgCtx.SavePNG(filePath)
	if err != nil {
		log.Println("Saving PNG in createImg() " + err.Error())
		os.Exit(0)
	}
	imgBuffer := make([]byte, 0)
	writer := bytes.NewBuffer(imgBuffer)
	err = imgCtx.EncodePNG(writer)
	if err != nil {
		log.Println("Encoding PNG in createImg() " + err.Error())
		os.Exit(0)
	}
	reader := bytes.NewReader(imgBuffer)
	imgToReturn.ImgReader = reader

	//Now that the photo is successfuly created set that photo to used
	photos.update(bson.M{"photo": img.Photo}, bson.M{"$set": bson.M{"used": true}})
	return imgToReturn

	// imgWand := imagick.NewMagickWand()
	// defer imgWand.Destroy()
	// drawImgText(imgWand, img, imgColors)
	// imgWand.SetImageFormat("JPEG")
	// imgToReturn.ImgReader = bytes.NewReader(imgWand.GetImageBlob())
	// imgWand.WriteImage(filePath)
	// imgToReturn.FilePath = filePath

	// txtToAppend := []byte("\nfile 'countdown/" + fileName + "'\nduration 0.5")
	// slideshowFile, err := os.OpenFile(path.Join(dataDir, "slideshow.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	return finalImg{}, errors.Wrap(err, "Slideshow File")
	// }
	// if _, err := slideshowFile.Write(txtToAppend); err != nil {
	// 	return finalImg{}, errors.Wrap(err, "Writing to Sliideshow File")
	// }
	// if err := slideshowFile.Close(); err != nil {
	// 	return finalImg{}, errors.Wrap(err, "Closing Slideshow File")
	// }
}

// func drawImgText(imgWand *imagick.MagickWand, img photo, colors *vibrant.Swatch) {
// 	loadedImg, _ := os.Open(path.Join(imgDir, img.Photo))
// 	defer loadedImg.Close()

// 	textWand := imagick.NewDrawingWand()
// 	textColor := imagick.NewPixelWand()
// 	defer textWand.Destroy()
// 	defer textColor.Destroy()
// 	imgWand.ReadImageFile(loadedImg)

// 	textColor.SetColor(colors.Color.RGBHex())
// 	textWand.SetFont(path.Join(dataDir, config.ImgSend.Font))
// 	textWand.SetFontSize(config.ImgSend.FontSize)
// 	textWand.SetFillColor(textColor)
// 	textWand.SetGravity(imagick.GRAVITY_SOUTH_WEST)
// 	textWand.Annotation(0, 0, strconv.Itoa(getDays(config.Date)))

// 	imgWand.DrawImage(textWand)
// }

//Searches the database for a photo that has not been used yet
func findImg() photo {
	//Name is used for the distinct value so that more different photos will be shown each day
	items := photos.distinct(bson.M{"used": false}, "name")
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	//Checks if the search is empty becuase all images have been used, or empty because no photos have been added
	if len(items) == 0 {
		items = photos.distinct(bson.M{}, "name")
		if len(items) == 0 {
			log.Println("Error: No photos found in database, add some photo to send images.")
			os.Exit(0)
		}
	}

	var nameToFind []photo
	nameToSearch := items[random.Intn(len(items))]
	nameToFind = photos.findAll(bson.M{"name": nameToSearch})

	return nameToFind[random.Intn(len(nameToFind))]
}

//Generates a color pallete from the selected image
func getImgColors(decodedImg image.Image) (*vibrant.Swatch, error) {
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

//Loads and image from a path and returns a decoded image
func loadImg(img photo) (image.Image, error) {
	loadedImg, err := os.Open(path.Join(imgDir, img.Photo))
	if err != nil {
		return nil, errors.Wrap(err, "Loading Img File Pallete")
	}
	defer loadedImg.Close()
	decodedImg, _, err := image.Decode(loadedImg)
	if err != nil {
		return nil, errors.Wrap(err, "Error Decoding Img")
	}
	return decodedImg, nil
}
