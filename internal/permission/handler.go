package permission

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.service.GetAllPermissions(c)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permissions not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permissions"})
		}
		return
	}

	c.JSON(http.StatusOK, permissions)
}

func (h *Handler) GetAssignablePermissions(c *gin.Context) {
	raw, _ := c.Get("user_id")
	id, _ := raw.(uuid.UUID)
	
	permissions, err := h.service.GetAssignablePermissions(c, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permissions not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permissions"})
		}
		return
	}

	c.JSON(http.StatusOK, permissions)
}
