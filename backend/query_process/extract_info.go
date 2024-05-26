package query_process

import (
	. "MariaInfoRetrieval/maria_types"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func ExtractInfo(doc_id string) (DocumentAbstract, error) {
	var result DocumentAbstract

	doc, ok := idDocMap[doc_id]
	if !ok {
		log.Error("Error getting doc ", doc_id)
	}
	data := map[string]string{"text": doc.Keywords, "language": doc.Lang.String()}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(python_server_url+"/extract_info", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return DocumentAbstract{}, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&result)

	// Check the response
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err.Error())
			return DocumentAbstract{}, err
		}
		return DocumentAbstract{}, errors.New(string(body))
	}

	log.Debug("Extract info for doc ", doc_id, " entities: ", result.Entities, " hot_words: ", result.HotWords)

	return result, nil
}
