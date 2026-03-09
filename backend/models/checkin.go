package models

import (
	"babyhabit/config"
	"database/sql"
	"errors"
	"time"
)

type CheckinRecord struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	HabitID       int64     `json:"habit_id"`
	AssignmentID  int64     `json:"assignment_id"`
	CheckinDate   string    `json:"checkin_date"`
	CheckinTime   time.Time `json:"checkin_time"`
	CheckinType   int       `json:"checkin_type"`
	Remark        string    `json:"remark"`
	Status        int       `json:"status"`
	PointsRewarded int       `json:"points_rewarded"`
	CreateTime    time.Time `json:"create_time"`
}

type StreakRecord struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	HabitID          int64     `json:"habit_id"`
	CurrentStreak    int       `json:"current_streak"`
	LongestStreak    int       `json:"longest_streak"`
	LastCheckinDate  string    `json:"last_checkin_date"`
	StreakStartDate  string    `json:"streak_start_date"`
	UpdateTime       time.Time `json:"update_time"`
}

// CreateCheckinRecord 创建打卡记录
func CreateCheckinRecord(record *CheckinRecord) error {
	query := `INSERT INTO checkin_record (user_id, habit_id, assignment_id, checkin_date, 
			 checkin_time, checkin_type, remark, status, points_rewarded, create_time) 
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	result, err := config.DB.Exec(query, record.UserID, record.HabitID, record.AssignmentID, 
		record.CheckinDate, record.CheckinTime, record.CheckinType, record.Remark, 
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
			 checkin_type, remark, status, points_rewarded, create_time 
			 FROM checkin_record 
			 WHERE user_id = ? AND habit_id = ? AND checkin_date = ?`
	err := config.DB.QueryRow(query, userID, habitID, date).Scan(
		&record.ID, &record.UserID, &record.HabitID, &record.AssignmentID, &record.CheckinDate, 
		&record.CheckinTime, &record.CheckinType, &record.Remark, &record.Status, 
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
	query := `SELECT id, user_id, habit_id, assignment_id, checkin_date, checkin_time, 
			 checkin_type, remark, status, points_rewarded, create_time 
			 FROM checkin_record 
			 WHERE user_id = ? AND checkin_date BETWEEN ? AND ? 
			 ORDER BY checkin_date DESC`
	rows, err := config.DB.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []*CheckinRecord{}
	for rows.Next() {
		record := &CheckinRecord{}
		err := rows.Scan(
			&record.ID, &record.UserID, &record.HabitID, &record.AssignmentID, &record.CheckinDate, 
			&record.CheckinTime, &record.CheckinType, &record.Remark, &record.Status, 
			&record.PointsRewarded, &record.CreateTime,
		)
		if err != nil {
			return nil, err
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