package logic

import (
	"KnowledgeAcquisition/model"
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

func SearchByImageLogic(imagePath string) (string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
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
	w.Close()

	req, err := http.NewRequest("POST", model.PYTHON_SERVER_URL+"/image_to_keywords", &b)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	log.Debug("FormDataContentType:", w.FormDataContentType())
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Error(err.Error())
			return "", err
		}
		return "", errors.New(string(body))
	}

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
