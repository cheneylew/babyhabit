package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheneylew/babyhabit/backend/models"
	"github.com/cheneylew/babyhabit/backend/utils"

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

	// 获取教材ID参数
	bookIDsStr := c.Query("book_ids")
	var bookIDs []int
	if bookIDsStr != "" {
		for _, idStr := range strings.Split(bookIDsStr, ",") {
			// 去除空格
			idStr = strings.TrimSpace(idStr)
			// 跳过空字符串
			if idStr == "" {
				continue
			}
			id, err := strconv.Atoi(idStr)
			if err == nil && id > 0 {
				bookIDs = append(bookIDs, id)
			}
		}
	}

	// 获取用户每日单词数量偏好设置，默认为5个
	dailyWordLimit := 5
	preference, err := models.GetUserPreference(userID, "daily_word_limit")
	if err == nil && preference != nil {
		if limit, err := strconv.Atoi(preference.PreferenceValue); err == nil && limit > 0 {
			dailyWordLimit = limit
		}
	}

	// 获取用户每日单词复习量上限偏好设置，默认为18个
	dailyReviewLimit := 18
	reviewPreference, err := models.GetUserPreference(userID, "daily_review_limit")
	if err == nil && reviewPreference != nil {
		if limit, err := strconv.Atoi(reviewPreference.PreferenceValue); err == nil && limit > 0 {
			dailyReviewLimit = limit
		}
	}

	// 获取今日已经学习的新单词数量
	learnedToday, err := models.GetTodayLearnedNewWordsCount(userID)
	if err != nil {
		learnedToday = 0
	}

	// 获取需要复习的单词
	reviewRecords, err := models.GetDueReviewVocabularies(userID, bookIDs)
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

	// 如果复习单词数量超过上限，则今日不学习新单词
	if len(reviewRecords)+stats["todayReviewedWords"].(int) >= dailyReviewLimit {
		dailyWordLimit = 0
	}

	// 计算剩余需要学习的新单词数量
	remainingNewWords := dailyWordLimit - learnedToday
	if remainingNewWords < 0 {
		remainingNewWords = 0
	}

	// 检查是否还有新单词需要学习
	if remainingNewWords > 0 {
		// 获取新单词
		_, err := models.GetNewVocabularies(userID, remainingNewWords, bookIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取新单词失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"plan": gin.H{
			"newWords":    dailyWordLimit,
			"reviewWords": len(reviewRecords),
		},
		"stats": stats,
	})
}

