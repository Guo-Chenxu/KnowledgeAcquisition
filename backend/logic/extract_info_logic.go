package logic

import (
	"KnowledgeAcquisition/model"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func ExtractInfo(doc_id string) (model.DocumentAbstract, error) {
	var result model.DocumentAbstract

	doc, ok := idDocMap[doc_id]
	if !ok {
		log.Error("Error getting doc ", doc_id)
	}
	data := map[string]string{"text": doc.Keywords, "language": doc.Lang.String()}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(model.PYTHON_SERVER_URL+"/extract_info", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return model.DocumentAbstract{}, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err.Error())
			return model.DocumentAbstract{}, err
		}
		return model.DocumentAbstract{}, errors.New(string(body))
	}

	log.Debug("Extract info for doc ", doc_id, " entities: ", result.Entities, " hot_words: ", result.HotWords)

	return result, nil
}

func ExtractInfoRegex(doc_id, pattern, word_class string) (model.DocumentExtractRegex, error) {
	var result model.DocumentExtractRegex

	doc, ok := idDocMap[doc_id]
	if !ok {
		log.Error("Error getting doc ", doc_id)
	}
	data := map[string]string{"text": doc.Keywords, "language": doc.Lang.String(), "pattern": pattern, "word_class": word_class}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(model.PYTHON_SERVER_URL+"/extract_info_regex", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return model.DocumentExtractRegex{}, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err.Error())
			return model.DocumentExtractRegex{}, err
		}
		return model.DocumentExtractRegex{}, errors.New(string(body))
	}

	log.Debug("Extract info for doc ", doc_id, " pattern: ", pattern, " word_class: ", word_class, " words: ", result.Words)

	return result, nil
}
