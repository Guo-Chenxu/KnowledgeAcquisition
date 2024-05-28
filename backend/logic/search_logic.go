package logic

import (
	"KnowledgeAcquisition/model"
	"errors"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Global search result cache
var cache_capacity = 10
var cache = NewCache(cache_capacity)

func PerformSearch(q string, page string, resultsPerPage string) (r model.SearchResponse, err error) {
	cacheKey := fmt.Sprintf("%s-%s-%s", q, page, resultsPerPage)

	// Return if hit cache
	if cachedResults, found := cache.Get(cacheKey); found {
		return model.SearchResponse{Code: 200, Results: cachedResults}, nil
	}

	// Else search
	intPage, err := strconv.Atoi(page)
	if err != nil {
		return model.SearchResponse{Code: 400}, errors.New("invalid page number")
	}

	intResultsPerPage, err := strconv.Atoi(resultsPerPage)
	if err != nil {
		return model.SearchResponse{Code: 400}, errors.New("invalid number of results per page")
	}

	queryWords := WordSplit(q)
	log.Info("queryWords: ", queryWords)

	// Search the index and calculate scores
	results, err := SearchIndex(queryWords, intPage, intResultsPerPage)
	if err != nil {
		return model.SearchResponse{Code: 500}, errors.New("error fetching documents")
	}

	// Store search results in cache
	cache.Set(cacheKey, results)

	return model.SearchResponse{Code: 200, Results: results}, nil
}
