package v1

import (
	"fmt"
	"log"
	"lyworder/global"
	"lyworder/internal/dao"
	"lyworder/internal/model"
	"lyworder/pkg/dify"
	"lyworder/pkg/meet"
	"lyworder/pkg/utils"

	"net/http"
	"strconv"
	"strings"
	"time"

	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// CreateTicket 创建工单
// @Summary 创建新工单
// @Description 创建新的工单，支持截图上传
// @Tags 工单管理
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "工单标题"
// @Param content formData string true "工单内容"
// @Param ticket_type formData string false "工单类型"
// @Param screenshots formData []file false "截图文件"
// @Success 201 {object} model.Response{data=map[string]interface{}} "创建成功"
// @Failure 400 {object} model.Response{error=string} "请求参数错误"
// @Failure 401 {object} model.Response{error=string} "未授权"
// @Failure 500 {object} model.Response{error=string} "服务器错误"
// @Router /api/v1/tickets [post]
func CreateTicket(c *gin.Context) {
	var ticket model.Ticket
	var robot meet.Robot
	// 解析多部分表单数据（最大10MB）
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析表单数据失败: " + err.Error()})
		return
	}

	// 提取表单字段
	ticket.Title = c.PostForm("title")
	ticket.Content = c.PostForm("content")
	ticket.TicketType = c.PostForm("ticket_type")

	// 验证必填字段
	if ticket.Title == "" || ticket.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题和内容为必填项"})
		return
	}

	// 从Cookie获取JWT令牌
	tokenString, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少访问令牌"})
		return
	}

	// 解析JWT令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(global.OauthSetting.ClientSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	// 从JWT声明中获取用户信息，处理可能的float64类型
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(float64); ok {
			ticket.UserID = fmt.Sprintf("%.0f", userID)
		} else if userID, ok := claims["user_id"].(string); ok {
			ticket.UserID = userID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID类型"})
			return
		}

		if userName, ok := claims["user_name"].(string); ok {
			ticket.UserName = userName
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户名类型"})
			return
		}
	}

	ticket.ID = time.Now().Format("20060102150405")
	ticket.CreatedAt = time.Now()
	ticket.Status = "open"
	ticket.ScreenShots = []string{}

	// 处理多张截图上传
	files, ok := c.Request.MultipartForm.File["screenshots"]
	if ok {
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无法打开文件: " + err.Error()})
				return
			}
			defer file.Close()

			// 验证文件大小 (5MB)
			if fileHeader.Size > 5*1024*1024 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过5MB"})
				return
			}

			// 创建上传目录
			uploadDir := filepath.Join("static/uploads")
			if err = os.MkdirAll(uploadDir, os.ModePerm); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败: " + err.Error()})
				return
			}

			// 生成唯一文件名
			ext := filepath.Ext(fileHeader.Filename)
			filename := fmt.Sprintf("%s_%d%s", ticket.ID, time.Now().UnixNano(), ext)
			filePath := filepath.Join(uploadDir, filename)

			// 保存文件
			dst, err := os.Create(filePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文件失败: " + err.Error()})
				return
			}
			defer dst.Close()
			if utils.IsImageFile(file) {
				file.Seek(0, 0)
				if _, err = io.Copy(dst, file); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
					return
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件格式"})
				return
			}

			err = global.ThumbnailGenerator.GenerateThumbnail(filePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "创建缩略图失败: " + err.Error()})
				return
			}
			// 添加截图路径到工单
			screenshotPath := "/uploads/" + filename
			ticket.ScreenShots = append(ticket.ScreenShots, screenshotPath)
		}
	}

	// 创建工单
	if err = dao.CreateTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建工单失败: " + err.Error()})
		return
	}

	// 异步执行 AI 分析和更新操作
	go func() {
		aiticket, err := dify.ScheduleTicket(&ticket)
		if err != nil {
			log.Printf("ScheduleTicket failed, err: %v", err)
			return
		}
		dao.UpdateTicketRemark(ticket.ID, aiticket.Remark)
		dao.UpdateTicketType(ticket.ID, aiticket.TicketType)
		if aiticket.OperatorName != "" {
			robot.ToUserName = []string{aiticket.OperatorName}
			ticketdetail := fmt.Sprintf("%s%s?id=%s", global.ServerSetting.Host, "/tickets/detail", ticket.ID)
			robot.Content = fmt.Sprintf("### %s\n%s\n工单详情：[%s](%s)\n此工单已指派给你", ticket.Title, ticket.Content, ticketdetail, ticketdetail)
			err = meet.Send(&robot)
			if err != nil {
				log.Printf("SendTicketOperator failed, err: %v", err)
			}
			robot.ToUserName = []string{ticket.UserName}
			robot.Content = fmt.Sprintf("### %s\n%s\n工单详情：[%s](%s)\n处理人: **%s**", ticket.Title, ticket.Content, ticketdetail, ticketdetail, aiticket.OperatorName)
			err = meet.Send(&robot)
			if err != nil {
				log.Printf("SendTicketOperator failed, err: %v", err)
			}
		} else {
			// 通知给admin角色的用户
			users, err := dao.GetAllUsers()
			if err != nil {
				log.Printf("GetAllUsers failed, err: %v", err)
			}
			for _, user := range users {
				if user.Role == "admin" {
					robot.ToUserName = append(robot.ToUserName, user.UserName)
				}
				ticketdetail := fmt.Sprintf("%s%s?id=%s", global.ServerSetting.Host, "/tickets/detail", ticket.ID)
				robot.Content = fmt.Sprintf("### %s\n%s\n工单详情：[%s](%s)\n请手动指派对应的处理人", ticket.Title, ticket.Content, ticketdetail, ticketdetail)
				err = meet.Send(&robot)
				if err != nil {
					log.Printf("SendTicketOpen failed, err: %v", err)
				}
			}

		}
		dao.UpdateTicketOperatorName(ticket.ID, aiticket.OperatorName)
	}()

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":         ticket.ID,
			"title":      ticket.Title,
			"content":    ticket.Content,
			"status":     ticket.Status,
			"screenshot": ticket.ScreenShots,
			"created_at": ticket.CreatedAt,
			"user_name":  ticket.UserName,
		},
		"msg": "工单创建成功",
	})
}

