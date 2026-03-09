package controllers

import (
	"babyhabit/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Checkin 打卡
func Checkin(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var req struct {
		HabitID      int64  `json:"habit_id" binding:"required"`
		AssignmentID int64  `json:"assignment_id"`
		CheckinDate  string `json:"checkin_date" binding:"required"`
		Remark       string `json:"remark"`
		CheckinType  int    `json:"checkin_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查是否已经打卡
	existingRecord, err := models.GetCheckinRecord(user.ID, req.HabitID, req.CheckinDate)
	if err == nil && existingRecord != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already checked in for this date"})
		return
	}

	// 获取习惯信息
	habit, err := models.GetHabitByID(req.HabitID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	// 验证打卡时间
	if req.CheckinType == 1 { // 正常打卡
		currentTime := time.Now().Format("15:04")
		if habit.CheckinTimeStart != "" && currentTime < habit.CheckinTimeStart {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Checkin time is too early"})
			return
		}
		if habit.CheckinTimeEnd != "" && currentTime > habit.CheckinTimeEnd {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Checkin time is too late"})
			return
		}
	}

	// 验证补卡
	if req.CheckinType == 2 { // 补卡
		if habit.AllowMakeup != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Makeup checkin is not allowed for this habit"})
			return
		}

		// 检查补卡时间范围
		checkinDate, _ := time.Parse("2006-01-02", req.CheckinDate)
		today := time.Now()
		dayDiff := int(today.Sub(checkinDate).Hours() / 24)
		if dayDiff > habit.MakeupDays {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Makeup checkin is only allowed within " + strconv.Itoa(habit.MakeupDays) + " days"})
			return
		}
	}

	// 创建打卡记录
	record := &models.CheckinRecord{
		UserID:        user.ID,
		HabitID:       req.HabitID,
		AssignmentID:  req.AssignmentID,
		CheckinDate:   req.CheckinDate,
		CheckinTime:   time.Now(),
		CheckinType:   req.CheckinType,
		Remark:        req.Remark,
		Status:        1,
		PointsRewarded: habit.RewardPoints,
	}

	if err := models.CreateCheckinRecord(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新连续打卡记录
	if err := models.UpdateStreakRecord(user.ID, req.HabitID, req.CheckinDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update streak record: " + err.Error()})
		return
	}

	// 发放积分
	if err := models.AddPoints(user.ID, habit.RewardPoints, "checkin", "checkin_record", record.ID, time.Time{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add points: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"record": record})
}

// GetCheckinRecords 获取打卡记录
func GetCheckinRecords(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date and end date are required"})
		return
	}

	records, err := models.GetCheckinRecordsByUserID(user.ID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"records": records})
}

// GetStreakRecord 获取连续打卡记录
func GetStreakRecord(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	habitIDStr := c.Param("habit_id")
	habitID, err := strconv.ParseInt(habitIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid habit ID"})
		return
	}

	streak, err := models.GetStreakRecord(user.ID, habitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"streak": streak})
}