package dao

import (
	"lyworder/global"
	"lyworder/internal/model"
	"time"
)

func CreateTicket(ticket *model.Ticket) error {
	return global.DBEngine.Create(ticket).Error
}

func UpdateTicketStatus(ticketID, status string) error {
	now := time.Now()
	return global.DBEngine.Model(&model.Ticket{}).
		Where("id = ?", ticketID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": now,
		}).Error
}

// 新增UpdateTicketOperatorName方法
func UpdateTicketOperatorName(ticketID string, operatorName string) error {
	now := time.Now()
	return global.DBEngine.Model(&model.Ticket{}).
		Where("id = ?", ticketID).
		Updates(map[string]interface{}{
			"operator_name": operatorName,
			"updated_at":    now,
		}).Error
}

// 修改GetTickets方法，确保返回operator_name字段并支持查询所有
func GetTickets(title, status, ticketType, startDate, endDate, userName string, operatorName string, page, limit int) ([]model.Ticket, int64, error) {
	db := global.DBEngine.Model(&model.Ticket{}).Select("*")
	if userName != "" {
		db = db.Where("user_name LIKE ?", userName+"%")
	}
	if operatorName != "" {
		db = db.Where("operator_name LIKE ?", "%"+operatorName+"%")
	}

	var tickets []model.Ticket
	var total int64

	if title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if ticketType != "" {
		db = db.Where("ticket_type = ?", ticketType)
	}

	if startDate != "" && endDate != "" {
		db = db.Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate)
	} else if startDate != "" {
		db = db.Where("DATE(created_at) >= ?", startDate)
	} else if endDate != "" {
		db = db.Where("DATE(created_at) <= ?", endDate)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 应用排序和分页
	if page > 0 && limit > 0 {
		db = db.Order("created_at desc")
		offset := (page - 1) * limit
		err := db.Offset(offset).Limit(limit).Find(&tickets).Error
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := db.Order("created_at desc").Find(&tickets).Error
		if err != nil {
			return nil, 0, err
		}
	}

	return tickets, total, nil
}

func GetTicketDetail(ticketID string) (*model.Ticket, error) {
	var ticket model.Ticket
	err := global.DBEngine.Model(&model.Ticket{}).
		Where("id = ?", ticketID).
		First(&ticket).Error

	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// UpdateTicketType 更新工单类型
func UpdateTicketType(ticketID string, ticketType string) error {
	now := time.Now()
	return global.DBEngine.Model(&model.Ticket{}).
		Where("id = ?", ticketID).
		Updates(map[string]interface{}{
			"ticket_type": ticketType,
			"updated_at":  now,
		}).Error
}

// UpdateTicketRemark 更新工单备注
func UpdateTicketRemark(ticketID string, remark string) error {
	now := time.Now()
	return global.DBEngine.Model(&model.Ticket{}).
		Where("id = ?", ticketID).
		Updates(map[string]interface{}{
			"remark":     remark,
			"updated_at": now,
		}).Error
}

// UpdateTicketSolution 更新工单解决方案
func UpdateTicketSolution(ticketID string, solution string) error {
	now := time.Now()
	return global.DBEngine.Model(&model.Ticket{}).
		Where("id = ?", ticketID).
		Updates(map[string]interface{}{
			"solution":   solution,
			"updated_at": now,
		}).Error
}