// GetTicketDetail 获取工单详情
// @Summary 获取工单详情
// @Description 根据工单ID获取详细信息，包含权限验证
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param id path string true "工单ID"
// @Success 200 {object} model.Ticket "工单详情"
// @Failure 401 {object} model.Response{error=string} "未授权"
// @Failure 403 {object} model.Response{error=string} "无权限"
// @Failure 404 {object} model.Response{error=string} "工单不存在"
// @Failure 500 {object} model.Response{error=string} "服务器错误"
// @Router /api/v1/tickets/{id}/detail [get]
func GetTicketDetail(c *gin.Context) {
	ticketID := c.Param("id")

	// 获取并解析JWT令牌
	tokenString, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少访问令牌"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(global.OauthSetting.ClientSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	// 从JWT声明中获取用户信息
	var currentUserName string
	var currentRole string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userName, ok := claims["user_name"].(string); ok {
			currentUserName = userName
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户名类型"})
			return
		}
		if role, ok := claims["role"].(string); ok {
			currentRole = role
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户角色"})
			return
		}
	}

	ticket, err := dao.GetTicketDetail(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if currentUserName != ticket.UserName {
		if currentRole != "admin" && currentRole != "operator" {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限访问此工单"})
			return
		}
	}

	response := gin.H{
		"id":            ticket.ID,
		"user_id":       ticket.UserID,
		"user_name":     ticket.UserName,
		"title":         ticket.Title,
		"ticket_type":   ticket.TicketType,
		"content":       ticket.Content,
		"status":        ticket.Status,
		"remark":        ticket.Remark,
		"solution":      ticket.Solution,
		"created_at":    ticket.CreatedAt,
		"updated_at":    ticket.UpdatedAt,
		"operator_name": ticket.OperatorName,
		"screenshots":   ticket.ScreenShots,
	}

	c.JSON(http.StatusOK, response)
}

