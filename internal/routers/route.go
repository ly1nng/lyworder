package routers

import (
	"lyworder/global"
	"lyworder/internal/middleware"
	v1 "lyworder/internal/routers/api/v1"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	// 添加Swagger支持
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// spaHandler 实现Vue Router History模式支持
func spaHandler(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求的路径是否是API路径或静态资源路径
		requestPath := c.Request.URL.Path

		// 如果是API路径，直接返回404，让其他路由处理
		if strings.HasPrefix(requestPath, "/api/") {
			c.Next()
			return
		}

		// 如果是认证相关路径，直接返回404，让其他路由处理
		if requestPath == "/login" || requestPath == "/oauth/callback" || requestPath == "/logout" {
			c.Next()
			return
		}

		// 如果是静态资源路径，直接返回404，让其他路由处理
		if requestPath == "/ticket.png" || strings.HasPrefix(requestPath, "/uploads/") || strings.HasPrefix(requestPath, "/vendor/") {
			c.Next()
			return
		}

		// 如果是assets路径，检查文件是否存在
		if strings.HasPrefix(requestPath, "/assets/") {
			// 构造文件路径
			filePath := "./static" + requestPath
			// 检查文件是否存在
			if _, err := os.Stat(filePath); err == nil {
				// 文件存在，让其他路由处理（静态文件服务）
				c.Next()
				return
			}
			// 文件不存在，继续处理
		}

		// 对于所有其他路径，返回index.html
		c.File(path)
	}
}

func NewRouter() *gin.Engine {
	r := gin.New()

	// 静态文件
	r.StaticFile("/ticket.png", "./static/ticket.png")
	r.Static("/uploads", "./static/uploads")
	// 添加对assets目录的静态文件支持
	r.Static("/assets", "./static/assets")

	// 中间件
	r.Use(middleware.CORSMiddleware())
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
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
		api.PUT("/operator/users/info", v1.UpdateUserInfo)
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

	// SPA路由 - 处理Vue Router History模式
	// 这些路由需要放在最后，以确保API路由优先匹配
	r.GET("/", middleware.AuthMiddleware(), spaHandler("./static/index.html"))
	r.GET("/tickets", middleware.AuthMiddleware(), spaHandler("./static/index.html"))
	r.GET("/admin", middleware.RoleCheckMiddleware("admin"), spaHandler("./static/index.html"))
	r.GET("/oamuser", middleware.RoleCheckMiddleware("admin"), spaHandler("./static/index.html"))
	r.GET("/operator", middleware.RoleCheckMiddleware("admin", "operator"), spaHandler("./static/index.html"))
	r.GET("/tickets/detail", middleware.AuthMiddleware(), spaHandler("./static/index.html"))

	// 添加一个通配符路由来处理所有其他前端路由
	r.NoRoute(spaHandler("./static/index.html"))

	return r
}
