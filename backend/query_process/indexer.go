package query_process

import (
	. "MariaInfoRetrieval/maria_types"
	"errors"
	"math"
	"sort"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type DocumentVector struct {
	Doc    Document
	Vector map[string]float64
}

var index = make(map[string][]DocumentVector)
var idDocMap = make(map[string]Document)
var docs []Document
var idfMap = make(map[string]float64) // To hold IDF values of all terms
var epsilon = 1e-10

// Words to ignore
var stopWords = map[string]struct{}{
	"_": {}, ",": {}, ".": {},
}

func isStopWord(word string) bool {
	_, ok := stopWords[word]
	return ok
}

// Build the index and calculate the IDF for all terms
func BuildIndex(documents []Document) {
	// Store documents
	docs = documents

	totalDocs := float64(len(documents))
	docIndex := make(map[string][]Document)

	log.Info("Number of documents loaded:", totalDocs)

	// First, build the index, and doc's summary at the same time
	for _, doc := range documents {
		// Build the idDocMap
		idDocMap[doc.Id] = doc

		// Build the index
		words := WordSplit(doc.Keywords)
		for _, word := range words {
			if !isStopWord(word) {
				// To lower case
				word = strings.ToLower(word)
				docIndex[word] = append(docIndex[word], doc)
			}
		}
	}

	// Second, calculate the IDF values for all words
	for word := range docIndex {
		if _, ok := idfMap[word]; !ok {
			word = strings.ToLower(word)
			idfMap[word] = math.Log(totalDocs / float64(len(docIndex[word])))
		}
	}

	// Third, build index of []DocumentVector
	for _, doc := range documents {
		docVector := buildDocumentVector(doc)
		words := WordSplit(doc.Keywords)

		for _, word := range words {
			word = strings.ToLower(word)
			index[word] = append(index[word], docVector)
		}
	}

}

func buildDocumentVector(doc Document) DocumentVector {
	vector := make(map[string]float64)
	words := WordSplit(doc.Keywords)
	wordCount := float64(len(words))

	// Calculate TF for each term in the document
	for _, word := range words {
		word = strings.ToLower(word)
		vector[word] += 1.0 / wordCount
	}

	// Multiply TF with IDF for each term to get TF-IDF score and normalize the
	// vector
	magnitude := 0.0
	for word, tf := range vector {
		word = strings.ToLower(word)
		tfIdf := tf * idfMap[word]
		vector[word] = tfIdf
		magnitude += tfIdf * tfIdf
	}

	if magnitude > 0.0 {

		sqrtMagnitude := math.Sqrt(magnitude + epsilon)

		// Divide each term's TF-IDF score with the magnitude to get the unit vector
		for word := range vector {
			word = strings.ToLower(word)
			vector[word] /= sqrtMagnitude
		}
	}

	return DocumentVector{Doc: doc, Vector: vector}
}

func buildSummaryDocument(doc Document) SummaryDocument {
	summaryDoc := SummaryDocument{
		Id:      doc.Id,
		Title:   doc.Title,
		URL:     doc.URL,
		Date:    doc.Date,
		Content: calculateSummary(doc.Keywords),
	}
	return summaryDoc
}

func SearchIndex(queryWords []string, page, resultsPerPage int) ([]SearchResult, error) {
	if len(queryWords) == 0 {
		return nil, errors.New("empty query")
	}

	queryVector := buildQueryVector(queryWords)
	log.Info("queryVector:", queryVector)

	// Handle the situation when the query words exist in all or none documents
	magnigude := 0.0
	for _, tfidf := range queryVector {
		magnigude += tfidf * tfidf
	}
	if magnigude == 0 {
		log.Info("Query made up of words in every or no documents. Returning all documents.")
		results := make([]SearchResult, 0, len(docs))
		for _, doc := range docs {
			results = append(results, SearchResult{Doc: buildSummaryDocument(doc), Score: 1.0})
		}

		return results, nil
	}

	vectorCounts := make(map[string]int)
	// Count the total number of vectors for each document ID across all query words
	for _, word := range queryWords {
		word = strings.ToLower(word)
		if vectors, ok := index[word]; ok {
			for _, vector := range vectors {
				vectorCounts[vector.Doc.Id]++
			}
		}
	}

	// Store the count of query words in each document. Let documents that
	// include more query words get a higher score
	queryWordCounts := make(map[string]int)
	// Count of query words in each doc's title
	titleQueryWordCounts := make(map[string]int)

	// Mutex to protect concurrent writes to queryWordCounts
	var mutex sync.Mutex

	// Parallel result computation
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

					// Ignore upper case
					wi := strings.ToLower(w)

					// Calculate score
					score := cosineSimilarity(queryVector, v.Vector)

					// Adjust the score based on:
					// -  frequency of query words
					// - document length
					// - position of first query word
					frequency := float64(len(WordSplit(v.Doc.Keywords)))
					position := float64(strings.Index(v.Doc.Keywords, wi))
					length := float64(len(v.Doc.Keywords))

					adjustment := (1 + math.Log(frequency+1)) * (1 / (1 + math.Log(length+1)) * (1 / (1 + math.Log(position+1))))
					score *= adjustment

					// Increase the count of query words in this document
					if strings.Contains(v.Doc.Keywords, wi) || strings.Contains(strings.ToLower(v.Doc.Title), wi) {
						// Guard the write with the mutex
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

	// Wait for all goroutines to finish, then close the results channel
	go func() {
		wg.Wait()
		for _, scoresChan := range scoresChansMap {
			close(scoresChan)
		}
	}()

	// Collect the results
	scoreMap := make(map[string]*SearchResult)
	for id, scoresChan := range scoresChansMap {
		totalScore := 0.0
		for score := range scoresChan {
			totalScore += score
		}
		// Boost the score of the document based on the number of query words it contains
		totalScore *= float64(1 + queryWordCounts[id])

		// Boost the score based on the number of query words in title
		totalScore *= 1.2 * float64(1+titleQueryWordCounts[id])

		// Build document summary
		summaryDoc := buildSummaryDocument(idDocMap[id])
		scoreMap[id] = &SearchResult{Doc: summaryDoc, Score: totalScore}
	}

	log.Info(len(scoreMap), " results")
	log.Debug(">>> scoreMap")
	for k, v := range scoreMap {
		log.Debug(k, ":", "Doc:", v.Doc, "Score:", v.Score)
	}
	log.Debug("<<< scoreMap")

	// Convert the scoreMap to a slice
	results := make([]SearchResult, 0, len(scoreMap))
	for _, result := range scoreMap {
		results = append(results, *result)
	}

	// Sort results by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Apply pagination
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

func GetFullDoc(id string) (Document, bool) {
	doc, ok := idDocMap[id]
	return doc, ok
}

func buildQueryVector(queryWords []string) map[string]float64 {
	vector := make(map[string]float64)
	wordCount := float64(len(queryWords))

	// Calculate TF for each term in the query
	for _, word := range queryWords {
		word = strings.ToLower(word)
		vector[word] += 1.0 / wordCount
	}

	// Multiply TF with IDF for each term to get TF-IDF score and normalize the vector
	magnitude := 0.0
	for word, tf := range vector {
		word = strings.ToLower(word)
		idf, ok := idfMap[word]
		if !ok {
			// Skip non-indexed words
			continue
		}
		tfIdf := idf * tf
		vector[word] = tfIdf
		magnitude += tfIdf * tfIdf
	}

	if magnitude > 0.0 {
		sqrtMagnitude := math.Sqrt(magnitude + epsilon)

		// Divide each term's TF-IDF score with the magnitude to get the unit vector
		// Only if magnitude is non-zero
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
