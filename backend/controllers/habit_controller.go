package controllers

import (
	"babyhabit/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateHabit 创建习惯
func CreateHabit(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var req struct {
		Name              string `json:"name" binding:"required"`
		Description       string `json:"description"`
		Icon              string `json:"icon"`
		Category          string `json:"category"`
		ScheduleType      int    `json:"schedule_type" binding:"required,oneof=1 2"`
		ScheduleDetail    string `json:"schedule_detail"`
		CheckinTimeStart  string `json:"checkin_time_start"`
		CheckinTimeEnd    string `json:"checkin_time_end"`
		RewardPoints      int    `json:"reward_points"`
		AllowMakeup       int    `json:"allow_makeup"`
		MakeupDays        int    `json:"makeup_days"`
		Status            int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	habit := &models.Habit{
		Name:              req.Name,
		Description:       req.Description,
		Icon:              req.Icon,
		Category:          req.Category,
		ScheduleType:      req.ScheduleType,
		ScheduleDetail:    req.ScheduleDetail,
		CheckinTimeStart:  req.CheckinTimeStart,
		CheckinTimeEnd:    req.CheckinTimeEnd,
		RewardPoints:      req.RewardPoints,
		AllowMakeup:       req.AllowMakeup,
		MakeupDays:        req.MakeupDays,
		CreatorID:         user.ID,
		Status:            req.Status,
	}

	if err := models.CreateHabit(habit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"habit": habit})
}

// GetHabits 获取习惯列表
func GetHabits(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	habits, err := models.GetHabitsByCreatorID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"habits": habits})
}

// GetHabit 获取习惯详情
func GetHabit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid habit ID"})
		return
	}

	habit, err := models.GetHabitByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"habit": habit})
}

// UpdateHabit 更新习惯
func UpdateHabit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid habit ID"})
		return
	}

	habit, err := models.GetHabitByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Name              string `json:"name"`
		Description       string `json:"description"`
		Icon              string `json:"icon"`
		Category          string `json:"category"`
		ScheduleType      int    `json:"schedule_type"`
		ScheduleDetail    string `json:"schedule_detail"`
		CheckinTimeStart  string `json:"checkin_time_start"`
		CheckinTimeEnd    string `json:"checkin_time_end"`
		RewardPoints      int    `json:"reward_points"`
		AllowMakeup       int    `json:"allow_makeup"`
		MakeupDays        int    `json:"makeup_days"`
		Status            int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新习惯信息
	habit.Name = req.Name
	habit.Description = req.Description
	habit.Icon = req.Icon
	habit.Category = req.Category
	habit.ScheduleType = req.ScheduleType
	habit.ScheduleDetail = req.ScheduleDetail
	habit.CheckinTimeStart = req.CheckinTimeStart
	habit.CheckinTimeEnd = req.CheckinTimeEnd
	habit.RewardPoints = req.RewardPoints
	habit.AllowMakeup = req.AllowMakeup
	habit.MakeupDays = req.MakeupDays
	habit.Status = req.Status

	if err := models.UpdateHabit(habit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"habit": habit})
}

// DeleteHabit 删除习惯
func DeleteHabit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid habit ID"})
		return
	}

	if err := models.DeleteHabit(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Habit deleted successfully"})
}

// AssignHabit 分配习惯给小孩
func AssignHabit(c *gin.Context) {
	var req struct {
		HabitID int `json:"habit_id" binding:"required"`
		ChildID int `json:"child_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignment := &models.HabitAssignment{
		HabitID: req.HabitID,
		ChildID: req.ChildID,
		Status:  1,
	}

	if err := models.AssignHabit(assignment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"assignment": assignment})
}

// GetChildHabits 获取小孩的习惯列表
func GetChildHabits(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	habits, err := models.GetHabitsByChildID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"habits": habits})
}