// StartVocabularyLearning 开始学习
func StartVocabularyLearning(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	userID := user.ID

	// 从环境变量获取前端静态访问地址
	frontendStaticURL := os.Getenv("FRONTEND_STATIC_URL")
	if frontendStaticURL == "" {
		frontendStaticURL = "http://localhost:8000"
	}

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取教材ID参数
	bookIDsStr := c.Query("book_ids")
	var bookIDs []int
	if bookIDsStr != "" {
		for _, idStr := range strings.Split(bookIDsStr, ",") {
			// 去除空格
			idStr = strings.TrimSpace(idStr)
			// 跳过空字符串
			if idStr == "" {
				continue
			}
			id, err := strconv.Atoi(idStr)
			if err == nil && id > 0 {
				bookIDs = append(bookIDs, id)
			}
		}
	}

	// 获取今天已经学习的新单词数量
	learnedToday, err := models.GetTodayLearnedNewWordsCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取今日学习记录失败"})
		return
	}

	// 获取用户每日单词数量偏好设置，默认为5个
	dailyWordLimit := 5
	preference, err := models.GetUserPreference(userID, "daily_word_limit")
	if err == nil && preference != nil {
		if limit, err := strconv.Atoi(preference.PreferenceValue); err == nil && limit > 0 {
			dailyWordLimit = limit
		}
	}

	// 获取用户每日单词复习量上限偏好设置，默认为18个
	dailyReviewLimit := 18
	reviewPreference, err := models.GetUserPreference(userID, "daily_review_limit")
	if err == nil && reviewPreference != nil {
		if limit, err := strconv.Atoi(reviewPreference.PreferenceValue); err == nil && limit > 0 {
			dailyReviewLimit = limit
		}
	}

	// 获取需要复习的单词
	reviewRecords, err := models.GetDueReviewVocabularies(userID, bookIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取复习单词失败"})
		return
	}

	// 学习统计
	todayReviewedWords, err := models.StatGetTodayReviewedWords(userID)
	if err != nil {
		todayReviewedWords = 0
	}

	// 如果复习单词数量超过上限，则今日不学习新单词
	if len(reviewRecords)+todayReviewedWords >= dailyReviewLimit {
		dailyWordLimit = 0
	}

	// 计算今天还可以学习的新单词数量
	remainingNewWords := dailyWordLimit - learnedToday
	if remainingNewWords < 0 {
		remainingNewWords = 0
	}

	// 获取新单词
	newWords, err := models.GetNewVocabularies(userID, remainingNewWords, bookIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取新单词失败"})
		return
	}

	// 合并单词列表
	var words []map[string]interface{}

	// 处理新单词
	for _, word := range newWords {
		audioURL := ""
		if word.AudioURL != nil && *word.AudioURL != "" {
			// 检查是否已经是完整的URL
			if strings.HasPrefix(*word.AudioURL, "http://") || strings.HasPrefix(*word.AudioURL, "https://") {
				audioURL = *word.AudioURL
			} else {
				audioURL = frontendStaticURL + *word.AudioURL
			}
		}

		// 处理 example_sentence
		exampleSentence := word.ExampleSentence
		if exampleSentence == nil || *exampleSentence == "null" {
			null := "[]"
			exampleSentence = &null
		}

		words = append(words, map[string]interface{}{
			"id":               word.ID,
			"english":          word.English,
			"chinese":          word.Chinese,
			"phonetic":         word.Phonetic,
			"audio_url":        audioURL,
			"example_sentence": exampleSentence,
			"type":             word.Type,
			"is_new":           true,
		})
	}

	// 处理复习单词
	for _, record := range reviewRecords {
		audioURL := ""
		if record.Vocabulary.AudioURL != nil && *record.Vocabulary.AudioURL != "" {
			// 检查是否已经是完整的URL
			if strings.HasPrefix(*record.Vocabulary.AudioURL, "http://") || strings.HasPrefix(*record.Vocabulary.AudioURL, "https://") {
				audioURL = *record.Vocabulary.AudioURL
			} else {
				audioURL = frontendStaticURL + *record.Vocabulary.AudioURL
			}
		}

		// 处理 example_sentence
		exampleSentence := record.Vocabulary.ExampleSentence
		if exampleSentence == nil || *exampleSentence == "null" {
			null := "[]"
			exampleSentence = &null
		}

		words = append(words, map[string]interface{}{
			"id":                 record.Vocabulary.ID,
			"english":            record.Vocabulary.English,
			"chinese":            record.Vocabulary.Chinese,
			"phonetic":           record.Vocabulary.Phonetic,
			"audio_url":          audioURL,
			"example_sentence":   exampleSentence,
			"type":               record.Vocabulary.Type,
			"is_new":             false,
			"learning_record_id": record.ID,
			"remark":             record.Vocabulary.Remark,
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
		Mastered  bool   `json:"mastered"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 获取或创建学习记录
	record, err := models.GetLearningRecordByUserAndVocab(userID, request.WordID)
	isNewRecord := false
	if err != nil {
		// 记录不存在，创建新记录
		err = models.CreateLearningRecord(userID, request.WordID, request.CheckType)
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
		isNewRecord = true
	}

	// 更新学习记录
	if request.Mastered {
		// 直接标记为掌握
		err = models.MarkVocabularyAsMastered(record.ID, request.CheckType)
	} else {
		err = models.UpdateLearningRecord(record.ID, request.IsCorrect)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新学习记录失败"})
		return
	}

	// 创建或更新学习打卡记录
	newWordsCount := 0
	if isNewRecord && request.CheckType != "skip" {
		newWordsCount = 1
	}
	reviewWordsCount := 0
	if !isNewRecord {
		reviewWordsCount = 1
	}
	correctCount := 0
	if request.IsCorrect {
		correctCount = 1
	}
	totalCount := 1
	durationMinutes := 1 // 假设每次学习单词花费1分钟

	err = models.CreateStudyCheckin(userID, newWordsCount, reviewWordsCount, correctCount, totalCount, durationMinutes)
	if err != nil {
		// 记录打卡失败不影响学习记录的更新
		fmt.Println("创建学习打卡记录失败:", err)
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

// GetVocabularyHistory 获取学习历史
func GetVocabularyHistory(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	userID := user.ID

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	history, err := models.GetLearningHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取学习历史失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"words": history,
	})
}

// CreateVocabulary 创建词汇（管理员）
func CreateVocabulary(c *gin.Context) {
	var vocab models.Vocabulary
	if err := c.ShouldBindJSON(&vocab); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 验证教材必填
	if vocab.BookID == nil || *vocab.BookID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "教材不能为空"})
		return
	}

	// 生成完整的单词信息
	wordInfo, err := utils.GenerateWordInfo(vocab.English)
	if err == nil {
		// 设置中文翻译
		vocab.Chinese = wordInfo.Chinese

		// 设置音标（将PhoneticInfo转换为JSON字符串）
		phoneticJSON, err := json.Marshal(wordInfo.Phonetic)
		if err == nil {
			phoneticStr := string(phoneticJSON)
			vocab.Phonetic = &phoneticStr
		}

		// 设置例句（将[]Example转换为JSON字符串）
		examplesJSON, err := json.Marshal(wordInfo.Examples)
		if err == nil {
			examplesStr := string(examplesJSON)
			vocab.ExampleSentence = &examplesStr
		}

		// 设置分类
		category := wordInfo.Category
		vocab.Category = &category

		// 设置音频URL
		if wordInfo.AudioURL != "" {
			vocab.AudioURL = &wordInfo.AudioURL
		}
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
	search := c.Query("search")
	bookIDStr := c.Query("book_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var bookID int
	if bookIDStr != "" {
		bookID, err = strconv.Atoi(bookIDStr)
		if err != nil {
			bookID = 0
		}
	}

	vocabularies, total, err := models.GetVocabularies(page, pageSize, search, bookID)
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

	// 检查是否需要生成完整的单词信息
	needGenerateInfo := false
	fmt.Println("Checking needGenerateInfo...")
	fmt.Println("vocab.Chinese == \"\":", vocab.Chinese == "")
	fmt.Println("vocab.ExampleSentence == nil:", vocab.ExampleSentence == nil)
	if vocab.ExampleSentence != nil {
		fmt.Println("*vocab.ExampleSentence == \"\":", *vocab.ExampleSentence == "")
	}
	fmt.Println("vocab.Category == nil:", vocab.Category == nil)
	if vocab.Category != nil {
		fmt.Println("*vocab.Category == \"\":", *vocab.Category == "")
	}

	if vocab.Chinese == "" || (vocab.ExampleSentence == nil || *vocab.ExampleSentence == "") || (vocab.Category == nil || *vocab.Category == "") {
		needGenerateInfo = true
		fmt.Println("needGenerateInfo set to true")
	} else {
		fmt.Println("needGenerateInfo set to false")
	}

	// 如果需要生成信息，调用大模型
	if needGenerateInfo {
		wordInfo, err := utils.GenerateWordInfo(vocab.English)
		if err == nil {
			// 如果中文为空，设置中文翻译
			if vocab.Chinese == "" {
				vocab.Chinese = wordInfo.Chinese
			}

			// 如果音标为空，设置音标
			if vocab.Phonetic == nil {
				phoneticJSON, err := json.Marshal(wordInfo.Phonetic)
				if err == nil {
					phoneticStr := string(phoneticJSON)
					vocab.Phonetic = &phoneticStr
				}
			}

			// 如果例句为空，设置例句
			if vocab.ExampleSentence == nil || *vocab.ExampleSentence == "" {
				examplesJSON, err := json.Marshal(wordInfo.Examples)
				if err == nil {
					examplesStr := string(examplesJSON)
					vocab.ExampleSentence = &examplesStr
				}
			}

			// 如果分类为空，设置分类
			if vocab.Category == nil || *vocab.Category == "" {
				category := wordInfo.Category
				vocab.Category = &category
			}

			// 如果音频URL为空，设置音频URL
			if vocab.AudioURL == nil && wordInfo.AudioURL != "" {
				vocab.AudioURL = &wordInfo.AudioURL
			}
		}
	}
	// 单独处理音频URL缺失的情况
	if vocab.AudioURL == nil {
		audioURL, err := utils.GenerateSpeech(vocab.English)
		if err == nil && audioURL != "" {
			vocab.AudioURL = &audioURL
		}
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
		BookID       int                 `json:"book_id"`
		Remark       string              `json:"remark"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 验证教材必填
	if request.BookID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "教材不能为空"})
		return
	}

	// 检查是否有重复词汇
	existingWords, err := models.GetVocabulariesByBookID(request.BookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查词汇重复失败"})
		return
	}

	// 构建已存在词汇的映射
	existingWordMap := make(map[string]bool)
	for _, word := range existingWords {
		existingWordMap[word.English] = true
	}

	// 过滤掉重复的词汇
	var filteredVocabularies []models.Vocabulary
	for _, vocab := range request.Vocabularies {
		if !existingWordMap[vocab.English] {
			filteredVocabularies = append(filteredVocabularies, vocab)
		}
	}

	// 如果没有需要导入的词汇
	if len(filteredVocabularies) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":    "所有词汇已存在，无需导入",
			"count":      0,
			"total":      len(request.Vocabularies),
			"duplicates": len(request.Vocabularies),
		})
		return
	}

	// 设置响应头为JSON流
	c.Header("Content-Type", "application/json")
	c.Header("Transfer-Encoding", "chunked")

	// 开始逐个处理词汇
	successCount := 0
	totalCount := len(filteredVocabularies)
	duplicateCount := len(request.Vocabularies) - totalCount

	for i, vocab := range filteredVocabularies {
		// 为每个词汇创建新的实例
		currentVocab := vocab
		currentVocab.BookID = &request.BookID
		// 设置备注
		if request.Remark != "" {
			currentVocab.Remark = &request.Remark
		}

		// 检查其他教材是否已经有同样的英文单词
		existingVocab, err := models.GetVocabularyByEnglish(currentVocab.English)
		var wordInfo *utils.WordInfo
		if err == nil && existingVocab != nil {
			// 其他教材已经有同样的英文单词，直接复制整行过来
			currentVocab.Chinese = existingVocab.Chinese
			currentVocab.Phonetic = existingVocab.Phonetic
			currentVocab.AudioURL = existingVocab.AudioURL
			currentVocab.ExampleSentence = existingVocab.ExampleSentence
			currentVocab.Type = existingVocab.Type
			currentVocab.Category = existingVocab.Category
		} else {
			// 生成完整的单词信息（重试机制）
			var infoErr error
			maxRetries := 3
			for retry := 0; retry < maxRetries; retry++ {
				wordInfo, infoErr = utils.GenerateWordInfo(currentVocab.English)
				if infoErr == nil && wordInfo != nil {
					break
				}
				// 重试间隔
				time.Sleep(time.Second * time.Duration(retry+1))
			}

			// 处理单词信息
			if infoErr == nil && wordInfo != nil {
				// 设置中文翻译
				currentVocab.Chinese = wordInfo.Chinese

				// 设置音标（将PhoneticInfo转换为JSON字符串）
				phoneticJSON, err := json.Marshal(wordInfo.Phonetic)
				if err == nil {
					phoneticStr := string(phoneticJSON)
					currentVocab.Phonetic = &phoneticStr
				}

				// 设置例句（将[]Example转换为JSON字符串）
				examplesJSON, err := json.Marshal(wordInfo.Examples)
				if err == nil {
					examplesStr := string(examplesJSON)
					currentVocab.ExampleSentence = &examplesStr
				}

				// 设置分类
				category := wordInfo.Category
				currentVocab.Category = &category

				// 设置音频URL
				if wordInfo.AudioURL != "" {
					currentVocab.AudioURL = &wordInfo.AudioURL
				}
			}
		}

		// 保存词汇
		err = models.CreateVocabulary(&currentVocab)
		if err == nil {
			successCount++
		}

		// 发送进度更新
		progress := map[string]interface{}{
			"status":     "processing",
			"current":    i + 1,
			"total":      totalCount,
			"success":    successCount,
			"duplicates": duplicateCount,
			"word":       currentVocab.English,
			"has_info":   wordInfo != nil,
		}

		// 序列化并发送
		if data, err := json.Marshal(progress); err == nil {
			c.Writer.Write(data)
			c.Writer.Write([]byte("\n"))
			c.Writer.Flush()
		}
	}

	// 发送完成状态
	finalStatus := map[string]interface{}{
		"status":     "completed",
		"total":      totalCount,
		"success":    successCount,
		"duplicates": duplicateCount,
		"message":    "词汇批量创建完成",
	}

	if data, err := json.Marshal(finalStatus); err == nil {
		c.Writer.Write(data)
		c.Writer.Write([]byte("\n"))
		c.Writer.Flush()
	}
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

