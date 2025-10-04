package test

import (
	"fmt"
	"testing"

	"server/model"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// 测试注册用户
func TestRegisterUser(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "李华",
		Email:    "lihua@zju.edu.cn",
		Password: "123456",
	}
	user, err := userService.RegisterUser(registerRequest)
	assert.NoError(t, err)
	err = db.Where("username = ?", "李华").First(&user).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Username, "李华")
	assert.Equal(t, user.Email, "lihua@zju.edu.cn")
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(registerRequest.Password))
	assert.NoError(t, err)
	assert.Equal(t, user.IsAdmin, false)
}

// 测试登录用户
func TestLoginUser(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "张三",
		Email:    "zhangsan@gmail.com",
		Password: "abcdef",
	}
	_, err = userService.RegisterUser(registerRequest)
	assert.NoError(t, err)
	loginRequest := &model.LoginRequest{
		Username: "张三",
		Password: "abcdef",
	}
	user, token, err := userService.LoginUser(loginRequest)
	assert.NoError(t, err)
	assert.NotEmpty(t, user)
	assert.NotEmpty(t, token)
}

// 测试忘记密码
func TestForgetPassword(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "张三",
		Email:    "zhangsan@gmail.com",
		Password: "abcdef",
	}
	user, err := userService.RegisterUser(registerRequest)
	assert.NoError(t, err)
	forgetRequest := &model.ForgetPasswordRequest{
		NewPassword: "123456",
		Username:    "张三",
		Email:       "zhangsan@gmail.com",
	}
	user, err = userService.ForgetPassword(forgetRequest)
	assert.NoError(t, err)
	assert.NotEmpty(t, user)
	assert.Equal(t, user.Username, "张三")
	assert.Equal(t, user.Email, "zhangsan@gmail.com")
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(forgetRequest.NewPassword))
	assert.NoError(t, err)
}

// 测试根据ID获取用户
func TestGetUserByID(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "李四",
		Email:    "lisi@qq.com",
		Password: "111111",
	}
	user, err := userService.RegisterUser(registerRequest)
	assert.NoError(t, err)
	user, err = userService.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, user)
	assert.Equal(t, user.Username, "李四")
	assert.Equal(t, user.Email, "lisi@qq.com")
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(registerRequest.Password))
	assert.NoError(t, err)
	assert.Equal(t, user.IsAdmin, false)
}

// 测试根据用户名获取用户
func TestGetUserByUsername(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	registerRequest := &model.RegisterRequest{
		Username: "王五",
		Email:    "wangwu@163.com",
		Password: "222222",
	}
	_, err = userService.RegisterUser(registerRequest)
	assert.NoError(t, err)
	user, err := userService.GetUserByUsername("王五")
	assert.NoError(t, err)
	assert.NotEmpty(t, user)
	assert.Equal(t, user.Username, "王五")
	assert.Equal(t, user.Email, "wangwu@163.com")
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(registerRequest.Password))
	assert.NoError(t, err)
	assert.Equal(t, user.IsAdmin, false)
}

// 测试更新用户
func TestUpdateUser(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 先注册一个用户
	registerRequest := &model.RegisterRequest{
		Username: "更新前用户",
		Email:    "before@test.com",
		Password: "oldpassword",
	}
	user, err := userService.RegisterUser(registerRequest)
	assert.NoError(t, err)

	// 更新用户信息
	updateRequest := &model.UpdateUserRequest{
		ID:       user.ID,
		Username: "更新后用户",
		Email:    "after@test.com",
		Password: "newpassword",
	}
	err = userService.UpdateUser(updateRequest)
	assert.NoError(t, err)

	// 验证更新结果
	updatedUser, err := userService.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "更新后用户", updatedUser.Username)
	assert.Equal(t, "after@test.com", updatedUser.Email)
	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.PasswordHash), []byte("newpassword"))
	assert.NoError(t, err)
}

// 测试获取不存在的用户的分页查询
func TestGetAllUsersEmpty(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 在没有任何用户的情况下测试分页查询
	users, pageTotal, err := userService.GetAllUsers(1, 10)
	assert.NoError(t, err)
	assert.Empty(t, users)
	assert.Equal(t, pageTotal, int64(0))
}

