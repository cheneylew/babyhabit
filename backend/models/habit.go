package models

import (
	"babyhabit/config"
	"database/sql"
	"errors"
	"time"
)

type Habit struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Icon             string    `json:"icon"`
	Category         string    `json:"category"`
	ScheduleType     int       `json:"schedule_type"`
	ScheduleDetail   string    `json:"schedule_detail"`
	CheckinTimeStart string    `json:"checkin_time_start"`
	CheckinTimeEnd   string    `json:"checkin_time_end"`
	RewardPoints     int       `json:"reward_points"`
	AllowMakeup      int       `json:"allow_makeup"`
	MakeupDays       int       `json:"makeup_days"`
	RequirePhoto     int       `json:"require_photo"`
	AllowSelfRate    int       `json:"allow_self_rate"`
	CheckinPrompt    string    `json:"checkin_prompt"`
	CreatorID        int64     `json:"creator_id"`
	CreateTime       time.Time `json:"create_time"`
	UpdateTime       time.Time `json:"update_time"`
	Status           int       `json:"status"`

	// 数据库扫描用（不导出到 JSON）
	CheckinPromptRaw sql.NullString `json:"-"`
}

type HabitAssignment struct {
	ID         int       `json:"id"`
	HabitID    int       `json:"habit_id"`
	ChildID    int       `json:"child_id"`
	AssignTime time.Time `json:"assign_time"`
	Status     int       `json:"status"`
	Habit      *Habit    `json:"habit,omitempty"`
}

// CreateHabit 创建习惯
func CreateHabit(habit *Habit) error {
	query := `INSERT INTO habit (name, description, icon, category, schedule_type, schedule_detail, 
			 checkin_time_start, checkin_time_end, reward_points, allow_makeup, makeup_days, require_photo,
			 allow_self_rate, checkin_prompt, creator_id, create_time, update_time, status) 
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), ?)`
	result, err := config.DB.Exec(query, habit.Name, habit.Description, habit.Icon, habit.Category,
		habit.ScheduleType, habit.ScheduleDetail, habit.CheckinTimeStart, habit.CheckinTimeEnd,
		habit.RewardPoints, habit.AllowMakeup, habit.MakeupDays, habit.RequirePhoto, habit.AllowSelfRate, habit.CheckinPrompt, habit.CreatorID, habit.Status)
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

// GetHabitByID 根据 ID 获取习惯
func GetHabitByID(id int64) (*Habit, error) {
	habit := &Habit{}
	query := `SELECT id, name, description, icon, category, schedule_type, schedule_detail, 
			 checkin_time_start, checkin_time_end, reward_points, allow_makeup, makeup_days, require_photo,
			 allow_self_rate, checkin_prompt, creator_id, create_time, update_time, status 
			 FROM habit WHERE id = ?`
	err := config.DB.QueryRow(query, id).Scan(
		&habit.ID, &habit.Name, &habit.Description, &habit.Icon, &habit.Category, &habit.ScheduleType,
		&habit.ScheduleDetail, &habit.CheckinTimeStart, &habit.CheckinTimeEnd, &habit.RewardPoints,
		&habit.AllowMakeup, &habit.MakeupDays, &habit.RequirePhoto, &habit.AllowSelfRate, &habit.CheckinPromptRaw, &habit.CreatorID, &habit.CreateTime, &habit.UpdateTime, &habit.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("habit not found")
		}
		return nil, err
	}

	// 将 NullString 转换为普通字符串
	if habit.CheckinPromptRaw.Valid {
		habit.CheckinPrompt = habit.CheckinPromptRaw.String
	}

	return habit, nil
}

