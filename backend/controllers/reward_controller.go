package controllers

import (
	"babyhabit/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateRewardItem 创建奖励物品
func CreateRewardItem(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var req struct {
		Name           string `json:"name" binding:"required"`
		Description    string `json:"description"`
		Image          string `json:"image"`
		Category       string `json:"category"`
		PointsRequired int    `json:"points_required" binding:"required,min=1"`
		Stock          int    `json:"stock"`
		ExchangeLimit  int    `json:"exchange_limit"`
		Status         int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := &models.RewardItem{
		Name:           req.Name,
		Description:    req.Description,
		Image:          req.Image,
		Category:       req.Category,
		PointsRequired: req.PointsRequired,
		Stock:          req.Stock,
		ExchangeLimit:  req.ExchangeLimit,
		CreatorID:      user.ID,
		Status:         req.Status,
	}

	if err := models.CreateRewardItem(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": item})
}

// GetRewardItems 获取奖励物品列表
func GetRewardItems(c *gin.Context) {
	statusStr := c.DefaultQuery("status", "1")
	status, _ := strconv.Atoi(statusStr)

	// 获取当前用户信息
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	currentUser := user.(*models.User)
	userID := currentUser.ID

	// 如果是小孩用户，使用父母ID查询奖励
	if currentUser.UserType == 2 && currentUser.ParentID > 0 {
		userID = currentUser.ParentID
	}

	items, err := models.GetRewardItemsByCreatorID(userID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// GetRewardItem 获取奖励物品详情
func GetRewardItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := models.GetRewardItemByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": item})
}

// UpdateRewardItem 更新奖励物品
func UpdateRewardItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := models.GetRewardItemByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
		Image          string `json:"image"`
		Category       string `json:"category"`
		PointsRequired int    `json:"points_required"`
		Stock          int    `json:"stock"`
		ExchangeLimit  int    `json:"exchange_limit"`
		Status         int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新物品信息
	item.Name = req.Name
	item.Description = req.Description
	item.Image = req.Image
	item.Category = req.Category
	item.PointsRequired = req.PointsRequired
	item.Stock = req.Stock
	item.ExchangeLimit = req.ExchangeLimit
	item.Status = req.Status

	if err := models.UpdateRewardItem(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": item})
}

// Exchange 兑换奖励
func Exchange(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var req struct {
		ItemID       int64  `json:"item_id" binding:"required"`
		Quantity     int    `json:"quantity" binding:"required,min=1"`
		DeliveryInfo string `json:"delivery_info"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取奖励物品
	item, err := models.GetRewardItemByID(req.ItemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 检查库存
	if item.Stock != -1 && item.Stock < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	// 检查兑换限制
	if item.ExchangeLimit > 0 {
		// 这里需要查询用户已兑换数量，简化处理
	}

	// 检查积分余额
	if user.PointsBalance < item.PointsRequired*req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient points"})
		return
	}

	// 创建兑换记录
	record := &models.ExchangeRecord{
		UserID:   user.ID,
		ItemID:   req.ItemID,
		Points:   item.PointsRequired * req.Quantity,
		Quantity: req.Quantity,
		Status:   2, // 处理中
	}
	// 如果 DeliveryInfo 不为空，设置指针
	if req.DeliveryInfo != "" {
		record.DeliveryInfo = &req.DeliveryInfo
	}

	if err := models.CreateExchangeRecord(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"record": record})
}

// GetExchangeRecords 获取兑换记录
func GetExchangeRecords(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	records, err := models.GetExchangeRecordsByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"records": records})
}

// GetAllExchangeRecords 获取所有兑换记录（管理员）
func GetAllExchangeRecords(c *gin.Context) {
	// 获取当前用户信息
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	userID := user.(*models.User).ID

	// 获取用户的孩子ID列表
	childIDs, err := models.GetChildIDsByParentID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建用户ID列表（包括自己和孩子）
	userIDs := append([]int64{userID}, childIDs...)

	records, err := models.GetExchangeRecordsByUserIDs(userIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"records": records})
}

// UpdateExchangeStatus 更新兑换状态
func UpdateExchangeStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=1 2"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateExchangeStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exchange status updated successfully"})
}

// DeleteRewardItem 删除奖励物品
func DeleteRewardItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	if err := models.DeleteRewardItem(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reward item deleted successfully"})
}

// DeleteExchangeRecord 删除兑换记录
func DeleteExchangeRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	if err := models.DeleteExchangeRecord(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exchange record deleted successfully"})
}
