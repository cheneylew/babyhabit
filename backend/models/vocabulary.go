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
	Grade           *string   `json:"grade"`
	Textbook        *string   `json:"textbook"`
	Category        *string   `json:"category"`
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
	var phonetic, audioURL, exampleSentence, grade, textbook, category *string
	query := `SELECT id, english, chinese, phonetic, audio_url, example_sentence, type, grade, textbook, category, create_time 
			FROM ab_vocabulary WHERE id = ?`
	err := config.DB.QueryRow(query, id).Scan(
		&vocab.ID, &vocab.English, &vocab.Chinese, &phonetic,
		&audioURL, &exampleSentence, &vocab.Type,
		&grade, &textbook, &category, &vocab.CreateTime,
	)
	if err != nil {
		return nil, err
	}
	vocab.Phonetic = phonetic
	vocab.AudioURL = audioURL
	vocab.ExampleSentence = exampleSentence
	vocab.Grade = grade
	vocab.Textbook = textbook
	vocab.Category = category
	return &vocab, nil
}

// GetNewVocabularies 获取新词汇（用户未学习过的）
func GetNewVocabularies(userID int64, limit int) ([]*Vocabulary, error) {
	query := `
		SELECT v.id, v.english, v.chinese, v.phonetic, v.audio_url, v.example_sentence, 
			   v.type, v.grade, v.textbook, v.category, v.create_time
		FROM ab_vocabulary v
		WHERE NOT EXISTS (
			SELECT 1 FROM ab_learning_record lr WHERE lr.user_id = ? AND lr.vocabulary_id = v.id
		)
		ORDER BY v.id LIMIT ?
	`
	rows, err := config.DB.Query(query, userID, limit)
	if err != nil {
		fmt.Printf("Error querying new vocabularies: %v\n", err)
		return nil, err
	}
	fmt.Printf("Query: %s\n", query)
	fmt.Printf("UserID: %d, Limit: %d\n", userID, limit)
	defer rows.Close()

	var vocabularies []*Vocabulary
	for rows.Next() {
		var vocab Vocabulary
		var phonetic, audioURL, exampleSentence, grade, textbook, category *string
		err := rows.Scan(
			&vocab.ID, &vocab.English, &vocab.Chinese, &phonetic,
			&audioURL, &exampleSentence, &vocab.Type,
			&grade, &textbook, &category, &vocab.CreateTime,
		)
		if err != nil {
			return nil, err
		}
		vocab.Phonetic = phonetic
		vocab.AudioURL = audioURL
		vocab.ExampleSentence = exampleSentence
		vocab.Grade = grade
		vocab.Textbook = textbook
		vocab.Category = category
		vocabularies = append(vocabularies, &vocab)
	}

	return vocabularies, nil
}

// GetDueReviewVocabularies 获取需要复习的词汇
func GetDueReviewVocabularies(userID int64) ([]*LearningRecord, error) {
	query := `
		SELECT lr.id, lr.user_id, lr.vocabulary_id, lr.status, lr.review_stage, 
			   lr.next_review_date, lr.correct_count, lr.wrong_count, lr.create_time, lr.update_time,
			   v.id, v.english, v.chinese, v.phonetic, v.audio_url, v.example_sentence, 
			   v.type, v.grade, v.textbook, v.category, v.create_time
		FROM ab_learning_record lr
		JOIN ab_vocabulary v ON lr.vocabulary_id = v.id
		WHERE lr.user_id = ? AND lr.status IN ('learning', 'reviewing') 
		AND (lr.next_review_date <= CURRENT_DATE OR lr.next_review_date IS NULL)
		ORDER BY lr.next_review_date ASC
	`
	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*LearningRecord
	for rows.Next() {
		var record LearningRecord
		var vocab Vocabulary
		err := rows.Scan(
			&record.ID, &record.UserID, &record.VocabularyID, &record.Status, &record.ReviewStage,
			&record.NextReviewDate, &record.CorrectCount, &record.WrongCount, &record.CreateTime, &record.UpdateTime,
			&vocab.ID, &vocab.English, &vocab.Chinese, &vocab.Phonetic, &vocab.AudioURL, &vocab.ExampleSentence,
			&vocab.Type, &vocab.Grade, &vocab.Textbook, &vocab.Category, &vocab.CreateTime,
		)
		if err != nil {
			return nil, err
		}
		record.Vocabulary = &vocab
		records = append(records, &record)
	}

	return records, nil
}

// CreateLearningRecord 创建学习记录
func CreateLearningRecord(userID int64, vocabularyID int) error {
	query := `
		INSERT INTO ab_learning_record (user_id, vocabulary_id, status, review_stage, next_review_date)
		VALUES (?, ?, 'learning', 0, CURRENT_DATE)
	`
	_, err := config.DB.Exec(query, userID, vocabularyID)
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
		INSERT INTO ab_vocabulary (english, chinese, phonetic, audio_url, example_sentence, type, grade, textbook, category)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := config.DB.Exec(
		query,
		vocab.English, vocab.Chinese, vocab.Phonetic, vocab.AudioURL,
		vocab.ExampleSentence, vocab.Type, vocab.Grade, vocab.Textbook, vocab.Category,
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
		SELECT id, english, chinese, phonetic, audio_url, example_sentence, type, grade, textbook, category, create_time
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
		var phonetic, audioURL, exampleSentence, grade, textbook, category *string
		err := rows.Scan(
			&vocab.ID, &vocab.English, &vocab.Chinese, &phonetic,
			&audioURL, &exampleSentence, &vocab.Type,
			&grade, &textbook, &category, &vocab.CreateTime,
		)
		if err != nil {
			return nil, 0, err
		}
		vocab.Phonetic = phonetic
		vocab.AudioURL = audioURL
		vocab.ExampleSentence = exampleSentence
		vocab.Grade = grade
		vocab.Textbook = textbook
		vocab.Category = category
		vocabularies = append(vocabularies, &vocab)
	}

	return vocabularies, total, nil
}

// UpdateVocabulary 更新词汇
func UpdateVocabulary(vocab *Vocabulary) error {
	query := `
		UPDATE ab_vocabulary
		SET english = ?, chinese = ?, phonetic = ?, audio_url = ?, example_sentence = ?, 
			type = ?, grade = ?, textbook = ?, category = ?
		WHERE id = ?
	`
	_, err := config.DB.Exec(
		query,
		vocab.English, vocab.Chinese, vocab.Phonetic, vocab.AudioURL,
		vocab.ExampleSentence, vocab.Type, vocab.Grade, vocab.Textbook, vocab.Category, vocab.ID,
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
		INSERT INTO ab_vocabulary (english, chinese, phonetic, audio_url, example_sentence, type, grade, textbook, category)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, vocab := range vocabularies {
		_, err = tx.Exec(
			query,
			vocab.English, vocab.Chinese, vocab.Phonetic, vocab.AudioURL,
			vocab.ExampleSentence, vocab.Type, vocab.Grade, vocab.Textbook, vocab.Category,
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