// GetHabitsByCreatorID 根据创建者 ID 获取习惯列表
func GetHabitsByCreatorID(creatorID int64) ([]*Habit, error) {
	query := `SELECT id, name, description, icon, category, schedule_type, schedule_detail, 
			 checkin_time_start, checkin_time_end, reward_points, allow_makeup, makeup_days, require_photo,
			 allow_self_rate, checkin_prompt, creator_id, create_time, update_time, status 
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
			&habit.AllowMakeup, &habit.MakeupDays, &habit.RequirePhoto, &habit.AllowSelfRate, &habit.CheckinPromptRaw, &habit.CreatorID, &habit.CreateTime, &habit.UpdateTime, &habit.Status,
		)
		if err != nil {
			println("=== Scan Error ===")
			println("Error:", err.Error())
			println("==================")
			return nil, err
		}

		// 将 NullString 转换为普通字符串
		if habit.CheckinPromptRaw.Valid {
			habit.CheckinPrompt = habit.CheckinPromptRaw.String
		}

		habits = append(habits, habit)
	}
	return habits, nil
}

// GetAssignedHabitsByChildID 获取小孩已分配的习惯列表（包含分配信息）
func GetAssignedHabitsByChildID(childID int64) ([]*HabitAssignment, error) {
	query := `SELECT ha.id, ha.habit_id, ha.child_id, ha.assign_time, ha.status,
			 h.name, h.description, h.icon, h.category, h.schedule_type, h.schedule_detail,
			 h.checkin_time_start, h.checkin_time_end, h.reward_points, h.allow_makeup, h.makeup_days, h.require_photo,
			 h.allow_self_rate, h.checkin_prompt
			 FROM habit_assignment ha
			 JOIN habit h ON ha.habit_id = h.id
			 WHERE ha.child_id = ? AND ha.status = 1
			 ORDER BY ha.assign_time DESC`
	rows, err := config.DB.Query(query, childID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := []*HabitAssignment{}
	for rows.Next() {
		assignment := &HabitAssignment{}
		habit := &Habit{}
		err := rows.Scan(
			&assignment.ID, &assignment.HabitID, &assignment.ChildID, &assignment.AssignTime, &assignment.Status,
			&habit.Name, &habit.Description, &habit.Icon, &habit.Category, &habit.ScheduleType, &habit.ScheduleDetail,
			&habit.CheckinTimeStart, &habit.CheckinTimeEnd, &habit.RewardPoints, &habit.AllowMakeup, &habit.MakeupDays, &habit.RequirePhoto,
			&habit.AllowSelfRate, &habit.CheckinPromptRaw,
		)
		if err != nil {
			return nil, err
		}

		// 将 NullString 转换为普通字符串
		if habit.CheckinPromptRaw.Valid {
			habit.CheckinPrompt = habit.CheckinPromptRaw.String
		}
		assignment.Habit = habit
		assignments = append(assignments, assignment)
	}
	return assignments, nil
}

// DeleteHabitAssignment 删除习惯分配记录
func DeleteHabitAssignment(assignmentID int) error {
	query := `UPDATE habit_assignment SET status = 0 WHERE id = ?`
	_, err := config.DB.Exec(query, assignmentID)
	return err
}

// UpdateHabit 更新习惯
func UpdateHabit(habit *Habit) error {
	query := `UPDATE habit SET name = ?, description = ?, icon = ?, category = ?, schedule_type = ?, 
			 schedule_detail = ?, checkin_time_start = ?, checkin_time_end = ?, reward_points = ?, 
			 allow_makeup = ?, makeup_days = ?, require_photo = ?, allow_self_rate = ?, checkin_prompt = ?, status = ?, update_time = NOW() 
			 WHERE id = ?`
	_, err := config.DB.Exec(query, habit.Name, habit.Description, habit.Icon, habit.Category,
		habit.ScheduleType, habit.ScheduleDetail, habit.CheckinTimeStart, habit.CheckinTimeEnd,
		habit.RewardPoints, habit.AllowMakeup, habit.MakeupDays, habit.RequirePhoto, habit.AllowSelfRate, habit.CheckinPrompt, habit.Status, habit.ID)
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

// DeleteHabitAssignmentsByChildID 按照小孩ID删除所有习惯分配
func DeleteHabitAssignmentsByChildID(childID int) error {
	query := `DELETE FROM habit_assignment WHERE child_id = ?`
	_, err := config.DB.Exec(query, childID)
	return err
}

// GetHabitsByChildID 获取小孩的习惯列表
func GetHabitsByChildID(childID int64) ([]*Habit, error) {
	query := `SELECT DISTINCT h.id, h.name, h.description, h.icon, h.category, h.schedule_type, 
			 h.schedule_detail, h.checkin_time_start, h.checkin_time_end, h.reward_points, 
			 h.allow_makeup, h.makeup_days, h.require_photo, h.allow_self_rate, h.checkin_prompt, h.creator_id, h.create_time, h.update_time, h.status 
			 FROM habit h 
			 JOIN habit_assignment ha ON h.id = ha.habit_id
			 WHERE ha.child_id = ? AND h.status = 1 AND ha.status = 1`
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
			&habit.AllowMakeup, &habit.MakeupDays, &habit.RequirePhoto, &habit.AllowSelfRate, &habit.CheckinPromptRaw, &habit.CreatorID, &habit.CreateTime, &habit.UpdateTime, &habit.Status,
		)
		if err != nil {
			return nil, err
		}

		// 将 NullString 转换为普通字符串
		if habit.CheckinPromptRaw.Valid {
			habit.CheckinPrompt = habit.CheckinPromptRaw.String
		}

		habits = append(habits, habit)
	}
	return habits, nil
}
