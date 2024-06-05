package controller

import (
	"KnowledgeAcquisition/logic"
	"KnowledgeAcquisition/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 分页查询
// @Description 根据关键词分页查询查询
// @Tags 查询接口
// @Produce application/json
// @Param q query string true "查询参数"
// @Param page query int true "页数"
// @Param limit query int true "每页大小"
// @Success 200 {object} model.SearchResponse
// @Router /search [get]
func Search(c *gin.Context) {
	result, err := logic.PerformSearch(c.Query("q"), c.Query("page"),
		c.DefaultQuery("limit", strconv.Itoa(model.DEAFULT_RESULT_PER_PAGE)))
	if err != nil {
		c.JSON(result.Code, err.Error())
		return
	}

	c.JSON(200, result.Results)
}

// @Summary 上传图片查询
// @Description 根据关键词分页查询查询
// @Tags 查询接口
// @Produce application/json
// @Param image formData file true "查询照片"
// @Success 200 {object} model.SearchImageResponse
// @Router /search_by_image [post]
func SearchByImage(c *gin.Context) {
	file, _ := c.FormFile("image")
	dst := "./upload_images/" + file.Filename
	c.SaveUploadedFile(file, dst)

	keywords, err := logic.GetKeywordsFromImage(dst)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	result, err := logic.PerformSearch(keywords, strconv.Itoa(1),
		strconv.Itoa(model.DEAFULT_RESULT_PER_PAGE))
	if err != nil {
		c.JSON(result.Code, err.Error())
		return
	}
	c.JSON(200, &model.SearchImageResponse{result.Results, keywords})
}
