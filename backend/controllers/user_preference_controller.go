package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"babyhabit/models"
)

// GetUserPreference 获取用户特定偏好设置
func GetUserPreference(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	userID := user.ID

	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "偏好设置键不能为空"})
		return
	}

	// 如果提供了user_id参数且当前用户是管理员，则使用提供的user_id
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if adminUserID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil && adminUserID > 0 && user.UserType == 1 {
			userID = adminUserID
		}
	}

	preference, err := models.GetUserPreference(userID, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取偏好设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"preference": preference,
	})
}

// GetUserPreferences 获取用户所有偏好设置
func GetUserPreferences(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	userID := user.ID

	preferences, err := models.GetUserPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取偏好设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"preferences": preferences,
	})
}

// SetUserPreference 设置用户偏好设置
func SetUserPreference(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	userID := user.ID

	type Request struct {
		Key    string `json:"key" binding:"required"`
		Value  string `json:"value" binding:"required"`
		UserID int64  `json:"user_id"`
	}

	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 如果提供了UserID且当前用户是管理员，则使用提供的UserID
	if request.UserID > 0 && user.UserType == 1 {
		userID = request.UserID
	}

	err := models.SetUserPreference(userID, request.Key, request.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "设置偏好设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "偏好设置更新成功",
	})
}

// DeleteUserPreference 删除用户偏好设置
func DeleteUserPreference(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	userID := user.ID

	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "偏好设置键不能为空"})
		return
	}

	err := models.DeleteUserPreference(userID, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除偏好设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "偏好设置删除成功",
	})
}
