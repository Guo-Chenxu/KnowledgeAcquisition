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
	w, ok := model.ResultWeight[feedback.ResultId]
	if ok {
		model.ResultWeight[feedback.ResultId] = w + float64(feedback.Score)
	} else {
		model.ResultWeight[feedback.ResultId] = 1 + float64(feedback.Score)
	}
}

func ResultWordsFeedbackLogic(feedback model.EntityFeedback) {
	m, ok := model.ResultWordsWeight[feedback.ResultId]
	if ok {
		if w, ok := m[feedback.Entity]; ok {
			model.ResultWordsWeight[feedback.ResultId][feedback.Entity] = w + float64(feedback.Score)
		}
	} else {
		model.ResultWordsWeight[feedback.ResultId] = make(map[string]float64)
		model.ResultWordsWeight[feedback.ResultId][feedback.Entity] = 1 + float64(feedback.Score)
	}

}
