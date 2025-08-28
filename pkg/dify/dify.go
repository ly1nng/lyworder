package dify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lyworder/global"
	"lyworder/internal/dao"
	"lyworder/internal/model"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

// 本地缓存
var (
	localCache = make(map[string]int64)
	cacheMutex sync.Mutex
)

// 调度算法函数
func ScheduleTicket(ticket *model.Ticket) (*model.Ticket, error) {
	// 检查必要的配置
	if global.DifySetting == nil {
		return ticket, fmt.Errorf("DifySetting未初始化")
	}
	if global.DifySetting.APIURL == "" || global.DifySetting.APIKey == "" {
		return ticket, fmt.Errorf("请配置Dify API地址和密钥")
	}

	// 1. 调用Dify API判断工单类型
	reqBody := map[string]interface{}{
		"inputs":          map[string]interface{}{},
		"query":           ticket.Content,
		"response_mode":   "blocking",
		"conversation_id": "",
		"user":            ticket.UserName,
		"files":           []interface{}{},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", global.DifySetting.APIURL+"/chat-messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+global.DifySetting.APIKey)

	// 设置超时
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 发送请求
	difyResp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("调用Dify API失败: %v", err)
	}
	defer difyResp.Body.Close()

	// 检查响应状态
	if difyResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Dify API返回非成功状态码: %d", difyResp.StatusCode)
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(difyResp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析Dify响应失败: %v", err)
	}

	answerStr, ok := result["answer"].(string)
	if !ok {
		return nil, fmt.Errorf("获取answer字段失败")
	}

	// 处理 JSON 代码块
	if jsonStart := strings.Index(answerStr, "```json"); jsonStart != -1 {
		// 移除开始标记及之前的内容
		answerStr = answerStr[jsonStart+len("```json"):]

		// 查找结束标记位置（支持带空格/换行的变体）
		jsonEnd := strings.Index(answerStr, "```")
		if jsonEnd == -1 {
			// 尝试带换行的结束标记
			if jsonEnd = strings.Index(answerStr, "```\n"); jsonEnd == -1 {
				jsonEnd = strings.Index(answerStr, "``` ")
			}
		}

		if jsonEnd != -1 {
			answerStr = answerStr[:jsonEnd]
		}
		answerStr = strings.TrimSpace(answerStr)
	} else {
		// 直接提取 JSON 对象
		jsonStart := strings.Index(answerStr, "{")
		jsonEnd := strings.LastIndex(answerStr, "}")
		if jsonStart != -1 && jsonEnd != -1 && jsonStart < jsonEnd {
			answerStr = answerStr[jsonStart : jsonEnd+1]
		} else {
			return nil, fmt.Errorf("answer字段不包含有效的JSON")
		}
	}

	// 解析内部JSON
	var answerJSON map[string]interface{}
	if err := json.Unmarshal([]byte(answerStr), &answerJSON); err != nil {
		return nil, fmt.Errorf("解析answer JSON失败: %v", err)
	}

	// 提取工单类型
	ticketType, ok := answerJSON["category"].(string)
	if !ok {
		return nil, fmt.Errorf("获取工单类型失败")
	}
	ticketRemark, ok := answerJSON["reason"].(string)
	if !ok {
		return nil, fmt.Errorf("获取工单备注失败")
	}

	// 验证工单类型
	ticketType = strings.TrimSpace(ticketType)
	if ticketType != "IT" && ticketType != "系统" {
		return nil, fmt.Errorf("无效的工单类型: %s", ticketType)
	}

	// 更新工单类型
	ticket.TicketType = ticketType
	ticket.Remark = ticketRemark

	// 2. 查找最佳处理人
	// 2.1 首先查找状态正常且ops_type匹配的用户
	switch ticketType {
	case "IT":
		activeUsers, err := dao.GetUsersByOpsType("it")
		if err != nil {
			return nil, fmt.Errorf("查询处理人失败: %v", err)
		}
		operateTime, _ := answerJSON["operate_time"].(float64)
		ticket.OperatorName = selectOperator(activeUsers, int(operateTime))

	case "系统":
		activeUsers, err := dao.GetUsersByOpsType("system")
		if err != nil {
			return nil, fmt.Errorf("查询处理人失败: %v", err)
		}
		operateTime, _ := answerJSON["operate_time"].(float64)
		ticket.OperatorName = selectOperator(activeUsers, int(operateTime))
	}

	return ticket, nil
}

func selectOperator(users []model.User, operateTime int) string {
	// 过滤状态为1的活跃用户
	activeUsers := []model.User{}
	for _, user := range users {
		if user.Status == 1 {
			activeUsers = append(activeUsers, user)
		}
	}

	// 没有活跃用户，返回空
	if len(activeUsers) == 0 {
		return ""
	}

	// 只有一个活跃用户，直接返回
	if len(activeUsers) == 1 {
		return activeUsers[0].UserName
	}

	// 基于负载的加权随机选择
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	type userLoad struct {
		user   model.User
		load   int64
		weight float64
	}

	var totalWeight float64
	var userLoads []userLoad

	// 计算每个用户的当前负载和权重
	for _, user := range activeUsers {
		loadKey := "user_load_" + user.UserName
		currentLoad, exists := localCache[loadKey]
		if !exists {
			currentLoad = 0
		}

		// 权重与当前负载+新工单处理时间成反比
		weight := 1.0 / float64(currentLoad+int64(operateTime))
		// 如果用户是管理员，降低权重
		if user.Role == "admin" {
			weight *= .5
		}
		userLoads = append(userLoads, userLoad{
			user:   user,
			load:   currentLoad,
			weight: weight,
		})
		totalWeight += weight
	}

	// 加权随机选择
	randValue := rand.Float64() * totalWeight
	var selectedUser model.User
	for _, ul := range userLoads {
		randValue -= ul.weight
		if randValue <= 0 {
			selectedUser = ul.user
			// 更新用户负载
			localCache["user_load_"+selectedUser.UserName] = ul.load + int64(operateTime)
			break
		}
	}

	return selectedUser.UserName
}