// GenerateSentenceAudio 生成例句音频
func GenerateSentenceAudio(c *gin.Context) {
	var request struct {
		Sentence string `json:"sentence"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if request.Sentence == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "例句不能为空"})
		return
	}

	// 生成例句音频
	fmt.Printf("Generating audio for sentence: %s\n", request.Sentence)
	audioURL, err := utils.GenerateSentenceSpeech(request.Sentence)
	if err != nil {
		fmt.Printf("Error generating sentence audio: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成例句音频失败"})
		return
	}

	// 从环境变量获取前端静态访问地址
	frontendStaticURL := os.Getenv("FRONTEND_STATIC_URL")
	if frontendStaticURL == "" {
		frontendStaticURL = "http://localhost:8000"
	}

	// 拼接完整的音频URL
	fullAudioURL := frontendStaticURL + audioURL
	fmt.Printf("Generated audio URL: %s\n", fullAudioURL)

	c.JSON(http.StatusOK, gin.H{
		"audio_url": fullAudioURL,
	})
}

// GetWordMeaning 获取单词的中文意思
func GetWordMeaning(c *gin.Context) {
	var request struct {
		Word string `json:"word"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if request.Word == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "单词不能为空"})
		return
	}

	// 使用大模型API查询单词意思
	meaning, err := utils.GetWordMeaning(request.Word)
	if err != nil {
		fmt.Printf("Error getting word meaning: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询单词意思失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meaning": meaning})
}

