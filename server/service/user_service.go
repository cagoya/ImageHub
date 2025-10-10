package service

import (
	"fmt"
	"math"
	"os"
	"path"
	"strconv"

	"server/model"
	"server/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// RegisterUser 注册用户
func (s *UserService) RegisterUser(registerRequest *model.RegisterRequest) (user *model.User, err error) {
	// 检查用户名是否已存在
	if err := s.db.Where("username = ?", registerRequest.Username).First(&model.User{}).Error; err == nil {
		return nil, fmt.Errorf("用户名已存在")
	}
	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", registerRequest.Email).First(&model.User{}).Error; err == nil {
		return nil, fmt.Errorf("邮箱已存在")
	}
	// 检查密码长度
	if len(registerRequest.Password) < 6 {
		return nil, fmt.Errorf("密码长度不能小于6位")
	}
	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("生成密码哈希失败: %v", err)
	}
	// 创建用户
	user = &model.User{
		Username:     registerRequest.Username,
		Email:        registerRequest.Email,
		PasswordHash: string(hashedPassword),
		IsAdmin:      false,
	}
	err = s.db.Create(&user).Error
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}
	// 创建目录
	os.MkdirAll(path.Join("./public", strconv.Itoa(user.ID), "origin"), 0755)
	os.MkdirAll(path.Join("./public", strconv.Itoa(user.ID), "thumbnail"), 0755)
	return user, nil
}

// LoginUser 登录用户
func (s *UserService) LoginUser(loginRequest *model.LoginRequest) (*model.User, string, error) {
	// 检查用户名是否存在
	if err := s.db.Where("username = ?", loginRequest.Username).First(&model.User{}).Error; err != nil {
		return nil, "", fmt.Errorf("用户名不存在")
	}
	// 检查密码是否正确
	user, err := s.GetUserByUsername(loginRequest.Username)
	if err != nil {
		return nil, "", fmt.Errorf("查询用户失败: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)); err != nil {
		return nil, "", fmt.Errorf("密码错误")
	}
	// 生成token
	token, err := utils.CreateToken(user.ID, user.Username, user.Email, user.IsAdmin)
	if err != nil {
		return nil, "", fmt.Errorf("生成token失败: %v", err)
	}
	return user, token, nil
}

// ForgetPassword 忘记密码
func (s *UserService) ForgetPassword(forgetPasswordRequest *model.ForgetPasswordRequest) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", forgetPasswordRequest.Username).First(&user).Error; err != nil {
		return nil, err
	}
	if user.Email != forgetPasswordRequest.Email {
		return nil, fmt.Errorf("邮箱不匹配")
	}
	if len(forgetPasswordRequest.NewPassword) < 6 {
		return nil, fmt.Errorf("密码长度不能小于6位")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(forgetPasswordRequest.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("生成密码哈希失败: %v", err)
	}
	user.PasswordHash = string(hashedPassword)
	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("更新密码失败: %v", err)
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id int) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers 获取所有用户，分页查询
func (s *UserService) GetAllUsers(page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	if err := s.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	var count int64
	s.db.Model(&model.User{}).Count(&count)
	pageTotal := int64(math.Ceil(float64(count) / float64(pageSize)))
	return users, pageTotal, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(updateUserRequest *model.UpdateUserRequest) error {
	// 检查用户名是否已存在
	if err := s.db.Where("username = ?", updateUserRequest.Username).First(&model.User{}).Error; err == nil {
		return fmt.Errorf("用户名已存在")
	}
	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", updateUserRequest.Email).First(&model.User{}).Error; err == nil {
		return fmt.Errorf("邮箱已存在")
	}
	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUserRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成密码哈希失败: %v", err)
	}
	user, err := s.GetUserByID(updateUserRequest.ID)
	if err != nil {
		return fmt.Errorf("获取用户失败: %v", err)
	}
	user.Username = updateUserRequest.Username
	user.Email = updateUserRequest.Email
	user.PasswordHash = string(hashedPassword)
	return s.db.Save(user).Error
}

// GrantAdminPermission 给予管理员权限
func (s *UserService) GrantAdminPermission(id int) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("is_admin", true).Error
}

// RevokeAdminPermission 撤销管理员权限
func (s *UserService) RevokeAdminPermission(id int) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("is_admin", false).Error
}

// IsAdmin 判断用户是否为管理员
func (s *UserService) IsAdmin(id int) (bool, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int) error {
	return s.db.Delete(&model.User{}, id).Error
}
