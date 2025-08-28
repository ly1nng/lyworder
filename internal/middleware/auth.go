package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"lyworder/global"
	"lyworder/internal/dao"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// 统一响应格式
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func LoginHandler(c *gin.Context) {
	// 获取原始访问路径，优先从state参数获取
	redirectURI := c.Query("state")
	if redirectURI == "" {
		redirectURI = c.Query("redirect_uri")
	}
	if redirectURI == "" {
		redirectURI = "/"
	}

	// 构造OAuth2认证地址
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&state=%s",
		global.OauthSetting.AuthURL,
		global.OauthSetting.ClientID,
		url.QueryEscape(global.OauthSetting.RedirectURI),
		url.QueryEscape(redirectURI),
	)

	c.Redirect(http.StatusFound, authURL)
}

func CallbackHandler(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少授权码"})
		return
	}

	params := url.Values{}
	params.Add("client_id", global.OauthSetting.ClientID)
	params.Add("client_secret", global.OauthSetting.ClientSecret)
	params.Add("redirect_uri", global.OauthSetting.RedirectURI)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")

	// 拼接完整的请求 URL
	tokenURL := fmt.Sprintf("%s?%s", global.OauthSetting.TokenURL, params.Encode())

	// 发起 GET 请求
	resp, err := http.Get(tokenURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌请求失败"})
		return
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	var tokenRes struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌解析失败"})
		return
	}

	// 获取用户信息
	userInfo, err := GetUserInfo(tokenRes.AccessToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 生成JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userInfo.UserID,
		"user_name": userInfo.UserName,
		"exp":       time.Now().Add(2 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(global.OauthSetting.ClientSecret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 直接跳转到首页
	redirectPath := "/"

	// 设置Cookie后跳转到首页
	c.SetCookie("access_token", tokenString, 3600, "/", global.ServerSetting.Domain, false, true)
	c.Redirect(http.StatusFound, redirectPath)
}

func GetUserInfo(accessToken string) (struct {
	UserID   int    `json:"id"`
	UserName string `json:"user_name"`
	Mobile   string `json:"mobile"`
}, error) {
	var userInfo struct {
		UserID   int    `json:"id"`
		UserName string `json:"user_name"`
		Mobile   string `json:"mobile"`
	}

	req, _ := http.NewRequest("POST", global.OauthSetting.UserInfoURL, nil)
	req.Header.Set("accessToken", accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return userInfo, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return userInfo, err
	}

	return userInfo, nil
}

// 新增角色检查中间件
func RoleCheckMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取访问令牌
		tokenString, err := c.Cookie("access_token")
		if err != nil || tokenString == "" {
			handleUnauthorized(c, "缺少访问令牌")
			return
		}

		// 解析JWT获取角色信息
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(global.OauthSetting.ClientSecret), nil
		})

		if err != nil || !token.Valid {
			handleUnauthorized(c, "无效的访问令牌")
			return
		}

		// 从JWT中获取角色
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handleUnauthorized(c, "无效的令牌声明")
			return
		}

		// 检查用户角色是否在允许的角色列表中
		userRole, ok := claims["role"].(string)
		if !ok {
			handleUnauthorized(c, "无效的用户角色")
			return
		}

		// 检查用户角色是否匹配任一要求的角色
		for _, role := range requiredRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		handleUnauthorized(c, "权限不足")

	}

}

// 在AuthMiddleware中添加角色声明
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取访问令牌
		tokenString, err := c.Cookie("access_token")
		if err != nil || tokenString == "" {
			handleUnauthorized(c, "缺少访问令牌")
			return
		}

		// AuthMiddleware 中验证令牌
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(global.OauthSetting.ClientSecret), nil
		})

		if err != nil || !token.Valid {
			handleUnauthorized(c, "无效的访问令牌")
			return
		}

		c.Next()
	}
}
func LogoutHandler(c *gin.Context) {
	// 1. 获取当前页面路径
	referer := c.Request.Referer()
	if referer == "" {
		referer = "/"
	}

	// 2. 清除本地Cookie
	c.SetCookie("access_token", "", -1, "/", global.ServerSetting.Domain, false, true)

	// 3. 构造OAuth服务端登出URL
	logoutURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s",
		global.OauthSetting.LogoutURL,
		global.OauthSetting.ClientID,
		url.QueryEscape(referer),
	)

	// 4. 重定向到OAuth登出页面
	c.Redirect(http.StatusFound, logoutURL)
}

func HandleUserInfo(c *gin.Context) {
	// 从Cookie获取JWT
	tokenString, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    http.StatusUnauthorized,
			Message: "无法获取用户信息",
		})
		return
	}

	// 解析JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(global.OauthSetting.ClientSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    http.StatusUnauthorized,
			Message: "无效的访问令牌",
		})
		return
	}

	// 从JWT中提取用户信息
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    http.StatusUnauthorized,
			Message: "无效的令牌声明",
		})
		return
	}

	// 查询数据库获取用户角色
	userName, ok := claims["user_name"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    http.StatusUnauthorized,
			Message: "无效的用户名",
		})
		return
	}

	user, err := dao.GetUserByUserName(userName)
	var role string = "user" // 默认角色
	if err == nil && user != nil {
		role = user.Role
	}

	// 更新token中的角色信息
	claims["role"] = role
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := newToken.SignedString([]byte(global.OauthSetting.ClientSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "无法更新令牌",
		})
		return
	}

	// 设置新的cookie
	c.SetCookie("access_token", newTokenString, 3600*24, "/", global.ServerSetting.Domain, false, true)

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "成功",
		Data: gin.H{
			"user_id":   claims["user_id"],
			"user_name": claims["user_name"],
			"role":      role,
		},
	})
}

// 统一处理未授权情况
func handleUnauthorized(c *gin.Context, message string) {
	// 根据请求类型返回不同响应
	if isAPIRequest(c) {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    http.StatusUnauthorized,
			Message: message,
		})
	} else {
		originalPath := url.QueryEscape(c.Request.URL.String())
		c.Redirect(http.StatusFound, "/login?redirect_uri="+originalPath)
	}
	c.Abort()
}

// 判断是否为API请求
func isAPIRequest(c *gin.Context) bool {
	return strings.HasPrefix(c.Request.URL.Path, "/api")
}