// GetTickets 获取工单列表
// @Summary 获取工单列表
// @Description 根据条件查询工单列表，支持分页和筛选
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param title query string false "工单标题模糊搜索"
// @Param status query string false "工单状态(open/closed)"
// @Param ticket_type query string false "工单类型(IT/系统)"
// @Param start_date query string false "开始日期(YYYY-MM-DD)"
// @Param end_date query string false "结束日期(YYYY-MM-DD)"
// @Param user_name query string false "创建人用户名"
// @Param operator_name query string false "处理人用户名"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} model.Response{data=map[string]interface{}} "工单列表"
// @Failure 401 {object} model.Response{error=string} "未授权"
// @Failure 403 {object} model.Response{error=string} "无权限"
// @Failure 500 {object} model.Response{error=string} "服务器错误"
// @Router /api/v1/tickets [get]
func GetTickets(c *gin.Context) {

	// 获取并解析JWT令牌
	tokenString, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少访问令牌"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(global.OauthSetting.ClientSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	// 从JWT声明中获取用户信息
	var currentUserName string
	var currentRole string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userName, ok := claims["user_name"].(string); ok {
			currentUserName = userName
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户名类型"})
			return
		}
		if role, ok := claims["role"].(string); ok {
			currentRole = role
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户角色"})
			return
		}
	}

	title := c.Query("title")
	status := c.Query("status")
	ticketType := c.Query("ticket_type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	userName := c.Query("user_name")
	operatorName := c.Query("operator_name")

	page := c.Query("page")
	limit := c.Query("limit")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	tickets, total, err := dao.GetTickets(title, status, ticketType, startDate, endDate, userName, operatorName, pageInt, limitInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if currentRole == "user" || currentRole == "operator" {
		if userName == "" && operatorName == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限访问此工单"})
			return
		}
	}

	if userName != "" {
		if userName != currentUserName && currentRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限访问此工单"})
			return
		}
	}

	if operatorName != "" {
		if currentRole != "admin" && currentRole != "operator" {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限访问此工单"})
			return
		}
		if currentRole == "operator" {
			if operatorName != currentUserName {
				c.JSON(http.StatusForbidden, gin.H{"error": "没有权限访问此工单"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"tickets": tickets,
		"total":   total,
		"page":    pageInt,
		"limit":   limitInt,
	})
}

// UpdateTicketStatus 更新工单状态
// @Summary 更新工单状态
// @Description 更新指定工单的状态
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param id path string true "工单ID"
// @Param status body model.StatusUpdate true "状态更新信息"
// @Success 200 {object} model.Response{message=string} "更新成功"
// @Failure 400 {object} model.Response{error=string} "请求参数错误"
// @Failure 401 {object} model.Response{error=string} "未授权"
// @Failure 500 {object} model.Response{error=string} "服务器错误"
// @Router /api/v1/tickets/{id}/status [put]
func UpdateTicketStatus(c *gin.Context) {
	ticketID := c.Param("id")
	var robot meet.Robot

	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证状态值
	if statusUpdate.Status != "open" && statusUpdate.Status != "closed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的状态值"})
		return
	}

	// 更新工单状态
	if err := dao.UpdateTicketStatus(ticketID, statusUpdate.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ticket, err := dao.GetTicketDetail(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if ticket.Status == "closed" {
		// 发送关闭工单消息
		robot.ToUserName = []string{ticket.UserName}
		ticketdetail := fmt.Sprintf("%s%s?id=%s", global.ServerSetting.Host, "/tickets/detail", ticket.ID)
		robot.Content = fmt.Sprintf("### %s\n%s\n工单详情：[%s](%s)\n工单状态：<font color=\"red\">已关闭</font>", ticket.Title, ticket.Content, ticketdetail, ticketdetail)
		err = meet.Send(&robot)
		if err != nil {
			log.Printf("SendTicketOpen failed, err: %v", err)
		}

	}

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

// UpdateTicketOperatorName 更新工单的处理人姓名
// @Summary 更新工单处理人
// @Description 更新指定工单的处理人姓名
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param id path string true "工单ID"
// @Param operator body model.TicketOperatorUpdate true "处理人更新信息"
// @Success 200 {object} model.Response
// @Router /api/v1/tickets/{id}/operators [put]
func UpdateTicketOperatorName(c *gin.Context) {
	var robot meet.Robot
	ticketID := c.Param("id")

	// 解析请求体
	var updateData struct {
		OperatorName string `json:"operator_name"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ticket, err := dao.GetTicketDetail(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	robot.ToUserName = strings.Split(updateData.OperatorName, ",")
	ticketdetail := fmt.Sprintf("%s%s?id=%s", global.ServerSetting.Host, "/tickets/detail", ticketID)
	robot.Content = fmt.Sprintf("### %s\n%s\n工单详情：[%s](%s)\n此工单已指派给你", ticket.Title, ticket.Content, ticketdetail, ticketdetail)
	err = meet.Send(&robot)
	if err != nil {
		log.Printf("SendTicketOperator failed, err: %v", err)
	}
	robot.ToUserName = []string{ticket.UserName}
	robot.Content = fmt.Sprintf("### %s\n%s\n工单详情：[%s](%s)\n处理人: **%s**", ticket.Title, ticket.Content, ticketdetail, ticketdetail, updateData.OperatorName)
	err = meet.Send(&robot)
	if err != nil {
		log.Printf("SendTicketOperator failed, err: %v", err)
	}
	// 更新处理人姓名
	if err := dao.UpdateTicketOperatorName(ticketID, updateData.OperatorName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "处理人更新成功"})
}

// UpdateTicketType 更新工单类型
// @Summary 更新工单类型
// @Description 更新指定工单的类型
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param id path string true "工单ID"
// @Param type body model.TicketTypeUpdate true "工单类型更新信息"
// @Success 200 {object} model.Response
// @Router /api/v1/tickets/{id}/type [put]
func UpdateTicketType(c *gin.Context) {
	ticketID := c.Param("id")

	// 解析请求体
	var updateData struct {
		TicketType string `json:"ticket_type"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证类型值
	if updateData.TicketType != "IT" && updateData.TicketType != "系统" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的类型值"})
		return
	}

	// 更新工单类型
	if err := dao.UpdateTicketType(ticketID, updateData.TicketType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "工单类型更新成功"})
}

// UpdateTicketRemark 更新工单备注
// @Summary 更新工单备注
// @Description 更新指定工单的备注信息
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param id path string true "工单ID"
// @Param remark body model.RemarkUpdate true "备注更新信息"
// @Success 200 {object} model.Response{message=string} "更新成功"
// @Failure 400 {object} model.Response{error=string} "请求参数错误"
// @Failure 401 {object} model.Response{error=string} "未授权"
// @Failure 500 {object} model.Response{error=string} "服务器错误"
// @Router /api/v1/tickets/{id}/remark [put]
func UpdateTicketRemark(c *gin.Context) {
	ticketID := c.Param("id")

	// 解析请求体
	var updateData struct {
		Remark string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新工单类型
	if err := dao.UpdateTicketRemark(ticketID, updateData.Remark); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "工单备注更新成功"})
}

// UpdateTicketSolution 更新工单解决方案
// @Summary 更新工单解决方案
// @Description 更新指定工单的解决方案
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param id path string true "工单ID"
// @Param solution body model.SolutionUpdate true "解决方案更新信息"
// @Success 200 {object} model.Response{message=string} "更新成功"
// @Failure 400 {object} model.Response{error=string} "请求参数错误"
// @Failure 401 {object} model.Response{error=string} "未授权"
// @Failure 500 {object} model.Response{error=string} "服务器错误"
// @Router /api/v1/tickets/{id}/solution [put]
func UpdateTicketSolution(c *gin.Context) {
	ticketID := c.Param("id")

	// 解析请求体
	var updateData struct {
		Solution string `json:"solution"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新工单类型
	if err := dao.UpdateTicketSolution(ticketID, updateData.Solution); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "工单解决方案更新成功"})
}