// 测试分页获取所有用户
func TestGetAllUsers(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 先创建几个测试用户
	for i := 1; i <= 5; i++ {
		registerRequest := &model.RegisterRequest{
			Username: fmt.Sprintf("用户%d", i),
			Email:    fmt.Sprintf("user%d@test.com", i),
			Password: "123456",
		}
		_, err := userService.RegisterUser(registerRequest)
		assert.NoError(t, err)
	}

	// 测试第一页
	users, pageTotal, err := userService.GetAllUsers(1, 3)
	assert.NoError(t, err)
	assert.Len(t, users, 3)
	assert.Equal(t, pageTotal, int64(2))

	// 测试第二页
	users, pageTotal, err = userService.GetAllUsers(2, 3)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, pageTotal, int64(2))

	// 测试超出范围页
	users, pageTotal, err = userService.GetAllUsers(3, 3)
	assert.NoError(t, err)
	assert.Len(t, users, 0)
	assert.Equal(t, pageTotal, int64(2))
}

// 测试给予管理员权限
func TestGrantAdminPermission(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 先注册一个普通用户
	registerRequest := &model.RegisterRequest{
		Username: "普通用户",
		Email:    "normal@test.com",
		Password: "123456",
	}
	user, err := userService.RegisterUser(registerRequest)
	assert.NoError(t, err)
	assert.Equal(t, false, user.IsAdmin)

	// 给予管理员权限
	err = userService.GrantAdminPermission(user.ID)
	assert.NoError(t, err)

	// 验证管理员权限是否生效
	adminUser, err := userService.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, adminUser.IsAdmin)
}

// 测试撤销管理员权限
func TestRevokeAdminPermission(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 先注册一个用户并给予管理员权限
	registerRequest := &model.RegisterRequest{
		Username: "管理员用户",
		Email:    "admin@test.com",
		Password: "123456",
	}
	user, err := userService.RegisterUser(registerRequest)
	assert.NoError(t, err)

	// 给予管理员权限
	err = userService.GrantAdminPermission(user.ID)
	assert.NoError(t, err)

	// 确认是管理员
	adminUser, err := userService.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, adminUser.IsAdmin)

	// 撤销管理员权限
	err = userService.RevokeAdminPermission(user.ID)
	assert.NoError(t, err)

	// 验证管理员权限已撤销
	normalUser, err := userService.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, false, normalUser.IsAdmin)
}

// 测试判断用户是否为管理员
func TestIsAdmin(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 先注册一个用户并给予管理员权限
	registerRequest := &model.RegisterRequest{
		Username: "管理员用户",
		Email:    "admin@test.com",
		Password: "123456",
	}
	user, err := userService.RegisterUser(registerRequest)
	assert.NoError(t, err)

	// 验证用户是否为管理员
	isAdmin, err := userService.IsAdmin(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, false, isAdmin)

	// 给予管理员权限
	err = userService.GrantAdminPermission(user.ID)
	assert.NoError(t, err)

	// 验证管理员权限是否生效
	adminUser, err := userService.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, adminUser.IsAdmin)

	// 验证用户是否为管理员
	isAdmin, err = userService.IsAdmin(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isAdmin)
}

// 测试删除用户
func TestDeleteUser(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 先注册一个用户
	registerRequest := &model.RegisterRequest{
		Username: "待删除用户",
		Email:    "delete@test.com",
		Password: "123456",
	}
	user, err := userService.RegisterUser(registerRequest)
	assert.NoError(t, err)

	// 确认用户存在
	_, err = userService.GetUserByID(user.ID)
	assert.NoError(t, err)

	// 删除用户
	err = userService.DeleteUser(user.ID)
	assert.NoError(t, err)

	// 验证用户已被删除
	_, err = userService.GetUserByID(user.ID)
	assert.Error(t, err) // 应该返回错误，因为用户不存在了
}

