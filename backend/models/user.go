package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/cheneylew/babyhabit/backend/config"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               int64      `json:"id"`
	Username         string     `json:"username"`
	Password         string     `json:"-"`
	Name             string     `json:"name"`
	Phone            string     `json:"phone"`
	Email            string     `json:"email"`
	Avatar           string     `json:"avatar"`
	RegisterTime     *time.Time `json:"register_time"`
	LastLoginTime    *time.Time `json:"last_login_time"`
	LastActivityTime *time.Time `json:"last_activity_time"`
	UserType         int        `json:"user_type"`
	Status           int        `json:"status"`
	ParentID         int64      `json:"parent_id"`
	LoginFailCount   int        `json:"login_fail_count"`
	LockedUntil      *time.Time `json:"locked_until"`
	PointsBalance    int        `json:"points_balance"`
	CreateTime       *time.Time `json:"create_time"`
	UpdateTime       *time.Time `json:"update_time"`
}

// CreateUser 创建用户
func CreateUser(user *User) error {
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	query := `INSERT INTO user (username, password, name, phone, email, avatar, user_type, parent_id, register_time, create_time) 
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
	result, err := config.DB.Exec(query, user.Username, user.Password, user.Name, user.Phone, user.Email, user.Avatar, user.UserType, user.ParentID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(id int64) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password, name, phone, email, avatar, register_time, last_login_time, last_activity_time, 
			 user_type, status, parent_id, login_fail_count, locked_until, points_balance, create_time, update_time 
			 FROM user WHERE id = ?`
	err := config.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Name, &user.Phone, &user.Email, &user.Avatar,
		&user.RegisterTime, &user.LastLoginTime, &user.LastActivityTime, &user.UserType, &user.Status,
		&user.ParentID, &user.LoginFailCount, &user.LockedUntil, &user.PointsBalance, &user.CreateTime, &user.UpdateTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password, name, phone, email, avatar, register_time, last_login_time, last_activity_time, 
			 user_type, status, parent_id, login_fail_count, locked_until, points_balance, create_time, update_time 
			 FROM user WHERE username = ?`
	err := config.DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Name, &user.Phone, &user.Email, &user.Avatar,
		&user.RegisterTime, &user.LastLoginTime, &user.LastActivityTime, &user.UserType, &user.Status,
		&user.ParentID, &user.LoginFailCount, &user.LockedUntil, &user.PointsBalance, &user.CreateTime, &user.UpdateTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(user *User) error {
	query := `UPDATE user SET name = ?, phone = ?, email = ?, avatar = ?, status = ?, update_time = NOW() WHERE id = ?`
	_, err := config.DB.Exec(query, user.Name, user.Phone, user.Email, user.Avatar, user.Status, user.ID)
	return err
}

// DeleteUser 删除用户
func DeleteUser(id int64) error {
	query := `DELETE FROM user WHERE id = ?`
	_, err := config.DB.Exec(query, id)
	return err
}

// UpdateUserLastLogin 更新用户最后登录时间
func UpdateUserLastLogin(id int64) error {
	query := `UPDATE user SET last_login_time = NOW(), last_activity_time = NOW(), login_fail_count = 0 WHERE id = ?`
	_, err := config.DB.Exec(query, id)
	return err
}

// UpdateLoginFailCount 更新登录失败次数
func UpdateLoginFailCount(id int64) error {
	query := `UPDATE user SET login_fail_count = login_fail_count + 1, locked_until = CASE 
			 WHEN login_fail_count >= 4 THEN NOW() + INTERVAL 15 MINUTE 
			 ELSE locked_until 
			 END 
			 WHERE id = ?`
	_, err := config.DB.Exec(query, id)
	return err
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsLocked 检查用户是否被锁定
func (u *User) IsLocked() bool {
	return u.LockedUntil != nil && time.Now().Before(*u.LockedUntil)
}

// GetChildrenByParentID 根据父母ID获取小孩列表
func GetChildrenByParentID(parentID int64) ([]*User, error) {
	query := `SELECT id, username, name, phone, email, avatar, user_type, status, points_balance, create_time 
			 FROM user WHERE parent_id = ? AND user_type = 2`
	rows, err := config.DB.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	children := []*User{}
	for rows.Next() {
		child := &User{}
		err := rows.Scan(
			&child.ID, &child.Username, &child.Name, &child.Phone, &child.Email, &child.Avatar,
			&child.UserType, &child.Status, &child.PointsBalance, &child.CreateTime,
		)
		if err != nil {
			return nil, err
		}
		children = append(children, child)
	}
	return children, nil
}

// GetAllChildren 获取所有小孩列表
func GetAllChildren() ([]*User, error) {
	query := `SELECT id, username, name, phone, email, avatar, user_type, status, points_balance, create_time 
			 FROM user WHERE user_type = 2`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	children := []*User{}
	for rows.Next() {
		child := &User{}
		err := rows.Scan(
			&child.ID, &child.Username, &child.Name, &child.Phone, &child.Email, &child.Avatar,
			&child.UserType, &child.Status, &child.PointsBalance, &child.CreateTime,
		)
		if err != nil {
			return nil, err
		}
		children = append(children, child)
	}
	return children, nil
}

// GetChildIDsByParentID 根据父母ID获取孩子ID列表
func GetChildIDsByParentID(parentID int64) ([]int64, error) {
	query := `SELECT id FROM user WHERE parent_id = ? AND user_type = 2 AND status = 1`
	rows, err := config.DB.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var childIDs []int64
	for rows.Next() {
		var childID int64
		err := rows.Scan(&childID)
		if err != nil {
			return nil, err
		}
		childIDs = append(childIDs, childID)
	}

	return childIDs, nil
}
