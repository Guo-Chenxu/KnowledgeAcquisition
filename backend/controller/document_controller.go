package controller

import (
	"KnowledgeAcquisition/logic"

	"github.com/gin-gonic/gin"
)

// @Summary 获取文档详细信息
// @Description 根据关键词分页查询查询
// @Tags 查询接口
// @Produce application/json
// @Param id query string true "文档id"
// @Success 200 {object} model.Document
// @Router /document [get]
func GetDocument(c *gin.Context) {
	id := c.Query("id")

	doc, found := logic.GetFullDoc(id)
	if !found {
		c.JSON(404, gin.H{"error": "Document" + id + " not found"})
		return
	}

	c.JSON(200, doc)
}
