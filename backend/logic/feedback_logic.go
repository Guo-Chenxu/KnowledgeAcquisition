package logic

import (
	"KnowledgeAcquisition/model"

	log "github.com/sirupsen/logrus"
)

func init() {
	model.ResultWeight = make(map[string]float64)
	model.ResultWordsWeight = make(map[string]map[string]float64)
	log.Info("init feedback logic")
}

func ResultFeedbackLogic(feedback model.Feedback) {
	model.ResultWeight[feedback.ResultId] = GetResultWeight(feedback.ResultId) + float64(feedback.Score)
}

func ResultWordsFeedbackLogic(feedback model.EntityFeedback) {
	model.ResultWordsWeight[feedback.ResultId][feedback.Entity] = GetResultWordsWeight(feedback.ResultId, feedback.Entity) + float64(feedback.Score)
}

func GetResultWeight(id string) float64 {
	if w, ok := model.ResultWeight[id]; ok {
		return w
	} else {
		return 1
	}
}

func GetResultWordsWeight(id string, word string) float64 {
	if m, ok := model.ResultWordsWeight[id]; ok {
		if w, ok := m[word]; ok {
			return w
		} else {
			return 1
		}
	} else {
		model.ResultWordsWeight[id] = make(map[string]float64)
		return 1
	}
}
