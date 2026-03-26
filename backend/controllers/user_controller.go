package controllers

import (
	"net/http"
	"strconv"

	"github.com/cheneylew/babyhabit/backend/middleware"
	"github.com/cheneylew/babyhabit/backend/models"

	"github.com/gin-gonic/gin"
)

// RegisterRequest 注册请求
func RegisterRequest(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		UserType int    `json:"user_type" binding:"required,oneof=1 2"`
		ParentID int64  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户类型和父ID
	if req.UserType == 2 && req.ParentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parent ID is required for child account"})
		return
	}

	// 检查用户名是否已存在
	existingUser, _ := models.GetUserByUsername(req.Username)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		UserType: req.UserType,
		ParentID: req.ParentID,
		Status:   1,
	}

	if err := models.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 生成token
	token, err := middleware.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"name":      user.Name,
			"phone":     user.Phone,
			"email":     user.Email,
			"user_type": user.UserType,
			"status":    user.Status,
		},
	})
}

// LoginRequest 登录请求
func LoginRequest(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户
	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		// 用户不存在，直接返回错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 检查用户是否被锁定
	if user.IsLocked() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is locked, please try again later"})
		return
	}

	// 验证密码
	if !user.VerifyPassword(req.Password) {
		// 记录登录失败
		models.UpdateLoginFailCount(user.ID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 更新最后登录时间
	models.UpdateUserLastLogin(user.ID)

	// 生成token
	token, err := middleware.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":             user.ID,
			"username":       user.Username,
			"name":           user.Name,
			"phone":          user.Phone,
			"email":          user.Email,
			"user_type":      user.UserType,
			"status":         user.Status,
			"points_balance": user.PointsBalance,
		},
	})
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":              user.ID,
			"username":        user.Username,
			"name":            user.Name,
			"phone":           user.Phone,
			"email":           user.Email,
			"avatar":          user.Avatar,
			"user_type":       user.UserType,
			"status":          user.Status,
			"points_balance":  user.PointsBalance,
			"register_time":   user.RegisterTime,
			"last_login_time": user.LastLoginTime,
		},
	})
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var req struct {
		Name   string `json:"name"`
		Phone  string `json:"phone"`
		Email  string `json:"email"`
		Avatar string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新用户信息
	user.Name = req.Name
	user.Phone = req.Phone
	user.Email = req.Email
	user.Avatar = req.Avatar

	if err := models.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"name":     user.Name,
			"phone":    user.Phone,
			"email":    user.Email,
			"avatar":   user.Avatar,
		},
	})
}

// GetChildren 获取小孩列表
func GetChildren(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// 只有父母账号或管理员账号可以获取小孩列表
	if user.UserType != 1 && user.UserType != 3 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var children []*models.User
	var err error

	// 如果是父母账号，只获取自己的小孩
	if user.UserType == 1 {
		children, err = models.GetChildrenByParentID(user.ID)
	} else {
		// 如果是管理员账号，获取所有小孩
		children, err = models.GetAllChildren()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"children": children})
}

// CreateChild 创建小孩账号
func CreateChild(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// 只有父母账号可以创建小孩账号
	if user.UserType != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	existingUser, _ := models.GetUserByUsername(req.Username)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// 创建小孩账号
	child := &models.User{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		UserType: 2,
		ParentID: user.ID,
		Status:   1,
	}

	if err := models.CreateUser(child); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"child": gin.H{
			"id":        child.ID,
			"username":  child.Username,
			"name":      child.Name,
			"phone":     child.Phone,
			"email":     child.Email,
			"user_type": child.UserType,
			"status":    child.Status,
		},
	})
}

// UpdateChild 更新小孩账号
func UpdateChild(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// 只有父母账号或管理员账号可以更新小孩账号
	if user.UserType != 1 && user.UserType != 3 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	childID := c.Param("id")
	childIDInt, err := strconv.ParseInt(childID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid child ID"})
		return
	}

	// 获取小孩信息
	child, err := models.GetUserByID(childIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Child not found"})
		return
	}

	// 检查权限：父母只能更新自己的小孩，管理员可以更新所有小孩
	if user.UserType == 1 && child.ParentID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var req struct {
		Name   string `json:"name"`
		Phone  string `json:"phone"`
		Email  string `json:"email"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新小孩账号信息
	child.Name = req.Name
	child.Phone = req.Phone
	child.Email = req.Email
	child.Status = req.Status

	if err := models.UpdateUser(child); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"child": gin.H{
			"id":        child.ID,
			"username":  child.Username,
			"name":      child.Name,
			"phone":     child.Phone,
			"email":     child.Email,
			"user_type": child.UserType,
			"status":    child.Status,
		},
	})
}

// DeleteChild 删除小孩账号
func DeleteChild(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// 只有父母账号或管理员账号可以删除小孩账号
	if user.UserType != 1 && user.UserType != 3 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	childID := c.Param("id")
	childIDInt, err := strconv.ParseInt(childID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid child ID"})
		return
	}

	// 获取小孩信息
	child, err := models.GetUserByID(childIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Child not found"})
		return
	}

	// 检查权限：父母只能删除自己的小孩，管理员可以删除所有小孩
	if user.UserType == 1 && child.ParentID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// 删除小孩账号
	if err := models.DeleteUser(childIDInt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Child deleted successfully"})
}

// GetChild 获取单个小孩信息
func GetChild(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// 只有父母账号或管理员账号可以获取小孩信息
	if user.UserType != 1 && user.UserType != 3 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	childID := c.Param("id")
	childIDInt, err := strconv.ParseInt(childID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid child ID"})
		return
	}

	// 获取小孩信息
	child, err := models.GetUserByID(childIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Child not found"})
		return
	}

	// 检查权限：父母只能查看自己的小孩，管理员可以查看所有小孩
	if user.UserType == 1 && child.ParentID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"child": gin.H{
			"id":        child.ID,
			"username":  child.Username,
			"name":      child.Name,
			"phone":     child.Phone,
			"email":     child.Email,
			"user_type": child.UserType,
			"status":    child.Status,
		},
	})
}