// 测试注册用户 - 用户名已存在的错误场景
func TestRegisterUserDuplicateUsername(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 先注册一个用户
	registerRequest1 := &model.RegisterRequest{
		Username: "重复用户名",
		Email:    "duplicate_user1@test.com",
		Password: "123456",
	}
	_, err = userService.RegisterUser(registerRequest1)
	assert.NoError(t, err)

	// 尝试用相同的用户名注册
	registerRequest2 := &model.RegisterRequest{
		Username: "重复用户名", // 相同的用户名
		Email:    "duplicate_user2@test.com",
		Password: "789012",
	}
	_, err = userService.RegisterUser(registerRequest2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户名已存在")
}

// 测试注册用户 - 邮箱已存在的错误场景
func TestRegisterUserDuplicateEmail(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 先注册一个用户
	registerRequest1 := &model.RegisterRequest{
		Username: "用户111",
		Email:    "duplicate@test.com",
		Password: "123456",
	}
	_, err = userService.RegisterUser(registerRequest1)
	assert.NoError(t, err)

	// 尝试用相同的邮箱注册
	registerRequest2 := &model.RegisterRequest{
		Username: "用户222",              // 不同的用户名
		Email:    "duplicate@test.com", // 相同的邮箱
		Password: "789012",
	}
	_, err = userService.RegisterUser(registerRequest2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "邮箱已存在")
}

// 测试登录用户 - 用户名不存在的错误场景
func TestLoginUserNotFound(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	loginRequest := &model.LoginRequest{
		Username: "不存在的用户名",
		Password: "123456",
	}
	_, _, err = userService.LoginUser(loginRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户名不存在")
}

// 测试登录用户 - 密码错误的错误场景
func TestLoginUserWrongPassword(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 先注册一个用户
	registerRequest := &model.RegisterRequest{
		Username: "密码测试用户",
		Email:    "password@test.com",
		Password: "correctpassword",
	}
	_, err = userService.RegisterUser(registerRequest)
	assert.NoError(t, err)

	// 尝试用错误的密码登录
	loginRequest := &model.LoginRequest{
		Username: "密码测试用户",
		Password: "wrongpassword",
	}
	_, _, err = userService.LoginUser(loginRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "密码错误")
}

// 测试根据ID获取不存在的用户
func TestGetUserByIDNotFound(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 尝试获取一个不存在的用户ID
	_, err = userService.GetUserByID(99999)
	assert.Error(t, err)
}

// 测试根据用户名获取不存在的用户
func TestGetUserByUsernameNotFound(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 尝试获取一个不存在的用户名
	_, err = userService.GetUserByUsername("不存在的用户名")
	assert.Error(t, err)
}

// 测试更新用户时用户名已存在的错误场景
func TestUpdateUserDuplicateUsername(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 先注册两个用户
	user1 := &model.RegisterRequest{
		Username: "用户A",
		Email:    "usera@test.com",
		Password: "123456",
	}
	_, err = userService.RegisterUser(user1)
	assert.NoError(t, err)

	user2 := &model.RegisterRequest{
		Username: "用户B",
		Email:    "userb@test.com",
		Password: "123456",
	}
	_, err = userService.RegisterUser(user2)
	assert.NoError(t, err)

	// 尝试将用户1的用户名改为用户2的用户名
	updateRequest := &model.UpdateUserRequest{
		Username: "用户B", // 与用户2的用户名重复
		Email:    "userc@test.com",
		Password: "789012",
	}
	err = userService.UpdateUser(updateRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户名已存在")
}

// 测试更新用户时邮箱已存在的错误场景
func TestUpdateUserDuplicateEmail(t *testing.T) {
	err := ClearTables(dbConfig, []string{"users"})
	if err != nil {
		panic(err)
	}
	// 先注册两个用户
	user1 := &model.RegisterRequest{
		Username: "用户C",
		Email:    "userc@test.com",
		Password: "123456",
	}
	_, err = userService.RegisterUser(user1)
	assert.NoError(t, err)

	user2 := &model.RegisterRequest{
		Username: "用户D",
		Email:    "userd@test.com",
		Password: "123456",
	}
	_, err = userService.RegisterUser(user2)
	assert.NoError(t, err)

	// 尝试将用户1的邮箱改为用户2的邮箱
	updateRequest := &model.UpdateUserRequest{
		Username: "用户E",
		Email:    "userd@test.com", // 与用户2的邮箱重复
		Password: "789012",
	}
	err = userService.UpdateUser(updateRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "邮箱已存在")
}
