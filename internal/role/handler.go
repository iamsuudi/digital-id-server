package role

import (
	"errors"
	"fmt"
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

func (h *Handler) GetRolesTree(c *gin.Context) {
	role, _ := c.Get("user_role")
	rolesTree, err := h.service.GetRolesTree(c, role.(string))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "RoleTree not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to role tree user"})
		}
		return
	}

	c.JSON(http.StatusOK, rolesTree)
}

func (h *Handler) GetRolesPermissions(c *gin.Context) {
	permissions, err := h.service.GetRolesPermissions(c)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permissions not found"})
		} else {
			fmt.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permissions"})
		}
		return
	}

	c.JSON(http.StatusOK, permissions)
}

func (h *Handler) UpdateRolePermission(c *gin.Context) {
	raw, _ := c.Get("user_id")
	id, _ := raw.(uuid.UUID)

	role_slug := c.Param("role_slug")
	permission := c.Param("permission_name")
	fmt.Println(raw)
	fmt.Println(id)
	fmt.Println(permission)
	
	myScope, err := h.service.CanManipulateRole(c, id, role_slug)
	if !myScope {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not your scope"})
		return
	}

	canGrant, err := h.service.CanGrantPermissionToRole(c, id, role_slug, permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
		return
	}
	if !canGrant {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update"})
		return
	}

	type GrantPermissionBody struct {
		Grant bool `json:"grant"`
	}
	var req GrantPermissionBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Grant {
		err = h.service.GrantPermissionToRole(c, role_slug, permission)
	} else {
		err = h.service.RevokePermissionFromRole(c, role_slug, permission)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role permission"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role permission updated successfully"})
}
