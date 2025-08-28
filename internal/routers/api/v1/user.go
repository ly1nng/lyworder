// 在包声明前添加
// Package v1 用户管理相关API
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"lyworder/internal/dao"
	"lyworder/internal/model"
)

// GetOperatorUsersName 获取操作员用户名列表
// @Summary 获取操作员用户名列表
// @Description 获取所有操作员用户的用户名列表
// @Tags 用户管理
// @Produce json
// @Success 200 {object} model.Response{data=[]string} "用户名列表"
// @Failure 401 {object} model.Response{error=string} "未授权"
// @Failure 500 {object} model.Response{error=string} "服务器错误"
// @Router /api/v1/operator/usersname [get]
func GetOperatorUsersName(c *gin.Context) {
	users, err := dao.GetAllUsernames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// GetAllUsers 获取所有用户
// @Summary 获取所有用户
// @Description 从users表获取所有用户信息，需要管理员权限
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} model.Response{data=map[string]interface{}} "用户列表"
// @Failure 401 {object} model.Response{error=string} "未授权"
// @Failure 500 {object} model.Response{error=string} "服务器错误"
// @Router /api/v1/operator/users [get]
func GetAllUsers(c *gin.Context) {
	users, err := dao.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Description 更新指定用户的工作状态，需要管理员权限
// @Tags 用户管理
// @Produce json
// @Param user_name query string true "用户名"
// @Param status query int true "状态值(1:工作中, 0:非工作中)"
// @Success 200 {object} model.Response{message=string} "更新成功"
// @Failure 400 {object} model.Response{error=string} "请求参数错误"
// @Failure 401 {object} model.Response{error=string} "未授权"
// @Failure 500 {object} model.Response{error=string} "服务器错误"
// @Router /api/v1/operator/users/status [put]
func UpdateUserStatus(c *gin.Context) {
	userName := c.Query("user_name")
	status := c.Query("status")

	if userName == "" || status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户名和状态不能为空",
		})
		return
	}

	err := dao.UpdateUserStatus(userName, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
	})
}

// UpdateUserType 更新用户运维类型
// @Summary 更新用户运维类型
// @Description 更新指定用户的运维类型，需要管理员权限
// @Tags 用户管理
// @Produce json
// @Param user_name query string true "用户名"
// @Param ops_type query string true "运维类型(system/it)"
// @Success 200 {object} model.Response{message=string} "更新成功"
// @Failure 400 {object} model.Response{error=string} "请求参数错误"
// @Failure 401 {object} model.Response{error=string} "未授权"
// @Failure 500 {object} model.Response{error=string} "服务器错误"
// @Router /api/v1/operator/users/type [put]
func UpdateUserType(c *gin.Context) {
	userName := c.Query("user_name")
	opsType := c.Query("ops_type")

	if userName == "" || opsType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户名和运维类型不能为空",
		})
		return
	}

	err := dao.UpdateUserType(userName, opsType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
	})
}

// GetUsersByStatus 根据状态获取用户
// @Summary 根据状态获取用户
// @Description 从users表获取指定状态的用户，需要管理员权限
// @Tags 用户管理
// @Produce json
// @Param status query int true "状态值(1:工作中, 0:非工作中)"
// @Success 200 {object} model.Response{data=[]model.User} "用户列表"
// @Failure 400 {object} model.Response{error=string} "请求参数错误"
// @Failure 401 {object} model.Response{error=string} "未授权"
// @Failure 500 {object} model.Response{error=string} "服务器错误"
// @Router /api/v1/operator/users/status [get]
func GetUsersByStatus(c *gin.Context) {
	status := c.Query("status")

	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "状态不能为空",
		})
		return
	}

	users, err := dao.GetUsersByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// CreateUser 新增用户
// @Summary 新增用户
// @Description 向users表添加新用户，需要管理员权限
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body model.UserCreate true "用户信息"
// @Success 200 {object} model.Response{message=string} "创建成功"
// @Failure 400 {object} model.Response{error=string} "请求参数错误"
// @Failure 401 {object} model.Response{error=string} "未授权"
// @Failure 500 {object} model.Response{error=string} "服务器错误"
// @Router /api/v1/operator/users [post]
func CreateUser(c *gin.Context) {
	var user model.User

	// 解析请求体
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求体格式错误: " + err.Error(),
		})
		return
	}

	// 检查必填字段
	if user.UserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户名不能为空",
		})
		return
	}

	// 调用DAO层函数添加用户
	if err := dao.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "添加用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "添加用户成功",
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 根据用户名删除用户，需要管理员权限
// @Tags 用户管理
// @Produce json
// @Param user_name query string true "用户名"
// @Success 200 {object} model.Response{message=string} "删除成功"
// @Failure 400 {object} model.Response{error=string} "请求参数错误"
// @Failure 401 {object} model.Response{error=string} "未授权"
// @Failure 500 {object} model.Response{error=string} "服务器错误"
// @Router /api/v1/operator/users [delete]
func DeleteUser(c *gin.Context) {
	userName := c.Query("user_name")

	if userName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户名不能为空",
		})
		return
	}

	// 检查用户是否存在
	existingUser, err := dao.GetUserByUserName(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询用户失败: " + err.Error(),
		})
		return
	}

	if existingUser == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 执行删除操作
	err = dao.DeleteUser(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除用户成功",
	})
}
