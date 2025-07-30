package auth

import (
	"fmt"
	"net/http"
	"time"

	"digital-id-server/shared/types"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) Login(c *gin.Context) {
	type LoginInput struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required,min=6"`
	}

	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.Authenticate(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, err := GenerateJWT(user.ID, user.RoleSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken := GenerateRandomToken(64)
	refreshExpires := time.Now().Add(7 * 24 * time.Hour)

	if err := h.Service.StoreRefreshToken(c.Request.Context(), user.ID, refreshToken, refreshExpires); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	// Set access token (JWT)
	c.SetCookie("oict_jwt", accessToken, 15*60*60, "/", "", true, true) // 15 minutes
	// Set refresh token securely
	c.SetCookie("oict_refresh_token", refreshToken, 7*24*60*60, "/", "", true, true) // 7 days

	c.JSON(http.StatusOK, user)
}

func (h *Handler) Logout(c *gin.Context) {
	rt, _ := c.Cookie("refresh_token")
	if rt != "" {
		_ = h.Service.DeleteRefreshToken(c.Request.Context(), rt) // Optional cleanup
	}

	c.SetCookie("oict_jwt", "", -1, "/", "", true, true)
	c.SetCookie("oict_refresh_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var input types.UserRegisterInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.RegisterUser(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	rt, err := c.Cookie("oict_refresh_token")
	if err != nil || rt == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token"})
		return
	}

	newAccessToken, newRefreshToken, err := h.Service.RefreshAccessToken(c.Request.Context(), rt)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Set access token (JWT)
	c.SetCookie("oict_jwt", newAccessToken, 15*60*60, "/", "", true, true) // 15 minutes
	// Set refresh token securely
	c.SetCookie("oict_refresh_token", newRefreshToken, 7*24*60*60, "/", "", true, true) // 7 days

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

func (h *Handler) Me(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	user, err := h.Service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Exclude password hash
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}
