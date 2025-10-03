package api

import (
	"fmt"
	"strconv"

	"server/global"
	"server/model"
	"server/service"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	userService := service.NewUserService(global.GetDB())
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var registerRequest model.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		utils.BadRequest(c, fmt.Sprintf("注册请求绑定失败: %v", err))
		return
	}
	user, err := h.userService.RegisterUser(&registerRequest)
	if err != nil {
		utils.ServerError(c, fmt.Sprintf("注册用户失败: %v", err))
		return
	}
	utils.Success(c, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginRequest model.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		utils.BadRequest(c, fmt.Sprintf("登录请求绑定失败: %v", err))
		return
	}
	user, token, err := h.userService.LoginUser(&loginRequest)
	if err != nil {
		utils.ServerError(c, fmt.Sprintf("登录用户失败: %v", err))
		return
	}
	utils.Success(c, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		utils.BadRequest(c, fmt.Sprintf("用户ID绑定失败: %v", err))
		return
	}
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		utils.ServerError(c, fmt.Sprintf("获取用户信息失败: %v", err))
		return
	}
	utils.Success(c, user)
}

func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		utils.BadRequest(c, "用户名不能为空")
		return
	}
	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		utils.ServerError(c, fmt.Sprintf("获取用户信息失败: %v", err))
		return
	}
	utils.Success(c, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// 用operatorID来和UserID区分
	operatorID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "操作者信息获取失败")
		return
	}
	isAdmin, err := h.userService.IsAdmin(operatorID.(int))
	if err != nil {
		utils.ServerError(c, fmt.Sprintf("判断用户是否为管理员失败: %v", err))
		return
	}
	if !isAdmin {
		utils.Forbidden(c, "无权限访问该接口")
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		utils.BadRequest(c, fmt.Sprintf("页码绑定失败: %v", err))
		return
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		utils.BadRequest(c, fmt.Sprintf("页大小绑定失败: %v", err))
		return
	}
	users, pageTotal, err := h.userService.GetAllUsers(page, pageSize)
	if err != nil {
		utils.ServerError(c, fmt.Sprintf("获取用户信息失败: %v", err))
		return
	}
	utils.Success(c, gin.H{
		"users":     users,
		"pageTotal": pageTotal,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var updateUserRequest model.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateUserRequest); err != nil {
		utils.BadRequest(c, fmt.Sprintf("更新用户请求绑定失败: %v", err))
		return
	}
	err := h.userService.UpdateUser(&updateUserRequest)
	if err != nil {
		utils.ServerError(c, fmt.Sprintf("更新用户失败: %v", err))
		return
	}
	utils.Success(c, "用户更新成功")
}

func (h *UserHandler) GrantAdminPermission(c *gin.Context) {
	operatorID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "操作者信息获取失败")
		return
	}
	isAdmin, err := h.userService.IsAdmin(operatorID.(int))
	if err != nil {
		utils.ServerError(c, fmt.Sprintf("判断用户是否为管理员失败: %v", err))
		return
	}
	if !isAdmin {
		utils.Forbidden(c, "无权限访问该接口")
		return
	}
	userID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		utils.BadRequest(c, fmt.Sprintf("用户ID绑定失败: %v", err))
		return
	}
	err = h.userService.GrantAdminPermission(userID)
	if err != nil {
		utils.ServerError(c, fmt.Sprintf("授予管理员权限失败: %v", err))
		return
	}
	utils.Success(c, "管理员权限授予成功")
}

func (h *UserHandler) RevokeAdminPermission(c *gin.Context) {
	operatorID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "操作者信息获取失败")
		return
	}
	isAdmin, err := h.userService.IsAdmin(operatorID.(int))
	if err != nil {
		utils.ServerError(c, fmt.Sprintf("判断用户是否为管理员失败: %v", err))
		return
	}
	if !isAdmin {
		utils.Forbidden(c, "无权限访问该接口")
		return
	}
	userID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		utils.BadRequest(c, fmt.Sprintf("用户ID绑定失败: %v", err))
		return
	}
	err = h.userService.RevokeAdminPermission(userID)
	if err != nil {
		utils.ServerError(c, fmt.Sprintf("撤销管理员权限失败: %v", err))
		return
	}
	utils.Success(c, "管理员权限撤销成功")
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	operatorID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "操作者信息获取失败")
		return
	}
	isAdmin, err := h.userService.IsAdmin(operatorID.(int))
	if err != nil {
		utils.ServerError(c, fmt.Sprintf("判断用户是否为管理员失败: %v", err))
		return
	}
	if !isAdmin {
		utils.Forbidden(c, "无权限访问该接口")
		return
	}
	userID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		utils.BadRequest(c, fmt.Sprintf("用户ID绑定失败: %v", err))
		return
	}
	err = h.userService.DeleteUser(userID)
	if err != nil {
		utils.ServerError(c, fmt.Sprintf("删除用户失败: %v", err))
		return
	}
	utils.Success(c, "用户删除成功")
}
