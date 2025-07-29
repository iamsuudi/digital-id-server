package user

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"digital-id-server/internal/cache"
	"digital-id-server/internal/repository"
	"digital-id-server/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	service   *Service
	cache *cache.Cache
}

func NewHandler(s *Service, c *cache.Cache) *Handler {
	return &Handler{service: s, cache: c}
}

func (h *Handler) RequirePermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		str, _ := c.Get("user_id")
		id, _ := str.(uuid.UUID)

		perms, err := h.cache.GetPerms(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "permission"})
			return
		}

		for _, reqPerm := range permissions {
			if !slices.ContainsFunc(perms, func(p repository.GetEffectivePermissionsForUserRow) bool {
				return p.Name == reqPerm
			}) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permission"})
				return
			}
		}
		c.Next()
	}
}

func (h *Handler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.GetUserById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetAll(c *gin.Context) {
	limit, offset, query, limitErr, pageErr := utils.PaginationHelper(c)
	if limitErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rows per page"})
		return
	}
	if pageErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	if strings.TrimSpace(query) == "" {
		str, _ := c.Get("user_id")
		id, _ := str.(uuid.UUID)
		user, _ := h.cache.GetUser(c, id)
		count, users, err := h.service.GetAllUnderScope(c.Request.Context(), limit, offset, query, user.RoleLevelRank, user.CityID, user.SubcityID, user.KebeleID)
			fmt.Println(*user.RoleLevelRank)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		if users == nil {
			users = []repository.ListUsersUnderScopeRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"count": count,
		})
	} else {
		count, users, err := h.service.SearchUsersUnderScope(c, limit, offset, query)
		fmt.Print(err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
			return
		}
		if users == nil {
			users = []repository.SearchUsersUnderScopeRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"count": count,
		})
	}
}

func (h *Handler) GetAllForSuperadmin(c *gin.Context) {
	limit, offset, query, limitErr, pageErr := utils.PaginationHelper(c)
	if limitErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rows per page"})
		return
	}
	if pageErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	if strings.TrimSpace(query) == "" {
		count, users, err := h.service.GetAllForSuperadmin(c.Request.Context(), limit, offset, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		if users == nil {
			users = []repository.ListAllUsersRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"count": count,
		})
	} else {
		count, users, err := h.service.SearchUsersForSuperadmin(c, limit, offset, query)
		if err != nil {
			fmt.Print("super:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
			return
		}
		if users == nil {
			users = []repository.SearchUsersRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"count": count,
		})
	}
}

func (h *Handler) GetByRole(c *gin.Context) {
	limit, offset, query, limitErr, pageErr := utils.PaginationHelper(c)
	role_slug := c.Query("role_slug")
	if limitErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rows per page"})
		return
	}
	if pageErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	if strings.TrimSpace(query) == "" {
		count, users, err := h.service.GetByRole(c.Request.Context(), limit, offset, query, role_slug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		if users == nil {
			users = []repository.ListByRoleRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"count": count,
		})
	} else {
		count, users, err := h.service.SearchByRole(c, limit, offset, query, role_slug)
		if err != nil {
			fmt.Print("super:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
			return
		}
		if users == nil {
			users = []repository.SearchByRoleRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"count": count,
		})
	}
}
