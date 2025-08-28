package model

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// StatusUpdate 状态更新请求结构
type StatusUpdate struct {
	Status string `json:"status" binding:"required" example:"open"`
}

// RemarkUpdate 备注更新请求结构
type RemarkUpdate struct {
	Remark string `json:"remark" binding:"required" example:"这是一个备注"`
}

// SolutionUpdate 解决方案更新请求结构
type SolutionUpdate struct {
	Solution string `json:"solution" binding:"required" example:"这是解决方案"`
}

// TicketTypeUpdate 工单类型更新请求结构
type TicketTypeUpdate struct {
	TicketType string `json:"ticket_type" binding:"required" example:"IT"`
}

// TicketOperatorUpdate 处理人更新请求结构
type TicketOperatorUpdate struct {
	OperatorName string `json:"operator_name" binding:"required" example:"admin"`
}

// UserCreate 用户创建请求结构
type UserCreate struct {
	UserName string `json:"user_name" binding:"required" example:"newuser"`
	Role     string `json:"role" binding:"required" example:"operator"`
	Status   int    `json:"status" binding:"required" example:"1"`
	OpsType  string `json:"ops_type" example:"system"`
}
