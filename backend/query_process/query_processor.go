package query_process

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/yanyiwu/gojieba"
)

var seg *gojieba.Jieba = gojieba.NewJieba()

func WordSplit(query string) []string {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			fmt.Printf("%v, %s", panicInfo, string(debug.Stack()))
		}
	}()

	words := seg.Cut(query, true) // Use Jeba for Chinese text segment

	words = whitespaceFilter(words)
	for i := range words {
		words[i] = strings.TrimSpace(words[i])

	}

	return words
}

func filter(slice []string, unwanted ...string) []string {
	unwantedSet := make(map[string]struct{}, len(unwanted))
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

func whitespaceFilter(slice []string) []string {
	unwanted := []string{" ", "\t"}
	return filter(slice, unwanted...)
}
