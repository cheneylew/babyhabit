package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"babyhabit/config"
)

// Vocabulary 词汇模型
type Vocabulary struct {
	ID              int       `json:"id"`
	English         string    `json:"english"`
	Chinese         string    `json:"chinese"`
	Phonetic        *string   `json:"phonetic"`
	AudioURL        *string   `json:"audio_url"`
	ExampleSentence *string   `json:"example_sentence"`
	Type            string    `json:"type"`
	BookID          *int      `json:"book_id"`
	Category        *string   `json:"category"`
	Remark          *string   `json:"remark"`
	CreateTime      time.Time `json:"create_time"`
}

// LearningRecord 学习记录模型
type LearningRecord struct {
	ID             int         `json:"id"`
	UserID         int64       `json:"user_id"`
	VocabularyID   int         `json:"vocabulary_id"`
	Status         string      `json:"status"`
	ReviewStage    int         `json:"review_stage"`
	NextReviewDate time.Time   `json:"next_review_date"`
	CorrectCount   int         `json:"correct_count"`
	WrongCount     int         `json:"wrong_count"`
	CreateTime     time.Time   `json:"create_time"`
	UpdateTime     time.Time   `json:"update_time"`
	Vocabulary     *Vocabulary `json:"vocabulary,omitempty"`
}

// StudyCheckin 学习打卡模型
type StudyCheckin struct {
	ID               int       `json:"id"`
	UserID           int64     `json:"user_id"`
	CheckinDate      time.Time `json:"checkin_date"`
	NewWordsCount    int       `json:"new_words_count"`
	ReviewWordsCount int       `json:"review_words_count"`
	CorrectCount     int       `json:"correct_count"`
	TotalCount       int       `json:"total_count"`
	DurationMinutes  int       `json:"duration_minutes"`
	CreateTime       time.Time `json:"create_time"`
}

// GetVocabularyByID 根据ID获取词汇
func GetVocabularyByID(id int) (*Vocabulary, error) {
	var vocab Vocabulary
	var phonetic, audioURL, exampleSentence, category, remark *string
	var bookID *int
	query := `SELECT id, english, chinese, phonetic, audio_url, example_sentence, type, book_id, category, remark, create_time 
			FROM ab_vocabulary WHERE id = ?`
	err := config.DB.QueryRow(query, id).Scan(
		&vocab.ID, &vocab.English, &vocab.Chinese, &phonetic,
		&audioURL, &exampleSentence, &vocab.Type,
		&bookID, &category, &remark, &vocab.CreateTime,
	)
	if err != nil {
		return nil, err
	}
	vocab.Phonetic = phonetic
	vocab.AudioURL = audioURL
	vocab.ExampleSentence = exampleSentence
	vocab.BookID = bookID
	vocab.Category = category
	vocab.Remark = remark
	return &vocab, nil
}

