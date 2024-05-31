package util

import (
	. "KnowledgeAcquisition/model"
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

		setLang(docs, filename)

		documents = append(documents, docs...)

		f.Close()
	}

	for id, doc := range documents {
		doc.Id = strconv.Itoa(id)
		sTitle := strconv.QuoteToASCII(doc.Title)
		doc.Title = sTitle[1 : len(sTitle)-1]
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
