package routers

import (
	"lyworder/global"
	"lyworder/internal/middleware"
	v1 "lyworder/internal/routers/api/v1"

	"github.com/gin-gonic/gin"
	
	// 添加Swagger支持
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	// 静态文件
	r.StaticFile("/ticket.png", "./static/ticket.png")
	r.Static("/uploads", "./static/uploads")
	r.Static("/vendor", "./static/vendor")

	// 中间件
	r.Use(middleware.CORSMiddleware())
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	}

	// 页面路由
	pageRoutes := []struct {
		path       string
		handler    gin.HandlerFunc
		middleware []gin.HandlerFunc
	}{
		{"/", fileHandler("./static/index.html"), []gin.HandlerFunc{middleware.AuthMiddleware()}},
		{"/tickets", fileHandler("./static/my_tickets.html"), []gin.HandlerFunc{middleware.AuthMiddleware()}},
		{"/admin", fileHandler("./static/admin.html"), []gin.HandlerFunc{middleware.RoleCheckMiddleware("admin")}},
		{"/oamuser", fileHandler("./static/oam_management.html"), []gin.HandlerFunc{middleware.RoleCheckMiddleware("admin")}},
		{"/operator", fileHandler("./static/my_operator.html"), []gin.HandlerFunc{middleware.RoleCheckMiddleware("admin", "operator")}},
		{"/tickets/detail", fileHandler("./static/ticket_detail.html"), []gin.HandlerFunc{middleware.AuthMiddleware()}},
	}

	for _, route := range pageRoutes {
		handlers := append(route.middleware, route.handler)
		r.GET(route.path, handlers...)
	}

	// 认证路由
	authRoutes := []struct {
		path    string
		handler gin.HandlerFunc
	}{
		{"/login", middleware.LoginHandler},
		{"/oauth/callback", middleware.CallbackHandler},
		{"/logout", middleware.LogoutHandler},
	}

	for _, route := range authRoutes {
		r.GET(route.path, route.handler)
	}

	// Swagger文档路由（无需认证）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/users", middleware.HandleUserInfo)
		api.GET("/operator/usersname", v1.GetOperatorUsersName)
		// 所有用户表操作都通过/operator/users路由处理
		api.GET("/operator/users", v1.GetAllUsers)
		api.POST("/operator/users", v1.CreateUser)
		api.DELETE("/operator/users", v1.DeleteUser)
		api.GET("/operator/users/status", v1.GetUsersByStatus)
		api.PUT("/operator/users/status", v1.UpdateUserStatus)
		api.PUT("/operator/users/type", v1.UpdateUserType)
		api.GET("/tickets", v1.GetTickets)
		api.GET("/tickets/:id/detail", v1.GetTicketDetail)
		api.POST("/tickets", v1.CreateTicket)
		api.PUT("/tickets/:id/status", v1.UpdateTicketStatus)
		api.PUT("/tickets/:id/operatorname", v1.UpdateTicketOperatorName)
		api.PUT("/tickets/:id/remark", v1.UpdateTicketRemark)
		api.PUT("/tickets/:id/solution", v1.UpdateTicketSolution)
		api.PUT("/tickets/:id/type", v1.UpdateTicketType)
	}

	return r
}

// 辅助函数：返回文件处理handler
func fileHandler(filePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.File(filePath)
	}
}