// GetNewVocabularies 获取新词汇（用户未学习过的）
func GetNewVocabularies(userID int64, limit int, bookIDs []int) ([]*Vocabulary, error) {
	// 如果没有选中教材，返回空
	if len(bookIDs) == 0 {
		return []*Vocabulary{}, nil
	}

	// 计算每个教材应该分配的单词数量
	bookCount := len(bookIDs)
	baseCount := limit / bookCount
	remainder := limit % bookCount

	var vocabularies []*Vocabulary
	var availableBooks []int

	// 第一次遍历：为每个教材获取分配的单词数量
	for i, bookID := range bookIDs {
		// 计算当前教材应该获取的单词数量
		currentLimit := baseCount
		if i < remainder {
			currentLimit++
		}

		// 构建查询
		query := `
			SELECT v.id, v.english, v.chinese, v.phonetic, v.audio_url, v.example_sentence, 
			   v.type, v.book_id, v.category, v.remark, v.create_time
			FROM ab_vocabulary v
			WHERE NOT EXISTS (
				SELECT 1 FROM ab_learning_record lr WHERE lr.user_id = ? AND lr.vocabulary_id = v.id
			) AND v.book_id = ?
			ORDER BY v.id LIMIT ?
		`

		// 执行查询
		rows, err := config.DB.Query(query, userID, bookID, currentLimit)
		if err != nil {
			fmt.Printf("Error querying new vocabularies for book %d: %v\n", bookID, err)
			return nil, err
		}

		// 处理结果
		bookVocabularies := []*Vocabulary{}
		for rows.Next() {
			var vocab Vocabulary
			var phonetic, audioURL, exampleSentence, category, remark *string
			var bookID *int
			err := rows.Scan(
				&vocab.ID, &vocab.English, &vocab.Chinese, &phonetic,
				&audioURL, &exampleSentence, &vocab.Type,
				&bookID, &category, &remark, &vocab.CreateTime,
			)
			if err != nil {
				rows.Close()
				return nil, err
			}
			vocab.Phonetic = phonetic
			vocab.AudioURL = audioURL
			vocab.ExampleSentence = exampleSentence
			vocab.BookID = bookID
			vocab.Category = category
			vocab.Remark = remark
			bookVocabularies = append(bookVocabularies, &vocab)
		}

		rows.Close()

		// 添加到总列表
		vocabularies = append(vocabularies, bookVocabularies...)

		// 检查当前教材是否还有未学习的单词
		checkQuery := `
			SELECT COUNT(*)
			FROM ab_vocabulary v
			WHERE NOT EXISTS (
				SELECT 1 FROM ab_learning_record lr WHERE lr.user_id = ? AND lr.vocabulary_id = v.id
			) AND v.book_id = ?
		`
		var count int
		var checkErr error
		checkErr = config.DB.QueryRow(checkQuery, userID, bookID).Scan(&count)
		if checkErr == nil && count > 0 {
			// 如果当前教材还有未学习的单词，添加到可用教材列表
			availableBooks = append(availableBooks, bookID)
		}
	}

	// 检查是否需要补充单词
	if len(vocabularies) < limit && len(availableBooks) > 0 {
		// 计算还需要补充的单词数量
		remaining := limit - len(vocabularies)

		// 从可用教材中补充单词
		for _, bookID := range availableBooks {
			if remaining <= 0 {
				break
			}

			// 构建查询，获取当前教材中已经学过的单词ID，排除已经在列表中的单词
			query := `
				SELECT v.id, v.english, v.chinese, v.phonetic, v.audio_url, v.example_sentence, 
				   v.type, v.book_id, v.category, v.remark, v.create_time
				FROM ab_vocabulary v
				WHERE NOT EXISTS (
					SELECT 1 FROM ab_learning_record lr WHERE lr.user_id = ? AND lr.vocabulary_id = v.id
				) AND v.book_id = ?
			`

			// 添加排除已经在列表中的单词的条件
			args := []interface{}{userID, bookID}
			if len(vocabularies) > 0 {
				query += " AND v.id NOT IN ("
				for i, vocab := range vocabularies {
					if i > 0 {
						query += ","
					}
					query += "?"
					args = append(args, vocab.ID)
				}
				query += ")"
			}

			// 添加排序和限制
			query += " ORDER BY v.id LIMIT ?"
			args = append(args, remaining)

			// 执行查询
			rows, err := config.DB.Query(query, args...)
			if err != nil {
				fmt.Printf("Error querying additional vocabularies for book %d: %v\n", bookID, err)
				return nil, err
			}

			// 处理结果
			addedCount := 0
			for rows.Next() {
				var vocab Vocabulary
				var phonetic, audioURL, exampleSentence, category, remark *string
				var bookID *int
				err := rows.Scan(
					&vocab.ID, &vocab.English, &vocab.Chinese, &phonetic,
					&audioURL, &exampleSentence, &vocab.Type,
					&bookID, &category, &remark, &vocab.CreateTime,
				)
				if err != nil {
					rows.Close()
					return nil, err
				}
				vocab.Phonetic = phonetic
				vocab.AudioURL = audioURL
				vocab.ExampleSentence = exampleSentence
				vocab.BookID = bookID
				vocab.Category = category
				vocab.Remark = remark

				// 检查是否已经存在相同的单词
				exists := false
				for _, v := range vocabularies {
					if v.ID == vocab.ID {
						exists = true
						break
					}
				}

				// 如果不存在，添加到列表
				if !exists {
					vocabularies = append(vocabularies, &vocab)
					addedCount++
					remaining--
					if remaining <= 0 {
						break
					}
				}
			}

			rows.Close()
		}
	}

	return vocabularies, nil
}

