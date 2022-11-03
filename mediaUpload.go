//mediaUpload.go
//Handles uploading an image to twitter so that it can be attached to a tweet

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/dghubble/oauth1"
	"github.com/pkg/errors"
)

type twitterMedia struct {
	MediaID          int64  `json:"media_id"`
	MediaIDString    string `json:"media_id_string"`
	Size             int    `json:"size"`
	ExpiresAfterSecs int    `json:"expires_after_secs"`
	Image            struct {
		ImageType string `json:"image_type"`
		W         int    `json:"w"`
		H         int    `json:"h"`
	} `json:"image"`
}

func uploadTwitterImg(imgPath string) (int64, error) {
	twitterConfig := oauth1.NewConfig(config.Twitter.ConsumerKey, config.Twitter.ConsumerSecret)
	twitterToken := oauth1.NewToken(config.Twitter.AccessToken, config.Twitter.AccessSecret)
	client := twitterConfig.Client(oauth1.NoContext, twitterToken)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("media", "dailySendImg")
	img, err := os.Open(imgPath)
	if err != nil {
		return 0, err
	}
	io.Copy(part, img)
	writer.Close()

	req, _ := http.NewRequest("POST", "https://upload.twitter.com/1.1/media/upload.json", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	resBody, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return 0, errors.New("HTTP Error got " + strconv.Itoa(res.StatusCode) + " " + string(resBody))
	}
	req.Body.Close()
	var mediaInfo twitterMedia
	json.Unmarshal(resBody, &mediaInfo)

	err = img.Close()
	if err != nil {
		return 0, errors.New("Closing image")
	}

	return mediaInfo.MediaID, nil
}

//Depricated for simpiler media upload option. No more video upload

// func uploadTwitterMedia(path string, mediaType string) (int64, error) {
// 	twitterConfig := oauth1.NewConfig(config.Twitter.ConsumerKey, config.Twitter.ConsumerSecret)
// 	twitterToken := oauth1.NewToken(config.Twitter.AccessToken, config.Twitter.AccessSecret)
// 	httpClient := twitterConfig.Client(oauth1.NoContext, twitterToken)

// 	mediaFile, err := os.ReadFile(path)
// 	if err != nil {
// 		return 0, errors.Wrap(err, "Loading Media File")
// 	}
// 	mediaID, err := initMediaUpload(strconv.Itoa(len(mediaFile)), mediaType, httpClient)
// 	if err != nil {
// 		return 0, errors.Wrap(err, "Init Media Upload")
// 	}
// 	fmt.Fprintln(out, "Got Media ID")
// 	err = appendMediaUpload(mediaFile, mediaID, httpClient)
// 	if err != nil {
// 		return 0, errors.Wrap(err, "Appending Chunks")
// 	}
// 	fmt.Fprintln(out, "Done Appending Chunks")
// 	err = finalMediaUpload(mediaID, httpClient)
// 	if err != nil {
// 		return 0, errors.Wrap(err, "Final Media Error")
// 	}
// 	fmt.Fprintln(out, "Done Finalizing Media")
// 	return mediaID, nil
// }

// func initMediaUpload(mediaSize string, mediaType string, client *http.Client) (int64, error) {
// 	form := url.Values{}
// 	form.Add("command", "INIT")
// 	form.Add("media_type", mediaType)
// 	form.Add("total_bytes", mediaSize)

// 	req, err := http.NewRequest("POST", "https://upload.twitter.com/1.1/media/upload.json", strings.NewReader(form.Encode()))
// 	if err != nil {
// 		return 0, errors.Wrap(err, "Creating Request")
// 	}
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	res, err := client.Do(req)
// 	if err != nil {
// 		return 0, errors.Wrap(err, "Making Request")
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode >= 400 {
// 		return 0, errors.New("Got Wrong HTTP Code in init " + strconv.Itoa(res.StatusCode))
// 	}

// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return 0, errors.Wrap(err, "Reading Body")
// 	}

// 	var mediaRes twitterMedia
// 	err = json.Unmarshal(body, &mediaRes)
// 	if err != nil {
// 		return 0, errors.Wrap(err, "Unmarshall JSON")
// 	}

// 	return mediaRes.MediaID, nil
// }

// func appendMediaUpload(mediaFile []byte, mediaID int64, client *http.Client) error {
// 	chunkSize := 1024 * 1024
// 	for seg := 0; seg*chunkSize < len(mediaFile); seg++ {
// 		start := seg * chunkSize
// 		end := (seg + 1) * chunkSize
// 		if end >= len(mediaFile) {
// 			end = len(mediaFile)
// 		}

// 		var body bytes.Buffer
// 		form := multipart.NewWriter(&body)
// 		form.WriteField("command", "APPEND")
// 		form.WriteField("media_id", fmt.Sprint(mediaID))
// 		form.WriteField("segment_index", strconv.Itoa(seg))
// 		formFile, _ := form.CreateFormFile("media", "countdownVideo.mp4")
// 		_, err := formFile.Write(mediaFile[start:end])
// 		if err != nil {
// 			return errors.Wrap(err, "Writing Chunk to Form")
// 		}

// 		form.Close()

// 		req, err := http.NewRequest("POST", "https://upload.twitter.com/1.1/media/upload.json", &body)
// 		if err != nil {
// 			return errors.Wrap(err, "Creating Append Request")
// 		}
// 		req.Header.Add("Content-Type", form.FormDataContentType())

// 		res, err := client.Do(req)
// 		if err != nil {
// 			return errors.Wrap(err, "Making Append Request")
// 		}
// 		defer res.Body.Close()
// 		fmt.Fprintln(out, "Appended To Media")
// 	}
// 	return nil
// }

// func finalMediaUpload(mediaID int64, client *http.Client) error {
// 	form := url.Values{}
// 	form.Add("command", "FINALIZE")
// 	form.Add("media_id", fmt.Sprint(mediaID))

// 	req, err := http.NewRequest("POST", "https://upload.twitter.com/1.1/media/upload.json", strings.NewReader(form.Encode()))
// 	if err != nil {
// 		return errors.Wrap(err, "Creating Request")
// 	}
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	res, err := client.Do(req)
// 	if err != nil {
// 		return errors.Wrap(err, "Making Request")
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode >= 400 {
// 		body, _ := io.ReadAll(res.Body)
// 		return errors.New("Got Wrong HTTP Code in Final " + strconv.Itoa(res.StatusCode) + " " + string(body))
// 	}
// 	return nil
// }
