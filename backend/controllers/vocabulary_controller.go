package controllers

import (
	"net/http"
	"strconv"

	"babyhabit/models"
	"babyhabit/utils"

	"github.com/gin-gonic/gin"
)

// GetVocabularyPlan 获取今日学习计划
func GetVocabularyPlan(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	userID := user.ID

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取新单词
	newWords, err := models.GetNewVocabularies(userID, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取新单词失败"})
		return
	}

	// 获取需要复习的单词
	reviewRecords, err := models.GetDueReviewVocabularies(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取复习单词失败"})
		return
	}

	// 学习统计
	stats, err := models.GetLearningStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取学习统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plan": gin.H{
			"newWords":    len(newWords),
			"reviewWords": len(reviewRecords),
		},
		"stats": stats,
	})
}

// StartVocabularyLearning 开始学习
func StartVocabularyLearning(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	userID := user.ID

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取新单词
	newWords, err := models.GetNewVocabularies(userID, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取新单词失败"})
		return
	}

	// 获取需要复习的单词
	reviewRecords, err := models.GetDueReviewVocabularies(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取复习单词失败"})
		return
	}

	// 合并单词列表
	var words []map[string]interface{}

	// 处理新单词
	for _, word := range newWords {
		// 创建学习记录
		err := models.CreateLearningRecord(userID, word.ID)
		if err != nil {
			continue
		}

		words = append(words, map[string]interface{}{
			"id":               word.ID,
			"english":          word.English,
			"chinese":          word.Chinese,
			"phonetic":         word.Phonetic,
			"audio_url":        word.AudioURL,
			"example_sentence": word.ExampleSentence,
			"type":             word.Type,
			"is_new":           true,
		})
	}

	// 处理复习单词
	for _, record := range reviewRecords {
		words = append(words, map[string]interface{}{
			"id":                 record.Vocabulary.ID,
			"english":            record.Vocabulary.English,
			"chinese":            record.Vocabulary.Chinese,
			"phonetic":           record.Vocabulary.Phonetic,
			"audio_url":          record.Vocabulary.AudioURL,
			"example_sentence":   record.Vocabulary.ExampleSentence,
			"type":               record.Vocabulary.Type,
			"is_new":             false,
			"learning_record_id": record.ID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"words": words,
	})
}

// GetVocabularyOptions 获取词汇选项
func GetVocabularyOptions(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	userID := user.ID

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	wordIDStr := c.Query("wordId")
	wordID, err := strconv.Atoi(wordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的单词ID"})
		return
	}

	optionType := c.Query("type")
	if optionType == "" {
		optionType = "chineseToEnglish"
	}

	options, err := models.GetVocabularyOptions(wordID, optionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取选项失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"options": options,
	})
}

// RecordVocabularyLearning 记录学习结果
func RecordVocabularyLearning(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	userID := user.ID

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var request struct {
		WordID    int    `json:"wordId"`
		IsCorrect bool   `json:"isCorrect"`
		CheckType string `json:"checkType"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 获取或创建学习记录
	record, err := models.GetLearningRecordByUserAndVocab(userID, request.WordID)
	if err != nil {
		// 记录不存在，创建新记录
		err = models.CreateLearningRecord(userID, request.WordID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建学习记录失败"})
			return
		}

		// 重新获取记录
		record, err = models.GetLearningRecordByUserAndVocab(userID, request.WordID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取学习记录失败"})
			return
		}
	}

	// 更新学习记录
	err = models.UpdateLearningRecord(record.ID, request.IsCorrect)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新学习记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "学习记录已更新",
	})
}

// GetVocabularyStats 获取词汇学习统计
func GetVocabularyStats(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	userID := user.ID

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	stats, err := models.GetLearningStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取学习统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// CreateVocabulary 创建词汇（管理员）
func CreateVocabulary(c *gin.Context) {
	var vocab models.Vocabulary
	if err := c.ShouldBindJSON(&vocab); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 生成语音文件
	audioURL, err := utils.GenerateSpeech(vocab.English)
	if err == nil {
		vocab.AudioURL = &audioURL
	}

	// 生成例句
	exampleSentence, err := utils.GenerateExampleSentence(vocab.English)
	if err == nil {
		vocab.ExampleSentence = &exampleSentence
	}

	err = models.CreateVocabulary(&vocab)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建词汇失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "词汇创建成功",
		"vocabulary": vocab,
	})
}

// GetVocabularies 获取词汇列表（管理员）
func GetVocabularies(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	vocabularies, total, err := models.GetVocabularies(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取词汇列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vocabularies": vocabularies,
		"total":        total,
	})
}

// GetVocabulary 获取单个词汇（管理员）
func GetVocabulary(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的词汇ID"})
		return
	}

	vocab, err := models.GetVocabularyByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取词汇失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vocabulary": vocab,
	})
}

// UpdateVocabulary 更新词汇（管理员）
func UpdateVocabulary(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的词汇ID"})
		return
	}

	var vocab models.Vocabulary
	if err := c.ShouldBindJSON(&vocab); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	vocab.ID = id
	err = models.UpdateVocabulary(&vocab)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新词汇失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "词汇更新成功",
		"vocabulary": vocab,
	})
}

// DeleteVocabulary 删除词汇（管理员）
func DeleteVocabulary(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的词汇ID"})
		return
	}

	err = models.DeleteVocabulary(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除词汇失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "词汇删除成功",
	})
}

// BatchCreateVocabulary 批量创建词汇（管理员）
func BatchCreateVocabulary(c *gin.Context) {
	var request struct {
		Vocabularies []models.Vocabulary `json:"vocabularies"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 转换为指针切片并生成音频和例句
	var vocabPtrs []*models.Vocabulary
	for i := range request.Vocabularies {
		vocab := &request.Vocabularies[i]

		// 生成语音文件
		audioURL, err := utils.GenerateSpeech(vocab.English)
		if err == nil {
			vocab.AudioURL = &audioURL
		}

		// 生成例句
		exampleSentence, err := utils.GenerateExampleSentence(vocab.English)
		if err == nil {
			vocab.ExampleSentence = &exampleSentence
		}

		vocabPtrs = append(vocabPtrs, vocab)
	}

	err := models.BatchCreateVocabulary(vocabPtrs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "批量创建词汇失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "词汇批量创建成功",
		"count":   len(request.Vocabularies),
	})
}

// BatchDeleteVocabulary 批量删除词汇（管理员）
func BatchDeleteVocabulary(c *gin.Context) {
	var request struct {
		IDs []int `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if len(request.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要删除的词汇"})
		return
	}

	err := models.BatchDeleteVocabulary(request.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "批量删除词汇失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "词汇批量删除成功",
		"count":   len(request.IDs),
	})
}
