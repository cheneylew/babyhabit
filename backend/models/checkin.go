package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/cheneylew/babyhabit/backend/config"
)

type CheckinRecord struct {
	ID                int64          `json:"id"`
	UserID            int64          `json:"user_id"`
	HabitID           int64          `json:"habit_id"`
	HabitName         string         `json:"habit_name"`
	AssignmentID      int64          `json:"assignment_id"`
	CheckinDate       string         `json:"checkin_date"`
	CheckinTime       time.Time      `json:"checkin_time"`
	CheckinType       int            `json:"checkin_type"`
	PhotoURL          string         `json:"photo_url"`
	SelfRate          int            `json:"self_rate"`
	Remark            string         `json:"remark"`
	Status            int            `json:"status"`
	PointsRewarded    int            `json:"points_rewarded"`
	IsRolledBack      int            `json:"is_rolled_back"`
	RollbackTime      *time.Time     `json:"rollback_time"`
	RollbackReason    sql.NullString `json:"-"`
	RollbackReasonStr string         `json:"rollback_reason"`
	CreateTime        time.Time      `json:"create_time"`
}

type StreakRecord struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	HabitID         int64     `json:"habit_id"`
	CurrentStreak   int       `json:"current_streak"`
	LongestStreak   int       `json:"longest_streak"`
	LastCheckinDate string    `json:"last_checkin_date"`
	StreakStartDate string    `json:"streak_start_date"`
	UpdateTime      time.Time `json:"update_time"`
}

