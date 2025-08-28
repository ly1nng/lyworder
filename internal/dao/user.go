package dao

import (
	"lyworder/global"
	"lyworder/internal/model"

	"gorm.io/gorm"
)

// 新增方法：获取所有用户名
func GetAllUsernames() ([]string, error) {
	var usernames []string
	err := global.DBEngine.Model(&model.User{}).
		Pluck("user_name", &usernames).Error
	if err != nil {
		return nil, err
	}
	return usernames, nil
}

// 获取所有用户
func GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := global.DBEngine.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 根据状态获取用户
func GetUsersByStatus(status string) ([]model.User, error) {
	var users []model.User
	err := global.DBEngine.Where("status = ?", status).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 更新用户类型
func UpdateUserType(userName, opsType string) error {
	return global.DBEngine.Model(&model.User{}).
		Where("user_name = ?", userName).
		Update("ops_type", opsType).Error
}

// 更新用户状态
func UpdateUserStatus(userName, status string) error {
	return global.DBEngine.Model(&model.User{}).
		Where("user_name = ?", userName).
		Update("status", status).Error
}

func GetUserByUserName(userName string) (*model.User, error) {
	var user model.User
	err := global.DBEngine.Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetUsersByOpsType(ops_type string) ([]model.User, error) {
	var users []model.User
	err := global.DBEngine.Where("ops_type = ?", ops_type).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func CreateUser(user *model.User) error {
	return global.DBEngine.Create(user).Error
}

func UpdateUserRole(userName, role string) error {
	return global.DBEngine.Model(&model.User{}).
		Where("user_name = ?", userName).
		Update("role", role).Error
}

// 在文件末尾添加删除用户函数
func DeleteUser(userName string) error {
	return global.DBEngine.Where("user_name = ?", userName).Delete(&model.User{}).Error
}
