package controller

import (
	"KnowledgeAcquisition/logic"

	"github.com/gin-gonic/gin"
)

// @Summary 提取关键信息
// @Description 对查询到的结果提取实体和关键词
// @Tags 提取接口
// @Produce application/json
// @Param id query string true "文档id"
// @Success 200 {object} model.DocumentAbstract
// @Router /extract_info [get]
func ExtractInfo(c *gin.Context) {
	doc_id := c.Query("id")

	result, err := logic.ExtractInfo(doc_id)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

// @Summary 提取关键信息
// @Description 对查询到的结果通过正则和词性提取指定信息点
// @Tags 提取接口
// @Produce application/json
// @Param id query string true "文档id"
// @Param pattern query string true "正则表达式"
// @Param word_class query string true "词性"
// @Success 200 {object} model.DocumentExtractRegex
// @Router /extract_info_regex [get]
func ExtractInfoRegex(c *gin.Context) {
	doc_id := c.Query("id")
	pattern := c.Query("pattern")
	word_class := c.Query("word_class")

	result, err := logic.ExtractInfoRegex(doc_id, pattern, word_class)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}
