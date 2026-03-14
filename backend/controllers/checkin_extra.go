package controllers

import (
	"babyhabit/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCheckinRecords 获取打卡记录
func GetCheckinRecords(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	habitID := c.Query("habit_id")

	records, err := models.GetCheckinRecords(user.ID, startDate, endDate, habitID)
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

	var habitID int64
	if _, err := strconv.ParseInt(habitIDStr, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid habit_id"})
		return
	}

	streak, err := models.GetStreakRecord(user.ID, habitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"streak": streak})
}
