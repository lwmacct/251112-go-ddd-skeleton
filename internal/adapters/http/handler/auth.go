package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-ddd-skeleton/internal/adapters/http/response"
	"github.com/yourusername/go-ddd-skeleton/internal/application/auth"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, resp)
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c)
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req auth.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, resp)
}

// Enable2FA 启用2FA
func (h *AuthHandler) Enable2FA(c *gin.Context) {
	userID := c.GetString("userID")
	email := c.Query("email")

	resp, err := h.authService.Enable2FA(c.Request.Context(), userID, email)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, resp)
}

// Verify2FA 验证2FA
func (h *AuthHandler) Verify2FA(c *gin.Context) {
	userID := c.GetString("userID")

	var req auth.Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	if err := h.authService.Verify2FA(c.Request.Context(), userID, req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "2FA enabled successfully"})
}

// CreatePAT 创建PAT
func (h *AuthHandler) CreatePAT(c *gin.Context) {
	userID := c.GetString("userID")

	var req auth.CreatePATRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	resp, err := h.authService.CreatePAT(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, resp)
}

// GetSessions 获取会话
func (h *AuthHandler) GetSessions(c *gin.Context) {
	userID := c.GetString("userID")

	sessions, err := h.authService.GetUserSessions(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, sessions)
}
