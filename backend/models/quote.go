package models

import (
	"babyhabit/config"
	"database/sql"
	"time"
)

// Quote 名言警句结构体
type Quote struct {
	ID         int64     `json:"id"`
	Content    string    `json:"content"`
	Meaning    string    `json:"meaning"`
	Author     string    `json:"author"`
	CreateTime time.Time `json:"create_time"`
}

// CreateQuote 创建名言警句
func CreateQuote(quote *Quote) error {
	query := `INSERT INTO quote (content, meaning, author, create_time) VALUES (?, ?, ?, NOW())`
	result, err := config.DB.Exec(query, quote.Content, quote.Meaning, quote.Author)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	quote.ID = id
	return nil
}

// GetQuoteByID 根据ID获取名言警句
func GetQuoteByID(id int64) (*Quote, error) {
	query := `SELECT id, content, meaning, author, create_time FROM quote WHERE id = ?`
	quote := &Quote{}
	err := config.DB.QueryRow(query, id).Scan(&quote.ID, &quote.Content, &quote.Meaning, &quote.Author, &quote.CreateTime)
	if err != nil {
		return nil, err
	}
	return quote, nil
}

// GetAllQuotes 获取所有名言警句
func GetAllQuotes() ([]*Quote, error) {
	query := `SELECT id, content, meaning, author, create_time FROM quote ORDER BY id DESC`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quotes := []*Quote{}
	for rows.Next() {
		quote := &Quote{}
		err := rows.Scan(&quote.ID, &quote.Content, &quote.Meaning, &quote.Author, &quote.CreateTime)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

// GetQuotesWithPagination 分页获取名言警句
func GetQuotesWithPagination(page, pageSize int) ([]*Quote, int, error) {
	// 获取总数
	var total int
	countQuery := `SELECT COUNT(*) FROM quote`
	err := config.DB.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	query := `SELECT id, content, meaning, author, create_time FROM quote ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := config.DB.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	quotes := []*Quote{}
	for rows.Next() {
		quote := &Quote{}
		err := rows.Scan(&quote.ID, &quote.Content, &quote.Meaning, &quote.Author, &quote.CreateTime)
		if err != nil {
			return nil, 0, err
		}
		quotes = append(quotes, quote)
	}

	return quotes, total, nil
}

// GetRandomQuote 获取随机名言警句
func GetRandomQuote() (*Quote, error) {
	query := `SELECT id, content, meaning, author, create_time FROM quote ORDER BY RAND() LIMIT 1`
	quote := &Quote{}
	err := config.DB.QueryRow(query).Scan(&quote.ID, &quote.Content, &quote.Meaning, &quote.Author, &quote.CreateTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return quote, nil
}

// UpdateQuote 更新名言警句
func UpdateQuote(quote *Quote) error {
	query := `UPDATE quote SET content = ?, meaning = ?, author = ? WHERE id = ?`
	_, err := config.DB.Exec(query, quote.Content, quote.Meaning, quote.Author, quote.ID)
	return err
}

// DeleteQuote 删除名言警句
func DeleteQuote(id int64) error {
	query := `DELETE FROM quote WHERE id = ?`
	_, err := config.DB.Exec(query, id)
	return err
}

// BatchDeleteQuotes 批量删除名言警句
func BatchDeleteQuotes(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	query := `DELETE FROM quote WHERE id IN (`
	params := make([]interface{}, len(ids))
	for i, id := range ids {
		if i > 0 {
			query += ", "
		}
		query += "?"
		params[i] = id
	}
	query += ")"

	_, err := config.DB.Exec(query, params...)
	return err
}

// BatchCreateQuotes 批量创建名言警句（去重）
func BatchCreateQuotes(quotes []*Quote) error {
	if len(quotes) == 0 {
		return nil
	}

	// 获取已存在的名言内容，用于去重
	existingContents := make(map[string]bool)
	existingQuery := `SELECT content FROM quote`
	rows, err := config.DB.Query(existingQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return err
		}
		existingContents[content] = true
	}

	// 过滤掉重复的名言
	var uniqueQuotes []*Quote
	for _, quote := range quotes {
		if !existingContents[quote.Content] {
			uniqueQuotes = append(uniqueQuotes, quote)
			existingContents[quote.Content] = true // 避免批量导入中重复
		}
	}

	// 如果没有新名言，直接返回
	if len(uniqueQuotes) == 0 {
		return nil
	}

	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO quote (content, meaning, author, create_time) VALUES (?, ?, ?, NOW())`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, quote := range uniqueQuotes {
		_, err := stmt.Exec(quote.Content, quote.Meaning, quote.Author)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
