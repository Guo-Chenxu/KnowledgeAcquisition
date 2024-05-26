package query_process

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type KeywordResponse struct {
	Keyword string `json:"keyword"`
}

func GetKeywordsFromImage(imagePath string) (string, error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file
	f, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fw, err := w.CreateFormFile("file", filepath.Base(imagePath))
	log.Debug("file: ", filepath.Base(imagePath))
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return "", err
	}
	// Close the multipart writer so that the form data can be sent
	w.Close()

	// Now that you have a form, you can submit it to your handler
	req, err := http.NewRequest("POST", python_server_url+"/image_to_keywords", &b)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	// Set content type, contain the boundary
	log.Debug("FormDataContentType:", w.FormDataContentType())
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Error(err.Error())
			return "", err
		}
		return "", errors.New(string(body))
	}

	// Parse the JSON response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var kr KeywordResponse
	err = json.Unmarshal(body, &kr)
	if err != nil {
		return "", err
	}

	return kr.Keyword, nil
}
