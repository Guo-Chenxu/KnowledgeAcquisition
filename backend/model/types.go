package model

type Document struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Keywords string `json:"keywords"`
	URL      string `json:"url"`
	Date     string `json:"date"`
	Lang     Language
}

type DocumentAbstract struct {
	Entities map[string]int `json:"entities"`
	HotWords map[string]int `json:"hot_words"`
}

type DocumentExtractRegex struct {
	Words []string `json:"words"`
}

type SummaryDocument struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	URL     string `json:"url"`
	Date    string `json:"date"`
}

type SearchResult struct {
	Score float64
	Doc   SummaryDocument
}

type Feedback struct {
	ResultId string `json:"resultId"`
	Score    int    `json:"Score"`
}

type EntityFeedback struct {
	ResultId string `json:"resultId"`
	Entity   string `json:"item"`
	Score    int    `json:"score"`
}

type SearchResponse struct {
	Code    int            `json:"code"`
	Results []SearchResult `json:"results"`
}

type Language int8

const (
	English Language = iota
	Chinese
)

func (l Language) String() string {
	if l == English {
		return "en"
	} else if l == Chinese {
		return "cn"
	}
	return "undefined"
}
