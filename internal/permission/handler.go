package permission

import (
	"digital-id-server/internal/cache"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

type Handler struct {
	service *Service
	cache   *cache.Cache
}

func NewHandler(s *Service, c *cache.Cache) *Handler {
	return &Handler{service: s, cache: c}
}

func (h *Handler) CreatePermission(c *gin.Context) {
	var input struct {
		Name        string `form:"name"   binding:"required"`
		Label       string `form:"label"  binding:"required"`
		Description string `form:"description"`
	}
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(input)

	err := h.service.CreatePermission(c, input.Name, input.Label, input.Description)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create permission"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Permission created successfully"})
}

func (h *Handler) DeletePermission(c *gin.Context) {
	name := c.Param("name")

	err := h.service.DeletePermission(c, name)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete permission"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission deleted successfully"})
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

	permissions, err := h.service.GetAssignablePermissionsForActor(c, id)

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

func (h *Handler) GetUniversalPermissionsForUser(c *gin.Context) {
	raw, _ := c.Get("user_id")
	id, _ := raw.(uuid.UUID)
	user, _ := h.cache.GetUser(c, id)
	target_role := c.Query("role")
	target_id, err := uuid.Parse(c.Query("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	perms, err := h.service.GetUniversalPermissionsForUser(c, user.RoleSlug, target_role, target_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permissions not found"})
		} else {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permissions"})
		}
		return
	}

	c.JSON(http.StatusOK, perms)
}

func (h *Handler) OverridePermission(c *gin.Context) {
	raw, _ := c.Get("user_id")
	id, _ := raw.(uuid.UUID)

	var input struct {
		TargetID   uuid.UUID `json:"id" binding:"required"`
		Permission string    `json:"permission" binding:"required"`
		Override   bool      `json:"override"`
	}
	err := c.BindJSON(&input)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = h.service.OverrideUserPermission(c, id, input.TargetID, input.Permission, input.Override)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to override permission"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Permission overridden successfully"})
}

func (h *Handler) RemoveOverride(c *gin.Context) {
	raw, _ := c.Get("user_id")
	id, _ := raw.(uuid.UUID)

	var input struct {
		TargetID   uuid.UUID `json:"id" binding:"required"`
		Permission string    `json:"permission" binding:"required"`
	}
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = h.service.RemoveUserPermissionOverride(c, id, input.TargetID, input.Permission)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove permission override"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Permission override removed successfully"})
}
