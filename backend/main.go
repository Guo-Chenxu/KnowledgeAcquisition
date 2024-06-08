package main

import (
	"KnowledgeAcquisition/controller"
	_ "KnowledgeAcquisition/docs"
	"KnowledgeAcquisition/logic"
	"KnowledgeAcquisition/util"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	// Load and index documents
	docs, err := util.LoadDocuments("./data")
	if err != nil {
		log.Errorf("Failed to load documents: %v", err)
		return
	}
	logic.BuildIndex(docs)
	log.Info("Indexing completed")
}

// @title 信息知识获取
// @version 1.0
// @description 信息知识获取后端接口

// @host 10.29.12.98:9011
// @BasePath /api/v1
func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		// Search with keywords
		v1.GET("/search", controller.Search)
		// Fetch SearchResult content details
		v1.GET("/document", controller.GetDocument)
		// Search by image
		v1.POST("/search_by_image", controller.SearchByImage)

		// Entities and hot words
		v1.GET("/extract_info", controller.ExtractInfo)
		// Entity and hot word feedback
		v1.POST("/extract_info_regex", controller.ExtractInfoRegex)

		// Feedback
		v1.POST("/feedback", controller.Feedback)
		// Entity Feedback
		v1.POST("/entity_feedback", controller.EntityFeedback)
		// Hotword Feedback
		v1.POST("/hotword_feedback", controller.HotwordFeedback)
		// Regex Feedback
		v1.POST("/extract_info_regex_feedback", controller.ExtractInfoRegexFeedback)
	}

	r.Run(":9011")
}
