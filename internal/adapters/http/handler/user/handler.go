package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/response"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/application/user"
	apperrors "github.com/lwmacct/251112-go-ddd-skeleton/internal/shared/errors"
)

// Handler 用户处理器（包含公开、认证、管理员端点）
type Handler struct {
	userService *user.Service
}

// NewHandler 创建用户处理器
func NewHandler(userService *user.Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

// getCurrentUserID 从 context 获取当前用户 ID
func getCurrentUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", apperrors.ErrUnauthorized
	}

	id, ok := userID.(string)
	if !ok {
		return "", apperrors.ErrUnauthorized
	}

	return id, nil
}

// ========== 公开端点 ==========

// Register 注册用户
// POST /api/auth/register
func (h *Handler) Register(c *gin.Context) {
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

// ========== 用户个人中心端点（需要认证）==========

// GetProfile 获取当前用户信息
// GET /api/user
func (h *Handler) GetProfile(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	dto, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto)
}

// UpdateProfile 更新当前用户信息
// PUT /api/user
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

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

// PatchProfile 部分更新当前用户信息
// PATCH /api/user
func (h *Handler) PatchProfile(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

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

// DeleteAccount 注销当前用户账号
// DELETE /api/user
func (h *Handler) DeleteAccount(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c)
}

// ChangePassword 修改密码
// PUT /api/user/password
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req user.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), userID, req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "密码已更新"})
}

// ChangeEmail 修改邮箱
// PUT /api/user/email
func (h *Handler) ChangeEmail(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// TODO: 实现邮箱修改逻辑（需要邮箱验证）
	_ = userID
	response.Success(c, gin.H{"message": "邮箱修改请求已发送"})
}

// UploadAvatar 上传头像
// POST /api/user/avatar
func (h *Handler) UploadAvatar(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		response.Error(c, err)
		return
	}

	// TODO: 实现头像上传逻辑（保存文件、更新用户信息）
	_ = userID
	_ = file
	response.Success(c, gin.H{"message": "头像已上传", "filename": file.Filename})
}

// ========== 管理员端点（需要admin权限）==========

// GetUser 获取指定用户信息
// GET /api/admin/users/:id
func (h *Handler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	dto, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto)
}

// ListUsers 列出所有用户
// GET /api/admin/users
func (h *Handler) ListUsers(c *gin.Context) {
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

// UpdateUser 更新指定用户信息
// PUT /api/admin/users/:id
func (h *Handler) UpdateUser(c *gin.Context) {
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

// DeleteUser 删除指定用户
// DELETE /api/admin/users/:id
func (h *Handler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c)
}

// BanUser 封禁用户
// POST /api/admin/users/:id/ban
func (h *Handler) BanUser(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// TODO: 实现封禁逻辑
	_ = userID
	response.Success(c, gin.H{"message": "用户已封禁", "reason": req.Reason})
}

// UnbanUser 解封用户
// POST /api/admin/users/:id/unban
func (h *Handler) UnbanUser(c *gin.Context) {
	userID := c.Param("id")

	// TODO: 实现解封逻辑
	_ = userID
	response.Success(c, gin.H{"message": "用户已解封"})
}
