package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/response"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/application/auth"
)

// Handler 认证处理器
type Handler struct {
	authService *auth.Service
}

// NewHandler 创建认证处理器
func NewHandler(authService *auth.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

// Login 登录
func (h *Handler) Login(c *gin.Context) {
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
func (h *Handler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c)
}

// RefreshToken 刷新令牌
func (h *Handler) RefreshToken(c *gin.Context) {
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
func (h *Handler) Enable2FA(c *gin.Context) {
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
func (h *Handler) Verify2FA(c *gin.Context) {
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
func (h *Handler) CreatePAT(c *gin.Context) {
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
func (h *Handler) GetSessions(c *gin.Context) {
	userID := c.GetString("userID")

	sessions, err := h.authService.GetUserSessions(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, sessions)
}

// RevokeSession 撤销会话
func (h *Handler) RevokeSession(c *gin.Context) {
	userID := c.GetString("userID")
	sessionID := c.Param("id")

	// TODO: 实现撤销会话逻辑
	_ = userID
	_ = sessionID
	response.Success(c, gin.H{"message": "会话已撤销"})
}

// GetPATs 获取个人访问令牌列表
func (h *Handler) GetPATs(c *gin.Context) {
	userID := c.GetString("userID")

	// TODO: 实现获取 PAT 列表逻辑
	_ = userID
	response.Success(c, []gin.H{})
}

// RevokePAT 撤销个人访问令牌
func (h *Handler) RevokePAT(c *gin.Context) {
	userID := c.GetString("userID")
	tokenID := c.Param("id")

	// TODO: 实现撤销 PAT 逻辑
	_ = userID
	_ = tokenID
	response.Success(c, gin.H{"message": "令牌已撤销"})
}

// Disable2FA 禁用双因素认证
func (h *Handler) Disable2FA(c *gin.Context) {
	userID := c.GetString("userID")

	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// TODO: 实现禁用 2FA 逻辑（需要验证密码）
	_ = userID
	response.Success(c, gin.H{"message": "双因素认证已禁用"})
}
