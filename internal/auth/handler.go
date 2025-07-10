package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler dependencies

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
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.Authenticate(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate access token (JWT)
	accessToken, err := GenerateJWT(int64(user.ID), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Generate refresh token
	refreshToken := GenerateRandomToken(64)
	refreshExpires := time.Now().Add(7 * 24 * time.Hour) // 7 days

	err = h.Service.StoreRefreshToken(c.Request.Context(), int64(user.ID), refreshToken, refreshExpires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	// Set both cookies
	c.SetCookie("jwt", accessToken, 3600*24, "", "", true, true)             // 1 day
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "", "", true, true) // 7 days

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}


// Logout handler
func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// RegisterUser allows anyone to register with a role for now.
// Later, only superadmins/managers will be allowed to assign certain roles.
func (h *Handler) RegisterUser(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to register
	err := h.Service.RegisterUser(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// handler.go
func (h *Handler) RefreshToken(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil || rt == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token"})
		return
	}

	newToken, err := h.Service.RefreshAccessToken(c.Request.Context(), rt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.SetCookie("jwt", newToken, 3600*24, "", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}
