package models

import (
	"babyhabit/config"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type RewardItem struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Image          string     `json:"image"`
	Category       string     `json:"category"`
	PointsRequired int        `json:"points_required"`
	Stock          int        `json:"stock"`
	ExchangeLimit  int        `json:"exchange_limit"`
	UserExchanged  int        `json:"user_exchanged"`
	StartTime      *time.Time `json:"start_time"`
	EndTime        *time.Time `json:"end_time"`
	CreatorID      int64      `json:"creator_id"`
	CreateTime     time.Time  `json:"create_time"`
	UpdateTime     time.Time  `json:"update_time"`
	Status         int        `json:"status"`
}

type ExchangeRecord struct {
	ID           int64       `json:"id"`
	UserID       int64       `json:"user_id"`
	ItemID       int64       `json:"item_id"`
	Points       int         `json:"points"`
	Quantity     int         `json:"quantity"`
	ExchangeTime time.Time   `json:"exchange_time"`
	DeliveryInfo *string     `json:"delivery_info"`
	Status       int         `json:"status"`
	CreateTime   time.Time   `json:"create_time"`
	UpdateTime   time.Time   `json:"update_time"`
	Item         *RewardItem `json:"item,omitempty"`
}

// CreateRewardItem 创建奖励物品
func CreateRewardItem(item *RewardItem) error {
	query := `INSERT INTO reward_item (name, description, image, category, points_required, 
			 stock, exchange_limit, user_exchanged, start_time, end_time, creator_id, 
			 create_time, update_time, status) 
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), ?)`
	result, err := config.DB.Exec(query, item.Name, item.Description, item.Image, item.Category,
		item.PointsRequired, item.Stock, item.ExchangeLimit, item.UserExchanged,
		item.StartTime, item.EndTime, item.CreatorID, item.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	item.ID = id
	return nil
}

// GetRewardItemByID 根据ID获取奖励物品
func GetRewardItemByID(id int64) (*RewardItem, error) {
	item := &RewardItem{}
	query := `SELECT id, name, description, image, category, points_required, stock, 
			 exchange_limit, user_exchanged, start_time, end_time, creator_id, 
			 create_time, update_time, status 
			 FROM reward_item WHERE id = ?`
	err := config.DB.QueryRow(query, id).Scan(
		&item.ID, &item.Name, &item.Description, &item.Image, &item.Category, &item.PointsRequired,
		&item.Stock, &item.ExchangeLimit, &item.UserExchanged, &item.StartTime, &item.EndTime,
		&item.CreatorID, &item.CreateTime, &item.UpdateTime, &item.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("reward item not found")
		}
		return nil, err
	}
	return item, nil
}

