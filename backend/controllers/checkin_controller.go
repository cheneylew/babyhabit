package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cheneylew/babyhabit/backend/models"
	"github.com/cheneylew/babyhabit/backend/utils"

	"github.com/gin-gonic/gin"
)

// Checkin 打卡
func Checkin(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// 获取打卡信息
	habitIDStr := c.PostForm("habit_id")
	assignmentIDStr := c.PostForm("assignment_id")
	checkinDateStr := c.PostForm("checkin_date")
	remark := c.PostForm("remark")
	checkinTypeStr := c.PostForm("checkin_type")

	// 调试日志
	log.Printf("Received form data: habit_id=%s, assignment_id=%s, checkin_date=%s, checkin_type=%s",
		habitIDStr, assignmentIDStr, checkinDateStr, checkinTypeStr)

	// 解析字段
	var habitID int64
	var assignmentID int64
	var checkinType int
	var err error

	if habitIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "habit_id is required"})
		return
	}

	habitID, err = strconv.ParseInt(habitIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid habit_id"})
		return
	}

	if assignmentIDStr != "" {
		assignmentID, err = strconv.ParseInt(assignmentIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment_id"})
			return
		}
	}

	if checkinTypeStr != "" {
		checkinType, err = strconv.Atoi(checkinTypeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid checkin_type"})
			return
		}
	} else {
		checkinType = 1 // 默认为正常打卡
	}

	// 处理上传的照片
	var photoURL string
	file, err := c.FormFile("photo")
	if err == nil && file != nil {
		// 保存到 files 目录
		saveDir := "files/checkin_photos"
		photoPath, err := utils.SaveUploadedFile(saveDir, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "上传照片失败：" + err.Error()})
			return
		}
		photoURL = photoPath
	}

	// 如果没有传入打卡日期，默认为今天
	if checkinDateStr == "" {
		checkinDateStr = time.Now().Format("2006-01-02")
	}

	// 检查是否已经打卡
	existingRecord, err := models.GetCheckinRecord(user.ID, habitID, checkinDateStr)
	if err == nil && existingRecord != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already checked in for this date"})
		return
	}

	// 获取习惯信息
	habit, err := models.GetHabitByID(habitID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	// 验证打卡时间
	if checkinType == 1 { // 正常打卡
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

	// 验证打卡日期
	checkinDate, err := time.Parse("2006-01-02", checkinDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid checkin date format"})
		return
	}

	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	checkinDate = time.Date(checkinDate.Year(), checkinDate.Month(), checkinDate.Day(), 0, 0, 0, 0, checkinDate.Location())

	// 计算日期差异
	dayDiff := int(today.Sub(checkinDate).Hours() / 24)

	// 明天及之后不允许打卡
	if dayDiff < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot check in for future dates"})
		return
	}

	// 昨天及更早的日期，只有允许补卡的习惯才能打卡
	if dayDiff > 0 {
		if habit.AllowMakeup != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Makeup checkin is not allowed for this habit"})
			return
		}

		// 只允许补卡最近 N 天（根据 makeup_days 配置）
		if dayDiff > habit.MakeupDays {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Makeup checkin is only allowed within " + strconv.Itoa(habit.MakeupDays) + " days"})
			return
		}
	}

	// 验证周期性习惯
	if habit.ScheduleType == 2 { // 周期性习惯
		if habit.ScheduleDetail == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Periodic habit has no schedule detail"})
			return
		}

		// 检查打卡日期是否在允许的星期内
		// schedule_detail 格式：逗号分隔的星期几 (0-6, 0 表示周日)
		allowedDays := strings.Split(habit.ScheduleDetail, ",")
		checkinWeekday := strconv.Itoa(int(checkinDate.Weekday()))

		found := false
		for _, day := range allowedDays {
			if day == checkinWeekday {
				found = true
				break
			}
		}

		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Today is not in the scheduled days for this habit"})
			return
		}
	}

	// 验证补卡
	if checkinType == 2 { // 补卡
		if habit.AllowMakeup != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Makeup checkin is not allowed for this habit"})
			return
		}

		// 检查补卡时间范围
		if dayDiff > habit.MakeupDays {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Makeup checkin is only allowed within " + strconv.Itoa(habit.MakeupDays) + " days"})
			return
		}
	}

	// 创建打卡记录
	record := &models.CheckinRecord{
		UserID:         user.ID,
		HabitID:        habitID,
		AssignmentID:   assignmentID,
		CheckinDate:    checkinDateStr,
		CheckinTime:    time.Now(),
		CheckinType:    checkinType,
		PhotoURL:       photoURL,
		Remark:         remark,
		Status:         1,
		PointsRewarded: habit.RewardPoints,
	}

	if err := models.CreateCheckinRecord(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新连续打卡记录
	if err := models.UpdateStreakRecord(user.ID, habitID, checkinDateStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update streak record: " + err.Error()})
		return
	}

	// 发放积分
	if err := models.AddPoints(user.ID, habit.RewardPoints, "checkin", "checkin_record", record.ID, time.Time{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add points: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Checkin successful",
		"record":  record,
		"points":  habit.RewardPoints,
	})
}

// SubmitSelfRate 提交自我评分
func SubmitSelfRate(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	println("=== SubmitSelfRate Called ===")
	println("User ID:", user.ID)

	var req struct {
		RecordID int64 `json:"record_id" binding:"required"`
		SelfRate int   `json:"self_rate" binding:"required,min=1,max=10"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		println("Bind error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	println("RecordID:", req.RecordID)
	println("SelfRate:", req.SelfRate)

	// 检查打卡记录是否存在且属于当前用户
	println("=== Getting Checkin Record ===")
	record, err := models.GetCheckinRecordByID(req.RecordID)
	if err != nil {
		println("GetCheckinRecordByID error:", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Checkin record not found"})
		return
	}
	println("Record found, UserID:", record.UserID)

	if record.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only rate your own checkin records"})
		return
	}

	// 更新自我评分
	if err := models.UpdateCheckinSelfRate(req.RecordID, req.SelfRate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update self rate: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Self rate submitted successfully",
		"record_id": req.RecordID,
		"self_rate": req.SelfRate,
	})
}

// GetDailyPointsStats 获取用户每天的积分统计
func GetDailyPointsStats(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// 获取查询参数
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	// 获取每天的积分统计
	stats, err := models.GetDailyPointsStats(user.ID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get daily points stats: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}
