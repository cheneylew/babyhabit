package models

import (
	"babyhabit/config"
	"time"
)

type PointsRecord struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Points      int       `json:"points"`
	Source      string    `json:"source"`
	RelatedType string    `json:"related_type"`
	RelatedID   int64     `json:"related_id"`
	ExpireTime  time.Time `json:"expire_time"`
	CreateTime  time.Time `json:"create_time"`
}

// AddPoints 添加积分
func AddPoints(userID int64, points int, source, relatedType string, relatedID int64, expireTime time.Time) error {
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

	// 插入积分记录
	var query string
	var args []interface{}
	if expireTime.IsZero() {
		// 如果没有过期时间，设置为 NULL
		query = `INSERT INTO points_record (user_id, points, source, related_type, related_id, expire_time, create_time) 
				 VALUES (?, ?, ?, ?, ?, NULL, NOW())`
		args = []interface{}{userID, points, source, relatedType, relatedID}
	} else {
		query = `INSERT INTO points_record (user_id, points, source, related_type, related_id, expire_time, create_time) 
				 VALUES (?, ?, ?, ?, ?, ?, NOW())`
		args = []interface{}{userID, points, source, relatedType, relatedID, expireTime}
	}
	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	// 更新用户积分余额
	updateQuery := `UPDATE user SET points_balance = points_balance + ? WHERE id = ?`
	_, err = tx.Exec(updateQuery, points, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeductPoints 扣除积分
func DeductPoints(userID int64, points int, source, relatedType string, relatedID int64) error {
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

	// 检查积分余额
	var balance int
	checkQuery := `SELECT points_balance FROM user WHERE id = ?`
	err = tx.QueryRow(checkQuery, userID).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < points {
		return ErrInsufficientPoints
	}

	// 插入积分记录（负数）
	query := `INSERT INTO points_record (user_id, points, source, related_type, related_id, create_time) 
			 VALUES (?, ?, ?, ?, ?, NOW())`
	_, err = tx.Exec(query, userID, -points, source, relatedType, relatedID)
	if err != nil {
		return err
	}

	// 更新用户积分余额
	updateQuery := `UPDATE user SET points_balance = points_balance - ? WHERE id = ?`
	_, err = tx.Exec(updateQuery, points, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetPointsRecords 获取用户积分记录
func GetPointsRecords(userID int64, limit, offset int) ([]*PointsRecord, error) {
	query := `SELECT id, user_id, points, source, related_type, related_id, expire_time, create_time 
			 FROM points_record 
			 WHERE user_id = ? 
			 ORDER BY create_time DESC 
			 LIMIT ? OFFSET ?`
	rows, err := config.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []*PointsRecord{}
	for rows.Next() {
		record := &PointsRecord{}
		err := rows.Scan(
			&record.ID, &record.UserID, &record.Points, &record.Source,
			&record.RelatedType, &record.RelatedID, &record.ExpireTime, &record.CreateTime,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

// GetUserPointsBalance 获取用户积分余额
func GetUserPointsBalance(userID int64) (int, error) {
	var balance int
	query := `SELECT points_balance FROM user WHERE id = ?`
	err := config.DB.QueryRow(query, userID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// Errors
var (
	ErrInsufficientPoints = NewError("insufficient points")
)

// Error 自定义错误
func NewError(message string) error {
	return &customError{message: message}
}

type customError struct {
	message string
}

func (e *customError) Error() string {
	return e.message
}
