package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-ddd-skeleton/internal/adapters/http/response"
	"github.com/yourusername/go-ddd-skeleton/internal/application/user"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *user.Service
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register 注册用户
func (h *UserHandler) Register(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	dto, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, dto)
}

// GetUser 获取用户
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	dto, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req user.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	dto, err := h.userService.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c)
}

// ListUsers 列出用户
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	req := user.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := h.userService.ListUsers(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, resp)
}
