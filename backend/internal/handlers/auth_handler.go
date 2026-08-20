package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanayo-dot/KIVU/backend/internal/models"
	"github.com/hanayo-dot/KIVU/backend/internal/services"
)

// AuthHandler handles user registration, login, token refresh, and logout.
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler initializes AuthHandler.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid registration payload: "+err.Error(), 400))
		return
	}

	tokens, err := h.authService.RegisterFarmer(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error(), 400))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(tokens))
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid login payload: "+err.Error(), 400))
		return
	}

	tokens, err := h.authService.LoginFarmer(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(err.Error(), 401))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(tokens))
}

// Refresh handles POST /auth/refresh.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid refresh token payload: "+err.Error(), 400))
		return
	}

	tokens, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(err.Error(), 401))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(tokens))
}

// Logout handles POST /auth/logout.
func (h *AuthHandler) Logout(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.RefreshToken != "" {
		_ = h.authService.LogoutRevokeToken(req.RefreshToken)
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{
		"message": "Logged out successfully. Refresh token revoked.",
	}))
}