// GetBooks 获取所有教材列表（管理员）
func GetBooks(c *gin.Context) {
	books, err := models.GetBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取教材列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books})
}

// GetBook 获取单个教材（管理员）
func GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的教材ID"})
		return
	}

	book, err := models.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取教材失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book})
}

// GetIncompleteVocabularies 获取缺少AI生成信息的词汇（管理员）
func GetIncompleteVocabularies(c *gin.Context) {
	vocabularies, err := models.GetIncompleteVocabularies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取缺少AI生成信息的词汇失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vocabularies": vocabularies,
	})
}

// RegenerateVocabularies 重新生成词汇的AI信息（管理员）
func RegenerateVocabularies(c *gin.Context) {
	var request struct {
		IDs []int `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 开始逐个处理词汇
	successCount := 0
	totalCount := len(request.IDs)

	// 使用c.Stream发送流式响应
	c.Stream(func(w io.Writer) bool {
		for i, id := range request.IDs {
			// 获取词汇信息
			vocab, err := models.GetVocabularyByID(id)
			if err != nil {
				// 发送错误信息
				progressData := map[string]interface{}{
					"status":  "progress",
					"current": i + 1,
					"total":   totalCount,
					"word":    "",
					"error":   err.Error(),
				}
				if data, err := json.Marshal(progressData); err == nil {
					fmt.Fprintf(w, "%s\n", data)
					w.(http.Flusher).Flush()
				}
				continue
			}

			// 生成完整的单词信息（重试机制）
			var wordInfo *utils.WordInfo
			var infoErr error
			maxRetries := 3
			for retry := 0; retry < maxRetries; retry++ {
				wordInfo, infoErr = utils.GenerateWordInfo(vocab.English)
				if infoErr == nil && wordInfo != nil {
					break
				}
				// 重试间隔
				time.Sleep(time.Second * time.Duration(retry+1))
			}

			// 处理单词信息
			if infoErr == nil && wordInfo != nil {
				// 设置中文翻译
				vocab.Chinese = wordInfo.Chinese

				// 设置音标（将PhoneticInfo转换为JSON字符串）
				phoneticJSON, err := json.Marshal(wordInfo.Phonetic)
				if err == nil {
					phoneticStr := string(phoneticJSON)
					vocab.Phonetic = &phoneticStr
				}

				// 设置例句（将[]Example转换为JSON字符串）
				examplesJSON, err := json.Marshal(wordInfo.Examples)
				if err == nil {
					examplesStr := string(examplesJSON)
					vocab.ExampleSentence = &examplesStr
				}

				// 设置分类
				category := wordInfo.Category
				vocab.Category = &category

				// 设置音频URL
				if wordInfo.AudioURL != "" {
					vocab.AudioURL = &wordInfo.AudioURL
				}

				// 更新词汇
				err = models.UpdateVocabulary(vocab)
				if err == nil {
					successCount++
				}
			}

			// 发送进度信息
			progressData := map[string]interface{}{
				"status":  "progress",
				"current": i + 1,
				"total":   totalCount,
				"word":    vocab.English,
			}
			if data, err := json.Marshal(progressData); err == nil {
				fmt.Fprintf(w, "%s\n", data)
				w.(http.Flusher).Flush()
			}
		}

		// 发送完成信息
		completedData := map[string]interface{}{
			"status":  "completed",
			"success": successCount,
			"total":   totalCount,
		}
		if data, err := json.Marshal(completedData); err == nil {
			fmt.Fprintf(w, "%s\n", data)
			w.(http.Flusher).Flush()
		}

		return false
	})
}

// Chat 聊天接口
func Chat(c *gin.Context) {
	var request struct {
		Prompt string `json:"prompt"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if request.Prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "提示词不能为空"})
		return
	}

	// 调用聊天函数
	response, err := utils.Chat(request.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "聊天失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}

// ChatStream 流式聊天接口
func ChatStream(c *gin.Context) {
	fmt.Println("[ChatStream] 开始处理流式聊天请求")

	var request struct {
		Prompt string `json:"prompt"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Printf("[ChatStream] 请求数据解析失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if request.Prompt == "" {
		fmt.Println("[ChatStream] 提示词为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "提示词不能为空"})
		return
	}

	fmt.Printf("[ChatStream] 收到提示词: %s\n", request.Prompt)

	// 设置响应头为流式响应
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	fmt.Println("[ChatStream] 设置响应头完成")

	// 创建一个通道，用于接收前台的结束信号
	done := make(chan bool)

	// 启动一个 goroutine，监听前台的结束信号
	go func() {
		fmt.Println("[ChatStream] 启动监听前台结束信号的 goroutine")
		// 读取前台发送的结束信号
		buffer := make([]byte, 1024)
		for {
			n, err := c.Request.Body.Read(buffer)
			if err != nil {
				fmt.Printf("[ChatStream] 读取请求体失败: %v，关闭 done 通道\n", err)
				// 当连接关闭时，发送结束信号
				close(done)
				break
			}
			if n > 0 {
				// 检查是否收到终止信号
				message := string(buffer[:n])
				fmt.Printf("[ChatStream] 收到前台消息: %s\n", message)
				if message == "stop" {
					// 收到终止信号，关闭通道
					fmt.Println("[ChatStream] 收到终止信号，关闭 done 通道")
					close(done)
					break
				}
			}
		}
	}()

	// 调用流式聊天函数
	fmt.Println("[ChatStream] 调用 utils.ChatStream 函数")
	err := utils.ChatStream(request.Prompt, func(chunk string) bool {
		fmt.Printf("[ChatStream] 收到 AI 响应 chunk: %s\n", chunk)
		// 发送数据到前台
		_, writeErr := fmt.Fprintf(c.Writer, "data: %s\n\n", chunk)
		if writeErr != nil {
			fmt.Printf("[ChatStream] 写入响应失败: %v\n", writeErr)
			return false
		}
		c.Writer.Flush()
		fmt.Println("[ChatStream] 响应已刷新")
		return true
	}, done)

	if err != nil {
		// 发送错误信息
		fmt.Printf("[ChatStream] 流式聊天失败: %v\n", err)
		_, writeErr := fmt.Fprintf(c.Writer, "data: {\"error\":\"%s\"}\n\n", err.Error())
		if writeErr != nil {
			fmt.Printf("[ChatStream] 写入错误响应失败: %v\n", writeErr)
		}
		c.Writer.Flush()
	}

	// 发送结束信号
	fmt.Println("[ChatStream] 发送结束信号")
	_, writeErr := fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
	if writeErr != nil {
		fmt.Printf("[ChatStream] 写入结束信号失败: %v\n", writeErr)
	}
	c.Writer.Flush()
	fmt.Println("[ChatStream] 请求处理完成")
}

// CreateBook 创建教材（管理员）
func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if book.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "教材名称不能为空"})
		return
	}

	err := models.CreateBook(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建教材失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "教材创建成功",
		"book":    book,
	})
}

// UpdateBook 更新教材（管理员）
func UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的教材ID"})
		return
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if book.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "教材名称不能为空"})
		return
	}

	book.ID = id
	err = models.UpdateBook(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新教材失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "教材更新成功",
		"book":    book,
	})
}

// DeleteBook 删除教材（管理员）
func DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的教材ID"})
		return
	}

	err = models.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除教材失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "教材删除成功",
	})
}

// GetBookOptions 获取教材选项列表
func GetBookOptions(c *gin.Context) {
	books, err := models.GetBookOptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取教材列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"books": books,
	})
}

// GetVocabularyDictation 获取今日需要默写的单词
func GetVocabularyDictation(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	userID := user.ID

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取教材ID参数
	bookIDsStr := c.Query("book_ids")
	var bookIDs []int
	if bookIDsStr != "" {
		for _, idStr := range strings.Split(bookIDsStr, ",") {
			// 去除空格
			idStr = strings.TrimSpace(idStr)
			// 跳过空字符串
			if idStr == "" {
				continue
			}
			id, err := strconv.Atoi(idStr)
			if err == nil && id > 0 {
				bookIDs = append(bookIDs, id)
			}
		}
	}

	// 从环境变量获取前端静态访问地址
	frontendStaticURL := os.Getenv("FRONTEND_STATIC_URL")
	if frontendStaticURL == "" {
		frontendStaticURL = "http://localhost:8000"
	}

	// 获取今天已经学习的单词（新单词）
	learnedWords, err := models.GetTodayLearnedWords(userID, bookIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取今日学习单词失败"})
		return
	}

	// 获取今天复习的单词
	reviewedWords, err := models.GetTodayReviewedWords(userID, bookIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取今日复习单词失败"})
		return
	}

	// 处理单词列表
	var words []map[string]interface{}

	// 处理新单词
	for _, record := range learnedWords {
		vocab := record["Vocabulary"].(*models.Vocabulary)
		isNew := record["IsNew"].(bool)

		audioURL := ""
		if vocab.AudioURL != nil && *vocab.AudioURL != "" {
			// 检查是否已经是完整的URL
			if strings.HasPrefix(*vocab.AudioURL, "http://") || strings.HasPrefix(*vocab.AudioURL, "https://") {
				audioURL = *vocab.AudioURL
			} else {
				audioURL = frontendStaticURL + *vocab.AudioURL
			}
		}

		// 处理 example_sentence 为 "null" 字符串的情况
		exampleSentence := vocab.ExampleSentence
		if exampleSentence != nil && *exampleSentence == "null" {
			emptyArray := "[]"
			exampleSentence = &emptyArray
		}

		words = append(words, map[string]interface{}{
			"id":               vocab.ID,
			"english":          vocab.English,
			"chinese":          vocab.Chinese,
			"phonetic":         vocab.Phonetic,
			"audio_url":        audioURL,
			"example_sentence": exampleSentence,
			"type":             vocab.Type,
			"is_new":           isNew,
		})
	}

	// 处理复习单词
	for _, record := range reviewedWords {
		vocab := record["Vocabulary"].(*models.Vocabulary)

		audioURL := ""
		if vocab.AudioURL != nil && *vocab.AudioURL != "" {
			// 检查是否已经是完整的URL
			if strings.HasPrefix(*vocab.AudioURL, "http://") || strings.HasPrefix(*vocab.AudioURL, "https://") {
				audioURL = *vocab.AudioURL
			} else {
				audioURL = frontendStaticURL + *vocab.AudioURL
			}
		}

		// 处理 example_sentence 为 "null" 字符串的情况
		exampleSentence := vocab.ExampleSentence
		if exampleSentence != nil && *exampleSentence == "null" {
			emptyArray := "[]"
			exampleSentence = &emptyArray
		}

		words = append(words, map[string]interface{}{
			"id":               vocab.ID,
			"english":          vocab.English,
			"chinese":          vocab.Chinese,
			"phonetic":         vocab.Phonetic,
			"audio_url":        audioURL,
			"example_sentence": exampleSentence,
			"type":             vocab.Type,
			"is_new":           false,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"words": words,
	})
}

// RecordVocabularyDictation 记录默写结果
func RecordVocabularyDictation(c *gin.Context) {
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
	isNewRecord := false
	if err != nil {
		// 记录不存在，创建新记录
		err = models.CreateLearningRecord(userID, request.WordID, request.CheckType)
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
		isNewRecord = true
	}

	fmt.Println("isNewRecord:", isNewRecord)

	// 只有默写失败时才更新艾宾浩斯记录
	if !request.IsCorrect {
		err = models.UpdateLearningRecord(record.ID, request.IsCorrect)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新学习记录失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "默写记录已更新",
	})
}