// CreateCheckinRecord 创建打卡记录
func CreateCheckinRecord(record *CheckinRecord) error {
	query := `INSERT INTO checkin_record (user_id, habit_id, assignment_id, checkin_date, 
			 checkin_time, checkin_type, photo_url, self_rate, remark, status, points_rewarded, create_time) 
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	result, err := config.DB.Exec(query, record.UserID, record.HabitID, record.AssignmentID,
		record.CheckinDate, record.CheckinTime, record.CheckinType, record.PhotoURL, record.SelfRate, record.Remark,
		record.Status, record.PointsRewarded)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	record.ID = id
	return nil
}

// GetCheckinRecord 获取打卡记录
func GetCheckinRecord(userID, habitID int64, date string) (*CheckinRecord, error) {
	record := &CheckinRecord{}
	query := `SELECT id, user_id, habit_id, assignment_id, checkin_date, checkin_time, 
			 checkin_type, photo_url, remark, status, points_rewarded, create_time 
			 FROM checkin_record 
			 WHERE user_id = ? AND habit_id = ? AND checkin_date = ?`
	err := config.DB.QueryRow(query, userID, habitID, date).Scan(
		&record.ID, &record.UserID, &record.HabitID, &record.AssignmentID, &record.CheckinDate,
		&record.CheckinTime, &record.CheckinType, &record.PhotoURL, &record.Remark, &record.Status,
		&record.PointsRewarded, &record.CreateTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("checkin record not found")
		}
		return nil, err
	}
	return record, nil
}

// GetCheckinRecordsByUserID 获取用户的打卡记录
func GetCheckinRecordsByUserID(userID int64, startDate, endDate string) ([]*CheckinRecord, error) {
	query := `SELECT cr.id, cr.user_id, cr.habit_id, h.name, cr.assignment_id, cr.checkin_date, cr.checkin_time, 
			 cr.checkin_type, cr.photo_url, cr.self_rate, cr.remark, cr.status, cr.points_rewarded, cr.create_time 
			 FROM checkin_record cr
			 LEFT JOIN habit h ON cr.habit_id = h.id
			 WHERE cr.user_id = ? AND cr.checkin_date BETWEEN ? AND ? 
			 ORDER BY cr.checkin_time DESC`
	rows, err := config.DB.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []*CheckinRecord{}
	for rows.Next() {
		record := &CheckinRecord{}
		err := rows.Scan(
			&record.ID, &record.UserID, &record.HabitID, &record.HabitName, &record.AssignmentID, &record.CheckinDate,
			&record.CheckinTime, &record.CheckinType, &record.PhotoURL, &record.SelfRate, &record.Remark,
			&record.Status, &record.PointsRewarded, &record.CreateTime,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

// GetChildCheckinRecords 获取小孩的打卡记录（管理员）
func GetChildCheckinRecords(childID int64, startDate, endDate, habitID string) ([]*CheckinRecord, error) {
	query := `SELECT cr.id, cr.user_id, cr.habit_id, h.name, cr.assignment_id, cr.checkin_date, cr.checkin_time, 
			 cr.checkin_type, cr.photo_url, cr.remark, cr.status, cr.points_rewarded, cr.is_rolled_back,
			 cr.rollback_time, cr.rollback_reason, cr.create_time 
			 FROM checkin_record cr
			 LEFT JOIN habit h ON cr.habit_id = h.id
			 WHERE cr.user_id = ?`

	args := []interface{}{childID}

	if startDate != "" && endDate != "" {
		query += " AND cr.checkin_date BETWEEN ? AND ?"
		args = append(args, startDate, endDate)
	}

	if habitID != "" {
		query += " AND cr.habit_id = ?"
		args = append(args, habitID)
	}

	query += " ORDER BY cr.checkin_time DESC"

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []*CheckinRecord{}
	for rows.Next() {
		record := &CheckinRecord{}
		var rollbackTime sql.NullTime
		err := rows.Scan(
			&record.ID, &record.UserID, &record.HabitID, &record.HabitName, &record.AssignmentID,
			&record.CheckinDate, &record.CheckinTime, &record.CheckinType, &record.PhotoURL,
			&record.Remark, &record.Status, &record.PointsRewarded, &record.IsRolledBack,
			&rollbackTime, &record.RollbackReason, &record.CreateTime,
		)
		if err != nil {
			return nil, err
		}

		if rollbackTime.Valid {
			record.RollbackTime = &rollbackTime.Time
		}

		// 将 NullString 转换为普通字符串
		if record.RollbackReason.Valid {
			record.RollbackReasonStr = record.RollbackReason.String
		}

		records = append(records, record)
	}
	return records, nil
}

// GetCheckinRecords 获取打卡记录（支持更多过滤条件）
func GetCheckinRecords(userID int64, startDate, endDate, habitID string) ([]*CheckinRecord, error) {
	query := `SELECT cr.id, cr.user_id, cr.habit_id, h.name, cr.assignment_id, cr.checkin_date, cr.checkin_time, 
			 cr.checkin_type, cr.photo_url, cr.self_rate, cr.remark, cr.status, cr.points_rewarded, 
			 cr.is_rolled_back, cr.rollback_time, cr.rollback_reason, cr.create_time 
			 FROM checkin_record cr
			 LEFT JOIN habit h ON cr.habit_id = h.id
			 WHERE cr.user_id = ?`

	args := []interface{}{userID}

	if startDate != "" && endDate != "" {
		query += " AND cr.checkin_date BETWEEN ? AND ?"
		args = append(args, startDate, endDate)
	}

	if habitID != "" {
		query += " AND cr.habit_id = ?"
		args = append(args, habitID)
	}

	query += " ORDER BY cr.checkin_time DESC"

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []*CheckinRecord{}
	for rows.Next() {
		record := &CheckinRecord{}
		var rollbackTime sql.NullTime
		err := rows.Scan(
			&record.ID, &record.UserID, &record.HabitID, &record.HabitName, &record.AssignmentID, &record.CheckinDate,
			&record.CheckinTime, &record.CheckinType, &record.PhotoURL, &record.SelfRate, &record.Remark,
			&record.Status, &record.PointsRewarded, &record.IsRolledBack, &rollbackTime, &record.RollbackReason, &record.CreateTime,
		)
		if err != nil {
			return nil, err
		}

		if rollbackTime.Valid {
			record.RollbackTime = &rollbackTime.Time
		}

		// 将 NullString 转换为普通字符串
		if record.RollbackReason.Valid {
			record.RollbackReasonStr = record.RollbackReason.String
		}

		records = append(records, record)
	}
	return records, nil
}

// UpdateStreakRecord 更新连续打卡记录
func UpdateStreakRecord(userID, habitID int64, checkinDate string) error {
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

	// 获取当前连续打卡记录
	streak := &StreakRecord{}
	query := `SELECT id, current_streak, longest_streak, last_checkin_date, streak_start_date 
			 FROM streak_record 
			 WHERE user_id = ? AND habit_id = ?`
	err = tx.QueryRow(query, userID, habitID).Scan(
		&streak.ID, &streak.CurrentStreak, &streak.LongestStreak,
		&streak.LastCheckinDate, &streak.StreakStartDate,
	)

	// 解析日期
	lastDate, _ := time.Parse("2006-01-02", streak.LastCheckinDate)
	currentDate, _ := time.Parse("2006-01-02", checkinDate)
	dayDiff := int(currentDate.Sub(lastDate).Hours() / 24)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 第一次打卡
			insertQuery := `INSERT INTO streak_record (user_id, habit_id, current_streak, 
						 longest_streak, last_checkin_date, streak_start_date, update_time) 
						 VALUES (?, ?, 1, 1, ?, ?, NOW())`
			_, err = tx.Exec(insertQuery, userID, habitID, checkinDate, checkinDate)
		} else {
			return err
		}
	} else {
		// 更新连续打卡记录
		var updateQuery string
		if dayDiff == 1 {
			// 连续打卡
			streak.CurrentStreak++
			if streak.CurrentStreak > streak.LongestStreak {
				streak.LongestStreak = streak.CurrentStreak
			}
			updateQuery = `UPDATE streak_record SET current_streak = ?, longest_streak = ?, 
						 last_checkin_date = ?, update_time = NOW() 
						 WHERE id = ?`
			_, err = tx.Exec(updateQuery, streak.CurrentStreak, streak.LongestStreak, checkinDate, streak.ID)
		} else if dayDiff > 1 {
			// 中断，重置
			updateQuery = `UPDATE streak_record SET current_streak = 1, longest_streak = ?, 
						 last_checkin_date = ?, streak_start_date = ?, update_time = NOW() 
						 WHERE id = ?`
			_, err = tx.Exec(updateQuery, streak.LongestStreak, checkinDate, checkinDate, streak.ID)
		}
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetStreakRecord 获取连续打卡记录
func GetStreakRecord(userID, habitID int64) (*StreakRecord, error) {
	streak := &StreakRecord{}
	query := `SELECT id, user_id, habit_id, current_streak, longest_streak, 
			 last_checkin_date, streak_start_date, update_time 
			 FROM streak_record 
			 WHERE user_id = ? AND habit_id = ?`
	err := config.DB.QueryRow(query, userID, habitID).Scan(
		&streak.ID, &streak.UserID, &streak.HabitID, &streak.CurrentStreak, &streak.LongestStreak,
		&streak.LastCheckinDate, &streak.StreakStartDate, &streak.UpdateTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &StreakRecord{UserID: userID, HabitID: habitID, CurrentStreak: 0, LongestStreak: 0}, nil
		}
		return nil, err
	}
	return streak, nil
}

// GetCheckinRecordByID 根据 ID 获取打卡记录
func GetCheckinRecordByID(id int64) (*CheckinRecord, error) {
	record := &CheckinRecord{}
	query := `SELECT id, user_id, habit_id, assignment_id, checkin_date, checkin_time, 
			 checkin_type, photo_url, self_rate, remark, status, points_rewarded, 
			 is_rolled_back, rollback_time, rollback_reason, create_time 
			 FROM checkin_record WHERE id = ?`
	err := config.DB.QueryRow(query, id).Scan(
		&record.ID, &record.UserID, &record.HabitID, &record.AssignmentID, &record.CheckinDate,
		&record.CheckinTime, &record.CheckinType, &record.PhotoURL, &record.SelfRate, &record.Remark,
		&record.Status, &record.PointsRewarded, &record.IsRolledBack, &record.RollbackTime,
		&record.RollbackReason, &record.CreateTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("checkin record not found")
		}
		return nil, err
	}

	// 将 NullString 转换为普通字符串
	if record.RollbackReason.Valid {
		record.RollbackReasonStr = record.RollbackReason.String
	}

	return record, nil
}

// UpdateCheckinSelfRate 更新自我评分
func UpdateCheckinSelfRate(recordID int64, selfRate int) error {
	query := `UPDATE checkin_record SET self_rate = ? WHERE id = ?`
	_, err := config.DB.Exec(query, selfRate, recordID)
	return err
}

// GetDailyPointsStats 获取用户每天的积分统计
func GetDailyPointsStats(userID int64, startDate, endDate string) ([]map[string]interface{}, error) {
	query := `SELECT DATE(checkin_date) as checkin_date, SUM(points_rewarded) as daily_points 
			 FROM checkin_record 
			 WHERE user_id = ? AND checkin_date BETWEEN ? AND ? AND is_rolled_back = 0 
			 GROUP BY DATE(checkin_date) 
			 ORDER BY checkin_date ASC`

	rows, err := config.DB.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := []map[string]interface{}{}
	for rows.Next() {
		var date string
		var points int
		if err := rows.Scan(&date, &points); err != nil {
			return nil, err
		}

		stats = append(stats, map[string]interface{}{
			"date":   date,
			"points": points,
		})
	}

	return stats, nil
}
