package controller

import (
	"KnowledgeAcquisition/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// @Summary 结果反馈
// @Description 对结果进行反馈
// @Tags 反馈接口
// @Accept application/json
// @Produce application/json
// @Param feedback body model.Feedback true "反馈参数"
// @Success 200 {object} string
// @Router /feedback [post]
func Feedback(c *gin.Context) {
	var feedback model.Feedback
	if err := c.BindJSON(&feedback); err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse request body: " + err.Error()})
		return
	}
	
	log.Infof("Received feedback: %v", feedback)
	
	c.JSON(200, gin.H{"message": "Feedback received successfully"})
}

// @Summary 实体反馈
// @Description 对实体类进行反馈
// @Tags 反馈接口
// @Accept application/json
// @Produce application/json
// @Param feedback body model.EntityFeedback true "反馈参数"
// @Success 200 {object} string
// @Router /entity_feedback [post]
func EntityFeedback(c *gin.Context) {
	var feedback model.EntityFeedback
	if err := c.BindJSON(&feedback); err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse request body: " + err.Error()})
		return
	}
	
	log.Infof("Received feedback: %v", feedback)
	
	c.JSON(200, gin.H{"message": "Feedback received successfully"})
}

// @Summary 热词反馈
// @Description 对热词进行反馈
// @Tags 反馈接口
// @Accept application/json
// @Produce application/json
// @Param feedback body model.EntityFeedback true "反馈参数"
// @Success 200 {object} string
// @Router /hotword_feedback [post]
func HotwordFeedback(c *gin.Context) {
	var feedback model.EntityFeedback
	if err := c.BindJSON(&feedback); err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse request body: " + err.Error()})
		return
	}
	
	log.Infof("Received feedback: %v", feedback)
	
	c.JSON(200, gin.H{"message": "Feedback received successfully"})
}

// @Summary 正则提取反馈
// @Description 对正则提取结果进行反馈
// @Tags 反馈接口
// @Accept application/json
// @Produce application/json
// @Param feedback body model.EntityFeedback true "反馈参数"
// @Success 200 {object} string
// @Router /extract_info_regex_feedback [post]
func ExtractInfoRegexFeedback(c *gin.Context) {
	var feedback model.EntityFeedback
	if err := c.BindJSON(&feedback); err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse request body: " + err.Error()})
		return
	}
	
	log.Infof("Received feedback: %v", feedback)
	
	c.JSON(200, gin.H{"message": "Feedback for " + feedback.Entity + "received successfully"})
}
