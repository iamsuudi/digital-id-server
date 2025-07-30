package user

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"digital-id-server/internal/cache"
	"digital-id-server/internal/repository"
	"digital-id-server/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	service *Service
	cache   *cache.Cache
}

func NewHandler(s *Service, c *cache.Cache) *Handler {
	return &Handler{service: s, cache: c}
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

func (h *Handler) GetUsers(c *gin.Context) {
	limit, offset, query := utils.PaginationHelper(c)

	str, _ := c.Get("user_id")
	id, _ := str.(uuid.UUID)
	user, _ := h.cache.GetUser(c, id)

	if strings.TrimSpace(query) == "" {
		count, users, err := h.service.ListUsersUnderScope(c.Request.Context(), limit, offset, query, user.RoleLevelRank, user.CityID, user.SubcityID, user.KebeleID)
		if err != nil {
			fmt.Print(err)
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
		count, users, err := h.service.SearchUsersUnderScope(c, limit, offset, query, user.RoleLevelRank, user.CityID, user.SubcityID, user.KebeleID)
		fmt.Println(limit, offset, query, count)
		if err != nil {
			fmt.Print(err)
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

func (h *Handler) GetUsersByRole(c *gin.Context) {
	limit, offset, query := utils.PaginationHelper(c)
	role_slug := c.Query("role_slug")

	if strings.TrimSpace(query) == "" {
		count, users, err := h.service.ListUsersByRole(c.Request.Context(), limit, offset, query, role_slug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		if users == nil {
			users = []repository.ListUsersByRoleRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"count": count,
		})
	} else {
		count, users, err := h.service.SearchUsersByRole(c, limit, offset, query, role_slug)
		if err != nil {
			fmt.Print("super:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
			return
		}
		if users == nil {
			users = []repository.SearchUsersByRoleRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"count": count,
		})
	}
}
