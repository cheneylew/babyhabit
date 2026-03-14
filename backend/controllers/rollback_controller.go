package controllers

import (
	"babyhabit/config"
	"babyhabit/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RollbackCheckin 回退打卡
func RollbackCheckin(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// 检查是否是管理员
	if user.UserType != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can rollback checkins"})
		return
	}

	// 获取请求参数
	var req struct {
		CheckinID int64  `json:"checkin_id"`
		Reason    string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.CheckinID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "checkin_id is required"})
		return
	}

	// 开始事务
	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 获取打卡记录
	checkin := &models.CheckinRecord{}
	query := `SELECT id, user_id, habit_id, points_rewarded, status, is_rolled_back 
			  FROM checkin_record WHERE id = ?`
	err = tx.QueryRow(query, req.CheckinID).Scan(
		&checkin.ID, &checkin.UserID, &checkin.HabitID,
		&checkin.PointsRewarded, &checkin.Status, &checkin.IsRolledBack,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Checkin record not found"})
		return
	}

	// 检查是否已经回退
	if checkin.IsRolledBack == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Checkin already rolled back"})
		return
	}

	// 1. 回退积分
	err = deductPointsInTx(tx, checkin.UserID, int(checkin.PointsRewarded),
		"rollback", "checkin_record", checkin.ID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rollback points: " + err.Error()})
		return
	}

	// 2. 更新打卡记录状态
	now := time.Now()
	updateQuery := `UPDATE checkin_record 
				   SET is_rolled_back = 1, rollback_time = ?, rollback_reason = ?, status = 0 
				   WHERE id = ?`
	_, err = tx.Exec(updateQuery, now, req.Reason, checkin.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update checkin record: " + err.Error()})
		return
	}

	// 3. 更新连续打卡记录（重新计算）
	err = recalculateStreak(tx, checkin.UserID, checkin.HabitID)
	if err != nil {
		// 回退失败不阻止主流程，只记录日志
		println("Warning: Failed to recalculate streak:", err.Error())
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Checkin rolled back successfully",
		"checkin_id":      checkin.ID,
		"points_deducted": checkin.PointsRewarded,
	})
}

// GetChildCheckinRecords 获取小孩的打卡记录（管理员）
func GetChildCheckinRecords(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// 检查是否是管理员
	if user.UserType != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can view child checkin records"})
		return
	}

	// 获取参数
	childIDStr := c.Query("child_id")
	habitIDStr := c.Query("habit_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if childIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "child_id is required"})
		return
	}

	childID, err := strconv.ParseInt(childIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid child_id"})
		return
	}

	// 查询打卡记录
	records, err := models.GetChildCheckinRecords(childID, startDate, endDate, habitIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"records": records})
}

// deductPointsInTx 在事务内扣除积分
func deductPointsInTx(tx *sql.Tx, userID int64, points int, source, relatedType string, relatedID int64, reason string) error {
	// 插入积分记录
	insertQuery := `INSERT INTO points_record (user_id, points, source, related_type, related_id, create_time, expire_time) 
				   VALUES (?, ?, ?, ?, ?, NOW(), ?)`

	_, err := tx.Exec(insertQuery, userID, -points, source, relatedType, relatedID, nil)
	if err != nil {
		return err
	}

	// 更新用户积分余额
	updateQuery := `UPDATE user SET points_balance = points_balance - ? WHERE id = ?`
	_, err = tx.Exec(updateQuery, points, userID)
	if err != nil {
		return err
	}

	return nil
}

// recalculateStreak 重新计算连续打卡记录
func recalculateStreak(tx *sql.Tx, userID, habitID int64) error {
	// 查询该习惯的所有有效打卡记录（未回退的）
	query := `SELECT checkin_date FROM checkin_record 
			  WHERE user_id = ? AND habit_id = ? AND is_rolled_back = 0 
			  ORDER BY checkin_date DESC`

	rows, err := tx.Query(query, userID, habitID)
	if err != nil {
		return err
	}
	defer rows.Close()

	dates := []string{}
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			return err
		}
		dates = append(dates, date)
	}

	// 计算连续打卡天数
	currentStreak := 0
	longestStreak := 0

	if len(dates) > 0 {
		currentStreak = 1
		longestStreak = 1

		for i := 1; i < len(dates); i++ {
			prevDate, _ := time.Parse("2006-01-02", dates[i-1])
			currDate, _ := time.Parse("2006-01-02", dates[i])

			diff := int(prevDate.Sub(currDate).Hours() / 24)
			if diff == 1 {
				currentStreak++
				if currentStreak > longestStreak {
					longestStreak = currentStreak
				}
			} else {
				currentStreak = 1
			}
		}
	}

	// 更新连续打卡记录
	streakStartDate := ""
	if len(dates) > 0 {
		streakStartDate = dates[len(dates)-1]
	}

	lastCheckinDate := ""
	if len(dates) > 0 {
		lastCheckinDate = dates[0]
	}

	// 检查是否存在连续打卡记录
	var streakID int64
	checkQuery := `SELECT id FROM streak_record WHERE user_id = ? AND habit_id = ?`
	err = tx.QueryRow(checkQuery, userID, habitID).Scan(&streakID)

	if err == sql.ErrNoRows {
		// 插入新记录
		insertQuery := `INSERT INTO streak_record (user_id, habit_id, current_streak, longest_streak, 
						   last_checkin_date, streak_start_date, update_time) 
					   VALUES (?, ?, ?, ?, ?, ?, NOW())`
		_, err = tx.Exec(insertQuery, userID, habitID, currentStreak, longestStreak, lastCheckinDate, streakStartDate)
	} else if err == nil {
		// 更新现有记录
		updateQuery := `UPDATE streak_record 
					   SET current_streak = ?, longest_streak = ?, last_checkin_date = ?, 
						   streak_start_date = ?, update_time = NOW() 
					   WHERE id = ?`
		_, err = tx.Exec(updateQuery, currentStreak, longestStreak, lastCheckinDate, streakStartDate, streakID)
	}

	return err
}
