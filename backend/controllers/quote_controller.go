package controllers

import (
	"babyhabit/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateQuote 创建名言警句
func CreateQuote(c *gin.Context) {
	var quote models.Quote
	if err := c.ShouldBindJSON(&quote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.CreateQuote(&quote); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quote": quote})
}

// GetQuotes 获取名言警句列表（支持分页）
func GetQuotes(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	quotes, total, err := models.GetQuotesWithPagination(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"quotes": quotes,
		"total":  total,
		"page":   page,
		"size":   pageSize,
	})
}

// GetQuote 获取单个名言警句
func GetQuote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
		return
	}

	quote, err := models.GetQuoteByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quote": quote})
}

// GetRandomQuote 获取随机名言警句
func GetRandomQuote(c *gin.Context) {
	quote, err := models.GetRandomQuote()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if quote == nil {
		c.JSON(http.StatusOK, gin.H{"quote": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quote": quote})
}

// UpdateQuote 更新名言警句
func UpdateQuote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
		return
	}

	var quote models.Quote
	if err := c.ShouldBindJSON(&quote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quote.ID = id
	if err := models.UpdateQuote(&quote); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quote": quote})
}

// DeleteQuote 删除名言警句
func DeleteQuote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
		return
	}

	if err := models.DeleteQuote(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quote deleted successfully"})
}

// BatchDeleteQuotes 批量删除名言警句
func BatchDeleteQuotes(c *gin.Context) {
	var req struct {
		IDs []int64 `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.BatchDeleteQuotes(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quotes deleted successfully"})
}

// BatchCreateQuotes 批量创建名言警句
func BatchCreateQuotes(c *gin.Context) {
	var req struct {
		Quotes []struct {
			Content string `json:"content" binding:"required"`
			Meaning string `json:"meaning"`
			Author  string `json:"author"`
		} `json:"quotes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	totalCount := len(req.Quotes)
	quotes := make([]*models.Quote, totalCount)
	for i, q := range req.Quotes {
		quotes[i] = &models.Quote{
			Content: q.Content,
			Meaning: q.Meaning,
			Author:  q.Author,
		}
	}

	// 获取导入前的数据量
	_, beforeCount, err := models.GetQuotesWithPagination(1, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := models.BatchCreateQuotes(quotes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取导入后的数据量
	_, afterCount, err := models.GetQuotesWithPagination(1, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 计算成功导入的数量和重复数量
	successCount := afterCount - beforeCount
	duplicateCount := totalCount - successCount

	c.JSON(http.StatusOK, gin.H{
		"message":         "Quotes created successfully",
		"total_count":     totalCount,
		"success_count":   successCount,
		"duplicate_count": duplicateCount,
	})
}
