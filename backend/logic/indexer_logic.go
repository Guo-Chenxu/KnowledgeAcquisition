package logic

import (
	"KnowledgeAcquisition/model"
	"errors"
	"math"
	"sort"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type DocumentVector struct {
	Doc    model.Document
	Vector map[string]float64
}

var index = make(map[string][]DocumentVector)
var idDocMap = make(map[string]model.Document)
var docs []model.Document
var idfMap = make(map[string]float64)
var epsilon = 1e-10

var stopWordMap = map[string]struct{}{
	"_": {}, ",": {}, ".": {},
}

func isStopWord(word string) bool {
	_, ok := stopWordMap[word]
	return ok
}

// BuildIndex 创建倒排索引并计算idf
func BuildIndex(documents []model.Document) {
	docs = documents

	totalDocs := float64(len(documents))
	docIndex := make(map[string][]model.Document)

	log.Info("Number of documents loaded:", totalDocs)

	// 建立索引
	for _, doc := range documents {
		idDocMap[doc.Id] = doc

		words := WordSplit(doc.Keywords)
		for _, word := range words {
			if !isStopWord(word) {
				word = strings.ToLower(word)
				docIndex[word] = append(docIndex[word], doc)
			}
		}
	}

	// 计算idf
	for word := range docIndex {
		if _, ok := idfMap[word]; !ok {
			word = strings.ToLower(word)
			idfMap[word] = math.Log(totalDocs / float64(len(docIndex[word])))
		}
	}

	// 根据tf-idf计算文档向量，创建tf-idf索引
	for _, doc := range documents {
		docVector := buildDocumentVector(doc)
		words := WordSplit(doc.Keywords)

		for _, word := range words {
			word = strings.ToLower(word)
			index[word] = append(index[word], docVector)
		}
	}

}

func buildDocumentVector(doc model.Document) DocumentVector {
	vector := make(map[string]float64)
	words := WordSplit(doc.Keywords)
	wordCount := float64(len(words))

	// 计算tf
	for _, word := range words {
		word = strings.ToLower(word)
		vector[word] += 1.0 / wordCount
	}

	// 计算tf-idf
	magnitude := 0.0
	for word, tf := range vector {
		word = strings.ToLower(word)
		tfIdf := tf * idfMap[word]
		vector[word] = tfIdf
		magnitude += tfIdf * tfIdf
	}

	// 归一化
	if magnitude > 0.0 {
		sqrtMagnitude := math.Sqrt(magnitude + epsilon)
		for word := range vector {
			word = strings.ToLower(word)
			vector[word] /= sqrtMagnitude
		}
	}

	return DocumentVector{Doc: doc, Vector: vector}
}

func buildSummaryDocument(doc model.Document) model.SummaryDocument {
	summaryDoc := model.SummaryDocument{
		Id:      doc.Id,
		Title:   doc.Title,
		URL:     doc.URL,
		Date:    doc.Date,
		Content: calculateSummary(doc.Keywords),
	}
	return summaryDoc
}

func SearchIndex(queryWords []string, page, resultsPerPage int) ([]model.SearchResult, error) {
	if len(queryWords) == 0 {
		return nil, errors.New("empty query")
	}

	queryVector := buildQueryVector(queryWords)
	log.Info("queryVector:", queryVector)

	magnigude := 0.0
	for _, tfidf := range queryVector {
		magnigude += tfidf * tfidf
	}
	if magnigude == 0 {
		log.Info("Query made up of words in every or no documents. Returning all documents.")
		results := make([]model.SearchResult, 0, len(docs))
		for _, doc := range docs {
			results = append(results, model.SearchResult{Doc: buildSummaryDocument(doc), Score: 1.0})
		}

		return results, nil
	}

	vectorCounts := make(map[string]int)
	for _, word := range queryWords {
		word = strings.ToLower(word)
		if vectors, ok := index[word]; ok {
			for _, vector := range vectors {
				vectorCounts[vector.Doc.Id]++
			}
		}
	}

	queryWordCounts := make(map[string]int)
	titleQueryWordCounts := make(map[string]int)

	var mutex sync.Mutex

	scoresChansMap := make(map[string]chan float64)
	for id, count := range vectorCounts {
		scoresChansMap[id] = make(chan float64, count)
	}

	var wg sync.WaitGroup

	for _, word := range queryWords {
		word = strings.ToLower(word)
		if vectors, ok := index[word]; ok {
			for _, vector := range vectors {

				wg.Add(1)
				go func(w string, v DocumentVector, scoresChan chan float64) {
					defer wg.Done()

					wi := strings.ToLower(w)

					score := cosineSimilarity(queryVector, v.Vector)

					frequency := float64(len(WordSplit(v.Doc.Keywords)))
					position := float64(strings.Index(v.Doc.Keywords, wi))
					length := float64(len(v.Doc.Keywords))

					adjustment := (1 + math.Log(frequency+1)) * (1 / (1 + math.Log(length+1)) * (1 / (1 + math.Log(position+1))))
					score *= adjustment

					if strings.Contains(v.Doc.Keywords, wi) || strings.Contains(strings.ToLower(v.Doc.Title), wi) {
						mutex.Lock()
						if strings.Contains(v.Doc.Keywords, wi) {
							queryWordCounts[v.Doc.Id]++
						}
						if strings.Contains(strings.ToLower(v.Doc.Title), wi) {
							titleQueryWordCounts[v.Doc.Id]++
						}
						mutex.Unlock()
					}
					scoresChan <- score
				}(word, vector, scoresChansMap[vector.Doc.Id])
			}
		}
	}

	go func() {
		wg.Wait()
		for _, scoresChan := range scoresChansMap {
			close(scoresChan)
		}
	}()

	scoreMap := make(map[string]*model.SearchResult)
	for id, scoresChan := range scoresChansMap {
		totalScore := 0.0
		for score := range scoresChan {
			totalScore += score
		}
		totalScore *= float64(1 + queryWordCounts[id])

		totalScore *= 1.2 * float64(1+titleQueryWordCounts[id])

		summaryDoc := buildSummaryDocument(idDocMap[id])
		scoreMap[id] = &model.SearchResult{Doc: summaryDoc, Score: totalScore}
	}

	log.Info(len(scoreMap), " results")
	log.Debug(">>> scoreMap")
	for k, v := range scoreMap {
		log.Debug(k, ":", "Doc:", v.Doc, "Score:", v.Score)
	}
	log.Debug("<<< scoreMap")

	results := make([]model.SearchResult, 0, len(scoreMap))
	for _, result := range scoreMap {
		results = append(results, *result)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	start := (page - 1) * resultsPerPage
	end := start + resultsPerPage
	if start > len(results) {
		start = len(results)
	}
	if end > len(results) {
		end = len(results)
	}

	results = results[start:end]

	return results, nil
}

func GetFullDoc(id string) (model.Document, bool) {
	doc, ok := idDocMap[id]
	return doc, ok
}

func buildQueryVector(queryWords []string) map[string]float64 {
	vector := make(map[string]float64)
	wordCount := float64(len(queryWords))

	for _, word := range queryWords {
		word = strings.ToLower(word)
		vector[word] += 1.0 / wordCount
	}

	magnitude := 0.0
	for word, tf := range vector {
		word = strings.ToLower(word)
		idf, ok := idfMap[word]
		if !ok {
			continue
		}
		tfIdf := idf * tf
		vector[word] = tfIdf
		magnitude += tfIdf * tfIdf
	}

	if magnitude > 0.0 {
		sqrtMagnitude := math.Sqrt(magnitude + epsilon)
		for word := range vector {
			vector[word] /= sqrtMagnitude
		}
	}

	return vector
}

func cosineSimilarity(vector1, vector2 map[string]float64) float64 {
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0
	for word, value := range vector1 {
		word = strings.ToLower(word)
		dotProduct += value * vector2[word]
		magnitude1 += value * value
	}
	for _, value := range vector2 {
		magnitude2 += value * value
	}

	sqrtEpsMag1 := math.Sqrt(magnitude1 + epsilon)
	sqrtEpsMag2 := math.Sqrt(magnitude2 + epsilon)
	return dotProduct / (sqrtEpsMag1 * sqrtEpsMag2)
}

func calculateSummary(content string) string {
	if len(content) > 100 {
		return content[:100] + "..."
	}
	return content
}