// GetRewardItems 获取奖励物品列表
func GetRewardItems(status int) ([]*RewardItem, error) {
	query := `SELECT id, name, description, image, category, points_required, stock, 
			 exchange_limit, user_exchanged, start_time, end_time, creator_id, 
			 create_time, update_time, status 
			 FROM reward_item 
			 WHERE status = ? 
			 ORDER BY create_time DESC`
	rows, err := config.DB.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*RewardItem{}
	for rows.Next() {
		item := &RewardItem{}
		err := rows.Scan(
			&item.ID, &item.Name, &item.Description, &item.Image, &item.Category, &item.PointsRequired,
			&item.Stock, &item.ExchangeLimit, &item.UserExchanged, &item.StartTime, &item.EndTime,
			&item.CreatorID, &item.CreateTime, &item.UpdateTime, &item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// GetRewardItemsByCreatorID 根据创建者ID获取奖励物品列表
func GetRewardItemsByCreatorID(creatorID int64, status int) ([]*RewardItem, error) {
	query := `SELECT id, name, description, image, category, points_required, stock, 
			 exchange_limit, user_exchanged, start_time, end_time, creator_id, 
			 create_time, update_time, status 
			 FROM reward_item 
			 WHERE creator_id = ? AND status = ? 
			 ORDER BY create_time DESC`
	rows, err := config.DB.Query(query, creatorID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*RewardItem{}
	for rows.Next() {
		item := &RewardItem{}
		err := rows.Scan(
			&item.ID, &item.Name, &item.Description, &item.Image, &item.Category, &item.PointsRequired,
			&item.Stock, &item.ExchangeLimit, &item.UserExchanged, &item.StartTime, &item.EndTime,
			&item.CreatorID, &item.CreateTime, &item.UpdateTime, &item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// UpdateRewardItem 更新奖励物品
func UpdateRewardItem(item *RewardItem) error {
	query := `UPDATE reward_item SET name = ?, description = ?, image = ?, category = ?, 
			 points_required = ?, stock = ?, exchange_limit = ?, user_exchanged = ?, 
			 start_time = ?, end_time = ?, status = ?, update_time = NOW() 
			 WHERE id = ?`
	_, err := config.DB.Exec(query, item.Name, item.Description, item.Image, item.Category,
		item.PointsRequired, item.Stock, item.ExchangeLimit, item.UserExchanged,
		item.StartTime, item.EndTime, item.Status, item.ID)
	return err
}

// CreateExchangeRecord 创建兑换记录
func CreateExchangeRecord(record *ExchangeRecord) error {
	// 开始事务
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 先查询奖励物品的库存和基本信息
	var stock int
	var itemName, itemImage, itemCategory string
	stockQuery := `SELECT stock, name, image, category FROM reward_item WHERE id = ?`
	err = tx.QueryRow(stockQuery, record.ItemID).Scan(&stock, &itemName, &itemImage, &itemCategory)
	if err != nil {
		return err
	}

	// 插入兑换记录
	// 如果 delivery_info 为空，则使用 NULL
	var deliveryInfo interface{}
	if record.DeliveryInfo == nil || *record.DeliveryInfo == "" {
		deliveryInfo = nil
	} else {
		deliveryInfo = *record.DeliveryInfo
	}

	query := `INSERT INTO exchange_record (user_id, item_id, item_name, item_image, item_category, points, quantity, exchange_time, 
			 delivery_info, status, create_time, update_time) 
			 VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), ?, ?, NOW(), NOW())`
	result, err := tx.Exec(query, record.UserID, record.ItemID, itemName, itemImage, itemCategory, record.Points, record.Quantity,
		deliveryInfo, record.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	record.ID = id

	// 更新奖励物品库存和已兑换数量
	// 如果库存为 -1（无限库存），则不减少库存
	var updateQuery string
	if stock == -1 {
		updateQuery = `UPDATE reward_item SET user_exchanged = user_exchanged + ? 
				 WHERE id = ?`
		_, err = tx.Exec(updateQuery, record.Quantity, record.ItemID)
	} else {
		updateQuery = `UPDATE reward_item SET stock = stock - ?, user_exchanged = user_exchanged + ? 
				 WHERE id = ?`
		_, err = tx.Exec(updateQuery, record.Quantity, record.Quantity, record.ItemID)
	}
	if err != nil {
		return err
	}

	// 在同一个事务中扣除用户积分
	// 1. 检查积分余额
	var balance int
	checkQuery := `SELECT points_balance FROM user WHERE id = ?`
	err = tx.QueryRow(checkQuery, record.UserID).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < record.Points {
		return ErrInsufficientPoints
	}

	// 2. 插入积分记录（负数）
	pointsQuery := `INSERT INTO points_record (user_id, points, source, related_type, related_id, create_time) 
			 VALUES (?, ?, ?, ?, ?, NOW())`
	_, err = tx.Exec(pointsQuery, record.UserID, -record.Points, "exchange", "exchange_record", id)
	if err != nil {
		return err
	}

	// 3. 更新用户积分余额
	updateBalanceQuery := `UPDATE user SET points_balance = points_balance - ? WHERE id = ?`
	_, err = tx.Exec(updateBalanceQuery, record.Points, record.UserID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetExchangeRecordsByUserID 获取用户的兑换记录
func GetExchangeRecordsByUserID(userID int64) ([]*ExchangeRecord, error) {
	query := `SELECT er.id, er.user_id, er.item_id, er.item_name, er.item_image, er.item_category, er.points, er.quantity, er.exchange_time, 
			 er.delivery_info, er.status, er.create_time, er.update_time 
			 FROM exchange_record er 
			 WHERE er.user_id = ? 
			 ORDER BY er.create_time DESC`
	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []*ExchangeRecord{}
	for rows.Next() {
		record := &ExchangeRecord{}
		item := &RewardItem{}
		err := rows.Scan(
			&record.ID, &record.UserID, &record.ItemID, &item.Name, &item.Image, &item.Category, &record.Points, &record.Quantity,
			&record.ExchangeTime, &record.DeliveryInfo, &record.Status, &record.CreateTime,
			&record.UpdateTime,
		)
		if err != nil {
			return nil, err
		}
		record.Item = item
		records = append(records, record)
	}
	return records, nil
}

// GetAllExchangeRecords 获取所有兑换记录（管理员）
func GetAllExchangeRecords() ([]*ExchangeRecord, error) {
	query := `SELECT er.id, er.user_id, er.item_id, er.item_name, er.item_image, er.item_category, er.points, er.quantity, er.exchange_time, 
			 er.delivery_info, er.status, er.create_time, er.update_time,
			 u.name as user_name
			 FROM exchange_record er 
			 JOIN user u ON er.user_id = u.id
			 ORDER BY er.create_time DESC`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []*ExchangeRecord{}
	for rows.Next() {
		record := &ExchangeRecord{}
		item := &RewardItem{}
		var userName string
		err := rows.Scan(
			&record.ID, &record.UserID, &record.ItemID, &item.Name, &item.Image, &item.Category, &record.Points, &record.Quantity,
			&record.ExchangeTime, &record.DeliveryInfo, &record.Status, &record.CreateTime,
			&record.UpdateTime, &userName,
		)
		if err != nil {
			return nil, err
		}
		record.Item = item
		// 将用户名放入 Item 的 Description 字段中（临时方案）
		item.Description = userName
		records = append(records, record)
	}
	return records, nil
}

// GetExchangeRecordsByUserIDs 根据用户ID列表获取兑换记录
func GetExchangeRecordsByUserIDs(userIDs []int64) ([]*ExchangeRecord, error) {
	// 构建IN查询的参数
	placeholders := make([]string, len(userIDs))
	args := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `SELECT er.id, er.user_id, er.item_id, er.item_name, er.item_image, er.item_category, er.points, er.quantity, er.exchange_time, 
			 er.delivery_info, er.status, er.create_time, er.update_time,
			 u.name as user_name
			 FROM exchange_record er 
			 JOIN user u ON er.user_id = u.id
			 WHERE er.user_id IN (` + strings.Join(placeholders, ",") + `)
			 ORDER BY er.create_time DESC`
	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []*ExchangeRecord{}
	for rows.Next() {
		record := &ExchangeRecord{}
		item := &RewardItem{}
		var userName string
		err := rows.Scan(
			&record.ID, &record.UserID, &record.ItemID, &item.Name, &item.Image, &item.Category, &record.Points, &record.Quantity,
			&record.ExchangeTime, &record.DeliveryInfo, &record.Status, &record.CreateTime,
			&record.UpdateTime, &userName,
		)
		if err != nil {
			return nil, err
		}
		record.Item = item
		// 将用户名放入 Item 的 Description 字段中（临时方案）
		item.Description = userName
		records = append(records, record)
	}
	return records, nil
}

// UpdateExchangeStatus 更新兑换状态
func UpdateExchangeStatus(id int64, status int) error {
	query := `UPDATE exchange_record SET status = ?, update_time = NOW() WHERE id = ?`
	_, err := config.DB.Exec(query, status, id)
	return err
}

// DeleteRewardItem 删除奖励物品
func DeleteRewardItem(id int64) error {
	// 直接删除奖励物品，不再删除相关的兑换记录
	// 因为兑换记录已经包含了冗余字段，与奖励管理解耦
	rewardQuery := `DELETE FROM reward_item WHERE id = ?`
	_, err := config.DB.Exec(rewardQuery, id)
	return err
}

// DeleteExchangeRecord 删除兑换记录
func DeleteExchangeRecord(id int64) error {
	query := `DELETE FROM exchange_record WHERE id = ?`
	_, err := config.DB.Exec(query, id)
	return err
}
