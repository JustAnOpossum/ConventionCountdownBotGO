//photo.go
//Creates the actual countdown image. All image generation steps are acomplished here.

package main

import (
	"image"
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
	FilePath   string
	FileName   string
	CreditName string
	CreditURL  string
	DaysLeft   int
}

//Creates the image
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
	//Set's the font color correctly
	R, G, B := imgColors.Color.RGB()
	imgCtx.SetRGB(float64(R), float64(G), float64(B))
	//Loads the font face
	err = imgCtx.LoadFontFace(path.Join(dataDir, config.ImgSend.Font), config.ImgSend.FontSize)
	if err != nil {
		log.Println("Setting font face in createImg() " + err.Error())
		os.Exit(0)
	}
	//Draws the text to the image
	imgCtx.DrawStringAnchored(strconv.Itoa(getDays(config.Date)), float64(imgCtx.Image().Bounds().Dx()/2), float64(imgCtx.Image().Bounds().Dy()/2), 0.5, 0.5)
	//Saves the image file
	err = imgCtx.SavePNG(filePath)
	if err != nil {
		log.Println("Saving PNG in createImg() " + err.Error())
		os.Exit(0)
	}
	imgToReturn.FilePath = filePath

	//Now that the photo is successfuly created set that photo to used
	photos.update(bson.M{"photo": img.Photo}, bson.M{"$set": bson.M{"used": true}})
	return imgToReturn
}

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
	photos.findAll(bson.M{"name": nameToSearch}, &nameToFind)

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
