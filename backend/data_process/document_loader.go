package data_process

import (
	. "MariaInfoRetrieval/maria_types"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func LoadDocuments(dir string) ([]Document, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var documents []Document

	for _, file := range files {
		// Ignore files whose filename are not ended with .json extension
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		filename := filepath.Join(dir, file.Name())

		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		var docs []Document
		dec := json.NewDecoder(f)
		if err := dec.Decode(&docs); err != nil {
			return nil, err
		}

		// Set language of documents
		setLang(docs, filename)

		documents = append(documents, docs...)

		f.Close()
	}

	// Convert document from escaped Unicode format to Unicode format
	for id, doc := range documents {
		// Reindex documents to avoid duplicate id
		doc.Id = strconv.Itoa(id)
		// Convert title
		sTitle := strconv.QuoteToASCII(doc.Title)
		doc.Title = sTitle[1 : len(sTitle)-1]
		// Convert content
		sUnicode := strconv.QuoteToASCII(doc.Content)
		doc.Content = sUnicode[1 : len(sUnicode)-1]
	}

	return documents, nil
}

func setLang(docs []Document, filename string) {
	var lang Language
	if filename == "oiwiki.json" {
		lang = Chinese
	} else {
		lang = English
	}

	for _, doc := range docs {
		doc.Lang = lang
	}
}
