package main

import (
	"MariaInfoRetrieval/data_process"
	"MariaInfoRetrieval/maria_types"
	"MariaInfoRetrieval/query_process"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	// log.SetLevel(log.DebugLevel)
	// log.SetReportCaller(true)
}

// Global variable to hold our "index"
// var index map[string][]Document

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	defaultResultsPerPage := 10

	// Load and index documents
	docs, err := data_process.LoadDocuments("./data")
	if err != nil {
		log.Errorf("Failed to load documents: %v", err)
		return
	}
	query_process.BuildIndex(docs)

	r.GET("/search", func(c *gin.Context) {
		result, err := query_process.PerformSearch(c.Query("q"), c.Query("page"), c.DefaultQuery("limit", strconv.Itoa(defaultResultsPerPage)))
		if err != nil {
			c.JSON(result.Code, err.Error())
			return
		}

		c.JSON(200, result.Results)
	})

	// Fetch SearchResult content details
	r.GET("/document", func(c *gin.Context) {
		id := c.Query("id")
		// Search for the document with the given id in docs
		doc, found := query_process.GetFullDoc(id)
		if !found {
			c.JSON(404, gin.H{"error": "Document" + id + " not found"})
			return
		}

		c.JSON(200, doc)
	})

	// Search by image
	r.POST("/search_by_image", func(c *gin.Context) {
		file, _ := c.FormFile("image")
		// Save file to the server
		dst := "/tmp/MariaInfoRetrieval/" + file.Filename
		c.SaveUploadedFile(file, dst)

		// Call a function to send this image to the python service, and get
		// back the keywords
		keywords, err := query_process.GetKeywordsFromImage(dst)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Then use these keywords to perform a search
		result, err := query_process.PerformSearch(keywords, strconv.Itoa(1), strconv.Itoa(defaultResultsPerPage))
		if err != nil {
			c.JSON(result.Code, err.Error())
			return
		}
		c.JSON(200, gin.H{"results": result.Results, "keywords": keywords})
	})

	// Entities and hot words
	r.GET("/extract_info", func(c *gin.Context) {
		doc_id := c.Query("id")

		// Extract entities and hot words
		result, err := query_process.ExtractInfo(doc_id)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
		c.JSON(200, gin.H{"entities": result.Entities, "hot_words": result.HotWords})
	})

	// Handle feedback
	r.POST("/feedback", func(c *gin.Context) {
		var feedback maria_types.Feedback
		if err := c.BindJSON(&feedback); err != nil {
			c.JSON(400, gin.H{"error": "Failed to parse request body: " + err.Error()})
			return
		}

		// Process the feedback here...
		log.Infof("Received feedback: %v", feedback)

		c.JSON(200, gin.H{"message": "Feedback received successfully"})
	})

	// Handle entity feedback
	r.POST("/entity_feedback", func(c *gin.Context) {
		var feedback maria_types.EntityFeedback
		if err := c.BindJSON(&feedback); err != nil {
			c.JSON(400, gin.H{"error": "Failed to parse request body: " + err.Error()})
			return
		}

		// Process the feedback here...
		log.Infof("Received feedback: %v", feedback)

		c.JSON(200, gin.H{"message": "Feedback for " + feedback.Entity + " received successfully"})
	})

	// Handle hot word feedback
	r.POST("/hotword_feedback", func(c *gin.Context) {
		var feedback maria_types.EntityFeedback
		if err := c.BindJSON(&feedback); err != nil {
			c.JSON(400, gin.H{"error": "Failed to parse request body: " + err.Error()})
			return
		}

		// Process the feedback here...
		log.Infof("Received feedback: %v", feedback)

		c.JSON(200, gin.H{"message": "Feedback for " + feedback.Entity + "received successfully"})
	})

	r.Run(":9011")
}
