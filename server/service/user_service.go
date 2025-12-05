package service

import (
	"context"
	"fmt"
	"math"
	"net/mail"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"

	"server/model"
	"server/utils"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewUserService(db *gorm.DB, redisClient *redis.Client) *UserService {
	return &UserService{db: db, redisClient: redisClient}
}

// RegisterUser 注册用户
func (s *UserService) RegisterUser(registerRequest *model.RegisterRequest) (user *model.User, err error) {
	// 检查用户名是否已存在
	if err := s.db.Where("username = ?", registerRequest.Username).First(&model.User{}).Error; err == nil {
		return nil, fmt.Errorf("用户名已存在")
	}
	// 检查邮箱格式
	if _, err := mail.ParseAddress(registerRequest.Email); err != nil {
		return nil, fmt.Errorf("邮箱格式不正确")
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
	// 检查用户是否被封禁（数据库字段）
	if user.IsSuspended {
		// 检查封禁是否已到期
		if user.SuspendedUntil != nil && time.Now().After(*user.SuspendedUntil) {
			// 封禁已到期，自动解封
			err = s.UnsuspendUser(user.ID)
			if err != nil {
				return nil, "", fmt.Errorf("自动解封失败: %v", err)
			}
			// 重新获取用户信息
			user, err = s.GetUserByUsername(loginRequest.Username)
			if err != nil {
				return nil, "", fmt.Errorf("查询用户失败: %v", err)
			}
		} else {
			// 封禁未到期（SuspendedUntil为nil表示永久封禁，或时间未到）
			return nil, "", fmt.Errorf("当前用户已被封禁")
		}
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
	// 检查邮箱格式
	if _, err := mail.ParseAddress(forgetPasswordRequest.Email); err != nil {
		return nil, fmt.Errorf("邮箱格式不正确")
	}
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
	// 先获取当前用户信息
	user, err := s.GetUserByID(updateUserRequest.ID)
	if err != nil {
		return fmt.Errorf("获取用户失败: %v", err)
	}

	// 检查用户名是否已存在（排除当前用户）
	var existingUser model.User
	if err := s.db.Where("username = ? AND id != ?", updateUserRequest.Username, updateUserRequest.ID).First(&existingUser).Error; err == nil {
		return fmt.Errorf("用户名已存在")
	}

	// 检查邮箱格式
	if _, err := mail.ParseAddress(updateUserRequest.Email); err != nil {
		return fmt.Errorf("邮箱格式不正确")
	}

	// 检查邮箱是否已存在（排除当前用户）
	if err := s.db.Where("email = ? AND id != ?", updateUserRequest.Email, updateUserRequest.ID).First(&existingUser).Error; err == nil {
		return fmt.Errorf("邮箱已存在")
	}

	// 更新用户信息
	user.Username = updateUserRequest.Username
	user.Email = updateUserRequest.Email

	// 如果提供了新密码，则更新密码
	if updateUserRequest.Password != "" {
		if len(updateUserRequest.Password) < 6 {
			return fmt.Errorf("密码长度不能小于6位")
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUserRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("生成密码哈希失败: %v", err)
		}
		user.PasswordHash = string(hashedPassword)
	}

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

// SuspendUser 封禁用户
func (s *UserService) SuspendUser(id int, expiration time.Duration) error {
	// 计算封禁到期时间
	suspendedUntil := time.Now().Add(expiration)

	// 将用户ID加到Redis中
	err := s.redisClient.Set(context.Background(), fmt.Sprintf("user%d_isSuspended", id), true, expiration).Err()
	if err != nil {
		return fmt.Errorf("添加被封禁用户ID到Redis失败: %v", err)
	}
	// 同时更新数据库字段
	err = s.db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_suspended":    true,
		"suspended_until": suspendedUntil,
	}).Error
	if err != nil {
		return fmt.Errorf("更新用户封禁状态到数据库失败: %v", err)
	}
	return nil
}

// UnsuspendUser 解封用户
func (s *UserService) UnsuspendUser(id int) error {
	// 从Redis中删除封禁记录
	err := s.redisClient.Del(context.Background(), fmt.Sprintf("user%d_isSuspended", id)).Err()
	if err != nil {
		return fmt.Errorf("删除封禁用户ID从Redis失败: %v", err)
	}
	// 同时更新数据库字段
	err = s.db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_suspended":    false,
		"suspended_until": nil,
	}).Error
	if err != nil {
		return fmt.Errorf("更新用户解封状态到数据库失败: %v", err)
	}
	return nil
}

// 设置是否解析 EXIF
func (s *UserService) SetParseExif(id int, parseExif bool) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("parse_exif", parseExif).Error
}

// StartSuspendExpirationListener 启动监听Redis键过期事件，自动解封用户
// 注意：需要在Redis配置中启用键空间通知，设置 notify-keyspace-events Ex
func (s *UserService) StartSuspendExpirationListener(ctx context.Context) {
	// 获取Redis数据库编号
	dbNum := s.redisClient.Options().DB
	// 订阅键过期事件频道
	pubsub := s.redisClient.PSubscribe(ctx, fmt.Sprintf("__keyevent@%d__:expired", dbNum))
	defer pubsub.Close()

	// 用于匹配封禁用户键的正则表达式
	keyPattern := regexp.MustCompile(`^user(\d+)_isSuspended$`)

	// 处理消息
	ch := pubsub.Channel()
	for msg := range ch {
		// 解析键名
		key := msg.Payload
		matches := keyPattern.FindStringSubmatch(key)
		if len(matches) == 2 {
			// 提取用户ID
			userID, err := strconv.Atoi(matches[1])
			if err != nil {
				// 如果解析失败，跳过这条消息
				continue
			}
			// 更新数据库，将用户解封
			err = s.db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
				"is_suspended":    false,
				"suspended_until": nil,
			}).Error
			if err != nil {
				// 记录错误，但不中断监听
				fmt.Printf("自动解封用户失败 (用户ID: %d): %v\n", userID, err)
			} else {
				fmt.Printf("用户 %d 的封禁时间已到期，已自动解封\n", userID)
			}
		}
	}
}