// GetTodayLearnedNewWordsCount 获取用户今天已经学习的新单词数量
func GetTodayLearnedNewWordsCount(userID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM ab_learning_record
		WHERE user_id = ? AND DATE(create_time) = DATE(NOW()) AND (check_type != 'skip' OR check_type IS NULL)
	`
	var count int
	err := config.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		fmt.Printf("Error querying today learned new words: %v\n", err)
		return 0, err
	}
	return count, nil
}

// GetDueReviewVocabularies 获取需要复习的词汇
func GetDueReviewVocabularies(userID int64, bookIDs []int) ([]*LearningRecord, error) {
	query := `
		SELECT lr.id, lr.user_id, lr.vocabulary_id, lr.status, lr.review_stage, 
		   lr.next_review_date, lr.correct_count, lr.wrong_count, lr.create_time, lr.update_time,
		   v.id, v.english, v.chinese, v.phonetic, v.audio_url, v.example_sentence, 
		   v.type, v.book_id, v.category, v.create_time
		FROM ab_learning_record lr
		JOIN ab_vocabulary v ON lr.vocabulary_id = v.id
		WHERE lr.user_id = ? AND lr.status IN ('learning', 'reviewing') 
		AND (lr.next_review_date <= CURRENT_DATE OR lr.next_review_date IS NULL)
	`

	// 添加教材ID过滤
	args := []interface{}{userID}
	if len(bookIDs) > 0 {
		query += " AND v.book_id IN (?"
		args = append(args, bookIDs[0])
		for i := 1; i < len(bookIDs); i++ {
			query += ", ?"
			args = append(args, bookIDs[i])
		}
		query += ")"
	}

	query += " ORDER BY lr.next_review_date ASC"

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*LearningRecord
	for rows.Next() {
		var record LearningRecord
		var vocab Vocabulary
		var bookID *int
		err := rows.Scan(
			&record.ID, &record.UserID, &record.VocabularyID, &record.Status, &record.ReviewStage,
			&record.NextReviewDate, &record.CorrectCount, &record.WrongCount, &record.CreateTime, &record.UpdateTime,
			&vocab.ID, &vocab.English, &vocab.Chinese, &vocab.Phonetic, &vocab.AudioURL, &vocab.ExampleSentence,
			&vocab.Type, &bookID, &vocab.Category, &vocab.CreateTime,
		)
		if err != nil {
			return nil, err
		}
		vocab.BookID = bookID
		record.Vocabulary = &vocab
		records = append(records, &record)
	}

	return records, nil
}

// CreateLearningRecord 创建学习记录
func CreateLearningRecord(userID int64, vocabularyID int, checkType string) error {
	query := `
		INSERT INTO ab_learning_record (user_id, vocabulary_id, status, check_type, review_stage, next_review_date)
		VALUES (?, ?, 'learning', ?, 0, CURRENT_DATE)
	`
	_, err := config.DB.Exec(query, userID, vocabularyID, checkType)
	return err
}

// UpdateLearningRecord 更新学习记录
func UpdateLearningRecord(recordID int, isCorrect bool) error {
	// 开始事务
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 获取当前记录
	var record LearningRecord
	query := `SELECT id, user_id, vocabulary_id, status, review_stage, correct_count, wrong_count FROM ab_learning_record WHERE id = ?`
	err = tx.QueryRow(query, recordID).Scan(
		&record.ID, &record.UserID, &record.VocabularyID, &record.Status, &record.ReviewStage,
		&record.CorrectCount, &record.WrongCount,
	)
	if err != nil {
		return err
	}

	// 更新记录
	if isCorrect {
		record.CorrectCount++
		record.ReviewStage++

		if record.ReviewStage >= 6 {
			// 已掌握
			_, err = tx.Exec(
				`UPDATE ab_learning_record SET status = 'mastered', correct_count = ?, review_stage = ? WHERE id = ?`,
				record.CorrectCount, record.ReviewStage, recordID,
			)
		} else {
			// 计算下次复习时间
			intervals := []int{1, 1, 3, 7, 15, 30}
			interval := intervals[record.ReviewStage-1]
			nextReviewDate := time.Now().AddDate(0, 0, interval)

			_, err = tx.Exec(
				`UPDATE ab_learning_record SET status = 'reviewing', correct_count = ?, review_stage = ?, next_review_date = ? WHERE id = ?`,
				record.CorrectCount, record.ReviewStage, nextReviewDate, recordID,
			)
		}
	} else {
		record.WrongCount++
		record.ReviewStage = 0
		// 错误后第二天复习
		nextReviewDate := time.Now().AddDate(0, 0, 1)

		_, err = tx.Exec(
			`UPDATE ab_learning_record SET status = 'reviewing', wrong_count = ?, review_stage = 0, next_review_date = ? WHERE id = ?`,
			record.WrongCount, nextReviewDate, recordID,
		)
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}

// MarkVocabularyAsMastered 直接标记单词为掌握
func MarkVocabularyAsMastered(recordID int, checkType string) error {
	// 开始事务
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 直接将单词标记为掌握
	_, err = tx.Exec(
		`UPDATE ab_learning_record SET status = 'mastered', check_type = ?, correct_count = 1, review_stage = 6 WHERE id = ?`,
		checkType, recordID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetLearningRecordByUserAndVocab 根据用户ID和词汇ID获取学习记录
func GetLearningRecordByUserAndVocab(userID int64, vocabularyID int) (*LearningRecord, error) {
	var record LearningRecord
	query := `
		SELECT id, user_id, vocabulary_id, status, review_stage, next_review_date, 
			   correct_count, wrong_count, create_time, update_time
		FROM ab_learning_record 
		WHERE user_id = ? AND vocabulary_id = ?
	`
	err := config.DB.QueryRow(query, userID, vocabularyID).Scan(
		&record.ID, &record.UserID, &record.VocabularyID, &record.Status, &record.ReviewStage,
		&record.NextReviewDate, &record.CorrectCount, &record.WrongCount, &record.CreateTime, &record.UpdateTime,
	)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// GetLearningStats 获取学习统计
func GetLearningStats(userID int64) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 累计学习单词数
	var totalWords int
	query := `SELECT COUNT(*) FROM ab_learning_record WHERE user_id = ?`
	err := config.DB.QueryRow(query, userID).Scan(&totalWords)
	if err != nil {
		return nil, err
	}
	stats["totalWords"] = totalWords

	// 已掌握单词数
	var masteredWords int
	query = `SELECT COUNT(*) FROM ab_learning_record WHERE user_id = ? AND status = 'mastered'`
	err = config.DB.QueryRow(query, userID).Scan(&masteredWords)
	if err != nil {
		return nil, err
	}
	stats["masteredWords"] = masteredWords

	// 今日学习单词数
	var todayLearnedWords int
	query = `SELECT COUNT(*) FROM ab_learning_record WHERE user_id = ? AND DATE(create_time) = CURRENT_DATE AND (check_type != 'skip' OR check_type IS NULL)`
	err = config.DB.QueryRow(query, userID).Scan(&todayLearnedWords)
	if err != nil {
		todayLearnedWords = 0
	}
	stats["todayLearnedWords"] = todayLearnedWords

	// 今日学习新单词数
	var todayLearnedNewWords int
	query = `SELECT COUNT(*) FROM ab_learning_record WHERE user_id = ? AND DATE(create_time) = CURRENT_DATE AND (check_type != 'skip' OR check_type IS NULL)`
	err = config.DB.QueryRow(query, userID).Scan(&todayLearnedNewWords)
	if err != nil {
		todayLearnedNewWords = 0
	}
	stats["todayLearnedNewWords"] = todayLearnedNewWords

	// 今日复习单词数
	var todayReviewedWords int
	query = `SELECT COUNT(*) FROM ab_learning_record WHERE user_id = ? AND DATE(update_time) = CURRENT_DATE AND DATE(create_time) != CURRENT_DATE`
	err = config.DB.QueryRow(query, userID).Scan(&todayReviewedWords)
	if err != nil {
		todayReviewedWords = 0
	}
	stats["todayReviewedWords"] = todayReviewedWords

	// 连续学习天数
	var learningStreak int
	query = `
		SELECT COUNT(*)
		FROM ab_study_checkin
		WHERE user_id = ?
		AND checkin_date >= DATE_SUB(CURRENT_DATE, INTERVAL (SELECT COUNT(*) FROM ab_study_checkin WHERE user_id = ? AND checkin_date >= DATE_SUB(CURRENT_DATE, INTERVAL 365 DAY)) DAY)
	`
	err = config.DB.QueryRow(query, userID, userID).Scan(&learningStreak)
	if err != nil {
		learningStreak = 0
	}
	stats["learningStreak"] = learningStreak

	// 正确率
	var correctCount, totalCount int
	query = `
		SELECT SUM(correct_count), SUM(total_count)
		FROM ab_study_checkin
		WHERE user_id = ?
		AND checkin_date >= DATE_SUB(CURRENT_DATE, INTERVAL 30 DAY)
	`
	err = config.DB.QueryRow(query, userID).Scan(&correctCount, &totalCount)
	if err != nil {
		correctCount = 0
		totalCount = 0
	}

	accuracyRate := 0.0
	if totalCount > 0 {
		accuracyRate = float64(correctCount) / float64(totalCount) * 100
	}
	stats["accuracyRate"] = fmt.Sprintf("%.1f", accuracyRate)

	return stats, nil
}

// CreateStudyCheckin 创建学习打卡记录
func CreateStudyCheckin(userID int64, newWordsCount, reviewWordsCount, correctCount, totalCount, durationMinutes int) error {
	query := `
		INSERT INTO ab_study_checkin (user_id, checkin_date, new_words_count, review_words_count, correct_count, total_count, duration_minutes)
		VALUES (?, CURRENT_DATE, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
		new_words_count = new_words_count + ?, 
		review_words_count = review_words_count + ?, 
		correct_count = correct_count + ?, 
		total_count = total_count + ?, 
		duration_minutes = duration_minutes + ?
	`
	_, err := config.DB.Exec(
		query, userID, newWordsCount, reviewWordsCount, correctCount, totalCount, durationMinutes,
		newWordsCount, reviewWordsCount, correctCount, totalCount, durationMinutes,
	)
	return err
}

// GetVocabularyOptions 获取词汇选项（用于选择题）
func GetVocabularyOptions(vocabularyID int, optionType string) ([]string, error) {
	var options []string
	var query string

	switch optionType {
	case "chineseToEnglish":
		// 获取中文对应的英文选项
		query = `
			SELECT english FROM ab_vocabulary 
			WHERE id != ? 
			ORDER BY RAND() LIMIT 3
		`
		rows, err := config.DB.Query(query, vocabularyID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var english string
			err := rows.Scan(&english)
			if err != nil {
				return nil, err
			}
			options = append(options, english)
		}

		// 添加正确答案
		var correctEnglish string
		err = config.DB.QueryRow("SELECT english FROM ab_vocabulary WHERE id = ?", vocabularyID).Scan(&correctEnglish)
		if err != nil {
			return nil, err
		}
		options = append(options, correctEnglish)

	case "englishToChinese":
		// 获取英文对应的中文选项
		query = `
			SELECT chinese FROM ab_vocabulary 
			WHERE id != ? 
			ORDER BY RAND() LIMIT 3
		`
		rows, err := config.DB.Query(query, vocabularyID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var chinese string
			err := rows.Scan(&chinese)
			if err != nil {
				return nil, err
			}
			options = append(options, chinese)
		}

		// 添加正确答案
		var correctChinese string
		err = config.DB.QueryRow("SELECT chinese FROM ab_vocabulary WHERE id = ?", vocabularyID).Scan(&correctChinese)
		if err != nil {
			return nil, err
		}
		options = append(options, correctChinese)

	case "listening":
		// 听力选项，返回英文单词
		query = `
			SELECT english FROM ab_vocabulary 
			WHERE id != ? 
			ORDER BY RAND() LIMIT 3
		`
		rows, err := config.DB.Query(query, vocabularyID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var english string
			err := rows.Scan(&english)
			if err != nil {
				return nil, err
			}
			options = append(options, english)
		}

		// 添加正确答案
		var correctEnglish string
		err = config.DB.QueryRow("SELECT english FROM ab_vocabulary WHERE id = ?", vocabularyID).Scan(&correctEnglish)
		if err != nil {
			return nil, err
		}
		options = append(options, correctEnglish)
	}

	// 打乱选项顺序
	for i := len(options) - 1; i > 0; i-- {
		j := time.Now().UnixNano() % int64(i+1)
		options[i], options[int(j)] = options[int(j)], options[i]
	}

	return options, nil
}

// CreateVocabulary 创建词汇
func CreateVocabulary(vocab *Vocabulary) error {
	query := `
		INSERT INTO ab_vocabulary (english, chinese, phonetic, audio_url, example_sentence, type, book_id, category, remark)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := config.DB.Exec(
		query,
		vocab.English, vocab.Chinese, vocab.Phonetic, vocab.AudioURL,
		vocab.ExampleSentence, vocab.Type, vocab.BookID, vocab.Category, vocab.Remark,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	vocab.ID = int(id)
	return nil
}

// GetLearningHistory 获取用户学习历史
func GetLearningHistory(userID int64) ([]*Vocabulary, error) {
	query := `
		SELECT v.id, v.english, v.chinese, v.phonetic, v.audio_url, v.example_sentence, 
		   v.type, v.book_id, v.category, v.remark, v.create_time
		FROM ab_vocabulary v
		JOIN ab_learning_record lr ON v.id = lr.vocabulary_id
		WHERE lr.user_id = ?
		ORDER BY lr.create_time DESC
		LIMIT 100
	`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vocabularies []*Vocabulary
	for rows.Next() {
		var vocab Vocabulary
		var phonetic, audioURL, exampleSentence, category, remark *string
		var bookID *int
		err := rows.Scan(
			&vocab.ID, &vocab.English, &vocab.Chinese, &phonetic,
			&audioURL, &exampleSentence, &vocab.Type,
			&bookID, &category, &remark, &vocab.CreateTime,
		)
		if err != nil {
			return nil, err
		}
		vocab.Phonetic = phonetic
		vocab.AudioURL = audioURL
		vocab.ExampleSentence = exampleSentence
		vocab.BookID = bookID
		vocab.Category = category
		vocab.Remark = remark
		vocabularies = append(vocabularies, &vocab)
	}

	return vocabularies, nil
}

// GetVocabularies 获取词汇列表（带分页）
func GetVocabularies(page, pageSize int) ([]*Vocabulary, int, error) {
	var vocabularies []*Vocabulary
	var total int

	// 获取总数
	countQuery := "SELECT COUNT(*) FROM ab_vocabulary"
	err := config.DB.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	query := `
		SELECT id, english, chinese, phonetic, audio_url, example_sentence, type, book_id, category, remark, create_time
		FROM ab_vocabulary
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`
	rows, err := config.DB.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var vocab Vocabulary
		var phonetic, audioURL, exampleSentence, category, remark *string
		var bookID *int
		err := rows.Scan(
			&vocab.ID, &vocab.English, &vocab.Chinese, &phonetic,
			&audioURL, &exampleSentence, &vocab.Type, &bookID, &category, &remark, &vocab.CreateTime,
		)
		if err != nil {
			return nil, 0, err
		}
		vocab.Phonetic = phonetic
		vocab.AudioURL = audioURL
		vocab.ExampleSentence = exampleSentence
		vocab.BookID = bookID
		vocab.Category = category
		vocab.Remark = remark
		vocabularies = append(vocabularies, &vocab)
	}

	return vocabularies, total, nil
}

// UpdateVocabulary 更新词汇
func UpdateVocabulary(vocab *Vocabulary) error {
	query := `
		UPDATE ab_vocabulary
		SET english = ?, chinese = ?, phonetic = ?, audio_url = ?, example_sentence = ?, 
			type = ?, book_id = ?, category = ?
		WHERE id = ?
	`
	_, err := config.DB.Exec(
		query,
		vocab.English, vocab.Chinese, vocab.Phonetic, vocab.AudioURL,
		vocab.ExampleSentence, vocab.Type, vocab.BookID, vocab.Category, vocab.ID,
	)
	return err
}

// DeleteVocabulary 删除词汇
func DeleteVocabulary(id int) error {
	// 开始事务
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 查询词汇的音频URL
	var audioURL *string
	err = tx.QueryRow("SELECT audio_url FROM ab_vocabulary WHERE id = ?", id).Scan(&audioURL)
	if err == nil && audioURL != nil && *audioURL != "" {
		// 提取文件名并删除文件
		filePath := filepath.Join("files", "words", filepath.Base(*audioURL))
		os.Remove(filePath)
	}

	// 删除相关的学习记录
	_, err = tx.Exec("DELETE FROM ab_learning_record WHERE vocabulary_id = ?", id)
	if err != nil {
		return err
	}

	// 删除词汇
	_, err = tx.Exec("DELETE FROM ab_vocabulary WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetVocabulariesByBookID 根据教材ID获取词汇列表
func GetVocabulariesByBookID(bookID int) ([]*Vocabulary, error) {
	query := `SELECT id, english FROM ab_vocabulary WHERE book_id = ?`
	rows, err := config.DB.Query(query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vocabularies []*Vocabulary
	for rows.Next() {
		var vocab Vocabulary
		if err := rows.Scan(&vocab.ID, &vocab.English); err != nil {
			return nil, err
		}
		vocabularies = append(vocabularies, &vocab)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vocabularies, nil
}

// BatchCreateVocabulary 批量创建词汇
func BatchCreateVocabulary(vocabularies []*Vocabulary) error {
	// 开始事务
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO ab_vocabulary (english, chinese, phonetic, audio_url, example_sentence, type, book_id, category)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, vocab := range vocabularies {
		_, err = tx.Exec(
			query,
			vocab.English, vocab.Chinese, vocab.Phonetic, vocab.AudioURL,
			vocab.ExampleSentence, vocab.Type, vocab.BookID, vocab.Category,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// BatchDeleteVocabulary 批量删除词汇
func BatchDeleteVocabulary(ids []int) error {
	// 开始事务
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 为IN子句生成占位符
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	// 查询所有要删除的词汇的音频URL
	query := "SELECT audio_url FROM ab_vocabulary WHERE id IN (" + strings.Join(placeholders, ",") + ")"
	rows, err := tx.Query(query, args...)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var audioURL *string
			if err := rows.Scan(&audioURL); err == nil && audioURL != nil && *audioURL != "" {
				// 提取文件名并删除文件
				filePath := filepath.Join("files", "words", filepath.Base(*audioURL))
				os.Remove(filePath)
			}
		}
	}

	// 删除相关的学习记录
	query = "DELETE FROM ab_learning_record WHERE vocabulary_id IN (" + strings.Join(placeholders, ",") + ")"
	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	// 删除词汇
	query = "DELETE FROM ab_vocabulary WHERE id IN (" + strings.Join(placeholders, ",") + ")"
	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetBookOptions 获取教材选项列表
func GetBookOptions() ([]map[string]interface{}, error) {
	// 获取所有教材
	bookQuery := `SELECT id, name FROM ab_book ORDER BY name`
	bookRows, err := config.DB.Query(bookQuery)
	if err != nil {
		return nil, err
	}
	defer bookRows.Close()

	var books []map[string]interface{}
	for bookRows.Next() {
		var id int
		var name string
		if err := bookRows.Scan(&id, &name); err != nil {
			return nil, err
		}
		books = append(books, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	return books, nil
}

// Book 教材模型
type Book struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	CreateTime time.Time `json:"create_time"`
}

// GetBooks 获取所有教材列表
func GetBooks() ([]*Book, error) {
	query := `SELECT id, name, create_time FROM ab_book ORDER BY name`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Name, &book.CreateTime)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, nil
}

// GetBookByID 根据ID获取教材
func GetBookByID(id int) (*Book, error) {
	var book Book
	query := `SELECT id, name, create_time FROM ab_book WHERE id = ?`
	err := config.DB.QueryRow(query, id).Scan(&book.ID, &book.Name, &book.CreateTime)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// CreateBook 创建教材
func CreateBook(book *Book) error {
	query := `INSERT INTO ab_book (name) VALUES (?)`
	result, err := config.DB.Exec(query, book.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	book.ID = int(id)
	return nil
}

// UpdateBook 更新教材
func UpdateBook(book *Book) error {
	query := `UPDATE ab_book SET name = ? WHERE id = ?`
	_, err := config.DB.Exec(query, book.Name, book.ID)
	return err
}

// DeleteBook 删除教材
func DeleteBook(id int) error {
	query := `DELETE FROM ab_book WHERE id = ?`
	_, err := config.DB.Exec(query, id)
	return err
}

// GetTodayLearnedWords 获取今日学习的单词
func GetTodayLearnedWords(userID int64, bookIDs []int) ([]map[string]interface{}, error) {
	query := `
		SELECT lr.id, lr.user_id, lr.vocabulary_id, lr.status, lr.review_stage, 
		   lr.next_review_date, lr.correct_count, lr.wrong_count, lr.create_time, lr.update_time,
		   v.id, v.english, v.chinese, v.phonetic, v.audio_url, v.example_sentence, 
		   v.type, v.book_id, v.category, v.create_time,
		   DATE(lr.create_time) = DATE(NOW()) AS is_new
		FROM ab_learning_record lr
		JOIN ab_vocabulary v ON lr.vocabulary_id = v.id
		WHERE lr.user_id = ? AND DATE(lr.create_time) = DATE(NOW()) AND (lr.check_type != 'skip' OR lr.check_type IS NULL)
	`

	// 添加教材ID过滤
	args := []interface{}{userID}
	if len(bookIDs) > 0 {
		query += " AND v.book_id IN (?"
		args = append(args, bookIDs[0])
		for i := 1; i < len(bookIDs); i++ {
			query += ", ?"
			args = append(args, bookIDs[i])
		}
		query += ")"
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []map[string]interface{}
	for rows.Next() {
		var record LearningRecord
		var vocab Vocabulary
		var bookID *int
		var isNew bool
		err := rows.Scan(
			&record.ID, &record.UserID, &record.VocabularyID, &record.Status, &record.ReviewStage,
			&record.NextReviewDate, &record.CorrectCount, &record.WrongCount, &record.CreateTime, &record.UpdateTime,
			&vocab.ID, &vocab.English, &vocab.Chinese, &vocab.Phonetic, &vocab.AudioURL, &vocab.ExampleSentence,
			&vocab.Type, &bookID, &vocab.Category, &vocab.CreateTime,
			&isNew,
		)
		if err != nil {
			return nil, err
		}
		vocab.BookID = bookID
		record.Vocabulary = &vocab
		records = append(records, map[string]interface{}{
			"ID":         record.ID,
			"UserID":     record.UserID,
			"Vocabulary": record.Vocabulary,
			"IsNew":      isNew,
		})
	}

	return records, nil
}

// GetTodayReviewedWords 获取今日复习的单词
func GetTodayReviewedWords(userID int64, bookIDs []int) ([]map[string]interface{}, error) {
	query := `
		SELECT lr.id, lr.user_id, lr.vocabulary_id, lr.status, lr.review_stage, 
		   lr.next_review_date, lr.correct_count, lr.wrong_count, lr.create_time, lr.update_time,
		   v.id, v.english, v.chinese, v.phonetic, v.audio_url, v.example_sentence, 
		   v.type, v.book_id, v.category, v.create_time
		FROM ab_learning_record lr
		JOIN ab_vocabulary v ON lr.vocabulary_id = v.id
		WHERE lr.user_id = ? AND DATE(lr.update_time) = DATE(NOW()) AND DATE(lr.create_time) != DATE(NOW()) AND (lr.check_type != 'skip' OR lr.check_type IS NULL)
	`

	// 添加教材ID过滤
	args := []interface{}{userID}
	if len(bookIDs) > 0 {
		query += " AND v.book_id IN (?"
		args = append(args, bookIDs[0])
		for i := 1; i < len(bookIDs); i++ {
			query += ", ?"
			args = append(args, bookIDs[i])
		}
		query += ")"
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []map[string]interface{}
	for rows.Next() {
		var record LearningRecord
		var vocab Vocabulary
		var bookID *int
		err := rows.Scan(
			&record.ID, &record.UserID, &record.VocabularyID, &record.Status, &record.ReviewStage,
			&record.NextReviewDate, &record.CorrectCount, &record.WrongCount, &record.CreateTime, &record.UpdateTime,
			&vocab.ID, &vocab.English, &vocab.Chinese, &vocab.Phonetic, &vocab.AudioURL, &vocab.ExampleSentence,
			&vocab.Type, &bookID, &vocab.Category, &vocab.CreateTime,
		)
		if err != nil {
			return nil, err
		}
		vocab.BookID = bookID
		record.Vocabulary = &vocab
		records = append(records, map[string]interface{}{
			"ID":         record.ID,
			"UserID":     record.UserID,
			"Vocabulary": record.Vocabulary,
		})
	}

	return records, nil
}

// UpdateLearningRecordForReview 更新学习记录为需要复习
func UpdateLearningRecordForReview(recordID int) error {
	// 错误后第二天复习
	nextReviewDate := time.Now().AddDate(0, 0, 1)

	query := `
		UPDATE ab_learning_record 
		SET status = 'reviewing', review_stage = 0, next_review_date = ? 
		WHERE id = ?
	`
	_, err := config.DB.Exec(query, nextReviewDate, recordID)
	return err
}
