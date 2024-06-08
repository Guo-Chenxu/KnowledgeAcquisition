package logic

import (
	"runtime/debug"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/yanyiwu/gojieba"
)

var seg *gojieba.Jieba = gojieba.NewJieba()

var stopWords = []string{" ", "\n", "\t"}

func WordSplit(query string) []string {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			log.Errorf("%v, %s", panicInfo, string(debug.Stack()))
		}
	}()

	words := seg.Cut(query, true)

	words = filter(words, stopWords)
	for i := range words {
		words[i] = strings.TrimSpace(words[i])

	}

	return words
}

func filter(slice []string, unwanted []string) []string {
	unwantedSet := make(map[string]any, len(unwanted))
	for _, s := range unwanted {
		unwantedSet[s] = struct{}{}
	}

	var result []string
	for _, s := range slice {
		if _, ok := unwantedSet[s]; !ok {
			result = append(result, s)
		}
	}

	return result
}
