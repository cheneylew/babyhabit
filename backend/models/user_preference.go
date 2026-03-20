package models

import (
	"database/sql"
	"time"

	"babyhabit/config"
)

// UserPreference 用户偏好设置模型
type UserPreference struct {
	ID             int       `json:"id"`
	UserID         int64     `json:"user_id"`
	PreferenceKey  string    `json:"preference_key"`
	PreferenceValue string    `json:"preference_value"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// GetUserPreference 获取用户特定偏好设置
func GetUserPreference(userID int64, key string) (*UserPreference, error) {
	query := `SELECT id, user_id, preference_key, preference_value, created_at, updated_at FROM user_preference WHERE user_id = ? AND preference_key = ?`
	var preference UserPreference
	err := config.DB.QueryRow(query, userID, key).Scan(
		&preference.ID, &preference.UserID, &preference.PreferenceKey, 
		&preference.PreferenceValue, &preference.CreatedAt, &preference.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &preference, nil
}

// GetUserPreferences 获取用户所有偏好设置
func GetUserPreferences(userID int64) ([]*UserPreference, error) {
	query := `SELECT id, user_id, preference_key, preference_value, created_at, updated_at FROM user_preference WHERE user_id = ?`
	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var preferences []*UserPreference
	for rows.Next() {
		var preference UserPreference
		err := rows.Scan(
			&preference.ID, &preference.UserID, &preference.PreferenceKey, 
			&preference.PreferenceValue, &preference.CreatedAt, &preference.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		preferences = append(preferences, &preference)
	}

	return preferences, nil
}

// SetUserPreference 设置用户偏好设置（如果不存在则创建，存在则更新）
func SetUserPreference(userID int64, key string, value string) error {
	// 检查是否存在
	preference, err := GetUserPreference(userID, key)
	if err != nil {
		return err
	}

	if preference != nil {
		// 更新
		query := `UPDATE user_preference SET preference_value = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
		_, err := config.DB.Exec(query, value, preference.ID)
		return err
	} else {
		// 创建
		query := `INSERT INTO user_preference (user_id, preference_key, preference_value) VALUES (?, ?, ?)`
		_, err := config.DB.Exec(query, userID, key, value)
		return err
	}
}

// DeleteUserPreference 删除用户偏好设置
func DeleteUserPreference(userID int64, key string) error {
	query := `DELETE FROM user_preference WHERE user_id = ? AND preference_key = ?`
	_, err := config.DB.Exec(query, userID, key)
	return err
}
