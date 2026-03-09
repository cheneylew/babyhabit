package models

import (
	"babyhabit/config"
	"database/sql"
	"errors"
	"time"
)

type Habit struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Icon              string    `json:"icon"`
	Category          string    `json:"category"`
	ScheduleType      int       `json:"schedule_type"`
	ScheduleDetail    string    `json:"schedule_detail"`
	CheckinTimeStart  string    `json:"checkin_time_start"`
	CheckinTimeEnd    string    `json:"checkin_time_end"`
	RewardPoints      int       `json:"reward_points"`
	AllowMakeup       int       `json:"allow_makeup"`
	MakeupDays        int       `json:"makeup_days"`
	CreatorID         int64     `json:"creator_id"`
	CreateTime        time.Time `json:"create_time"`
	UpdateTime        time.Time `json:"update_time"`
	Status            int       `json:"status"`
}

type HabitAssignment struct {
	ID         int       `json:"id"`
	HabitID    int       `json:"habit_id"`
	ChildID    int       `json:"child_id"`
	AssignTime time.Time `json:"assign_time"`
	Status     int       `json:"status"`
}

// CreateHabit 创建习惯
func CreateHabit(habit *Habit) error {
	query := `INSERT INTO habit (name, description, icon, category, schedule_type, schedule_detail, 
			 checkin_time_start, checkin_time_end, reward_points, allow_makeup, makeup_days, 
			 creator_id, create_time, update_time, status) 
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), ?)`
	result, err := config.DB.Exec(query, habit.Name, habit.Description, habit.Icon, habit.Category, 
		habit.ScheduleType, habit.ScheduleDetail, habit.CheckinTimeStart, habit.CheckinTimeEnd, 
		habit.RewardPoints, habit.AllowMakeup, habit.MakeupDays, habit.CreatorID, habit.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	habit.ID = id
	return nil
}

// GetHabitByID 根据ID获取习惯
func GetHabitByID(id int64) (*Habit, error) {
	habit := &Habit{}
	query := `SELECT id, name, description, icon, category, schedule_type, schedule_detail, 
			 checkin_time_start, checkin_time_end, reward_points, allow_makeup, makeup_days, 
			 creator_id, create_time, update_time, status 
			 FROM habit WHERE id = ?`
	err := config.DB.QueryRow(query, id).Scan(
		&habit.ID, &habit.Name, &habit.Description, &habit.Icon, &habit.Category, &habit.ScheduleType, 
		&habit.ScheduleDetail, &habit.CheckinTimeStart, &habit.CheckinTimeEnd, &habit.RewardPoints, 
		&habit.AllowMakeup, &habit.MakeupDays, &habit.CreatorID, &habit.CreateTime, &habit.UpdateTime, &habit.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("habit not found")
		}
		return nil, err
	}
	return habit, nil
}

// GetHabitsByCreatorID 根据创建者ID获取习惯列表
func GetHabitsByCreatorID(creatorID int64) ([]*Habit, error) {
	query := `SELECT id, name, description, icon, category, schedule_type, schedule_detail, 
			 checkin_time_start, checkin_time_end, reward_points, allow_makeup, makeup_days, 
			 creator_id, create_time, update_time, status 
			 FROM habit WHERE creator_id = ? ORDER BY create_time DESC`
	rows, err := config.DB.Query(query, creatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits := []*Habit{}
	for rows.Next() {
		habit := &Habit{}
		err := rows.Scan(
			&habit.ID, &habit.Name, &habit.Description, &habit.Icon, &habit.Category, &habit.ScheduleType, 
			&habit.ScheduleDetail, &habit.CheckinTimeStart, &habit.CheckinTimeEnd, &habit.RewardPoints, 
			&habit.AllowMakeup, &habit.MakeupDays, &habit.CreatorID, &habit.CreateTime, &habit.UpdateTime, &habit.Status,
		)
		if err != nil {
			return nil, err
		}
		habits = append(habits, habit)
	}
	return habits, nil
}

// UpdateHabit 更新习惯
func UpdateHabit(habit *Habit) error {
	query := `UPDATE habit SET name = ?, description = ?, icon = ?, category = ?, schedule_type = ?, 
			 schedule_detail = ?, checkin_time_start = ?, checkin_time_end = ?, reward_points = ?, 
			 allow_makeup = ?, makeup_days = ?, status = ?, update_time = NOW() 
			 WHERE id = ?`
	_, err := config.DB.Exec(query, habit.Name, habit.Description, habit.Icon, habit.Category, 
		habit.ScheduleType, habit.ScheduleDetail, habit.CheckinTimeStart, habit.CheckinTimeEnd, 
		habit.RewardPoints, habit.AllowMakeup, habit.MakeupDays, habit.Status, habit.ID)
	return err
}

// DeleteHabit 删除习惯
func DeleteHabit(id int64) error {
	query := `DELETE FROM habit WHERE id = ?`
	_, err := config.DB.Exec(query, id)
	return err
}

// AssignHabit 分配习惯给小孩
func AssignHabit(assignment *HabitAssignment) error {
	query := `INSERT INTO habit_assignment (habit_id, child_id, assign_time, status) 
			 VALUES (?, ?, NOW(), ?)`
	result, err := config.DB.Exec(query, assignment.HabitID, assignment.ChildID, assignment.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	assignment.ID = int(id)
	return nil
}

// GetHabitsByChildID 获取小孩的习惯列表
func GetHabitsByChildID(childID int64) ([]*Habit, error) {
	query := `SELECT h.id, h.name, h.description, h.icon, h.category, h.schedule_type, 
			 h.schedule_detail, h.checkin_time_start, h.checkin_time_end, h.reward_points, 
			 h.allow_makeup, h.makeup_days, h.creator_id, h.create_time, h.update_time, h.status 
			 FROM habit h 
			 JOIN habit_assignment ha ON h.id = ha.habit_id 
			 WHERE ha.child_id = ? AND ha.status = 1 AND h.status = 1`
	rows, err := config.DB.Query(query, childID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits := []*Habit{}
	for rows.Next() {
		habit := &Habit{}
		err := rows.Scan(
			&habit.ID, &habit.Name, &habit.Description, &habit.Icon, &habit.Category, &habit.ScheduleType, 
			&habit.ScheduleDetail, &habit.CheckinTimeStart, &habit.CheckinTimeEnd, &habit.RewardPoints, 
			&habit.AllowMakeup, &habit.MakeupDays, &habit.CreatorID, &habit.CreateTime, &habit.UpdateTime, &habit.Status,
		)
		if err != nil {
			return nil, err
		}
		habits = append(habits, habit)
	}
	return habits, nil
}