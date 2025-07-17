package user

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamsuudi/digital-id-server/internal/repository"
	"github.com/iamsuudi/digital-id-server/shared/utils"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetUser(c *gin.Context) {
	by := c.Query("by")
	value := c.Query("value")

	switch by {
	case "id":
		id, err := uuid.Parse(value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city ID"})
			return
		}

		resident, err := h.service.GetUserById(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Resident not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			}
			return
		}
		c.JSON(http.StatusOK, resident)

	case "email":
		resident, err := h.service.GetUserByEmail(c.Request.Context(), value)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Resident not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resident"})
			}
			return
		}
		c.JSON(http.StatusOK, resident)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query ?by= must be 'id' or 'email'"})
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	limit, _, offset, limitErr, pageErr := utils.PaginationHelper(c)
	query := c.Query("query")

	if limitErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rows per page"})
		return
	}
	if pageErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	if strings.TrimSpace(query) == "" {
		count, users, err := h.service.GetAll(c.Request.Context(), limit, offset, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		if users == nil {
			users = []repository.ListUsersRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"count": count,
		})
	} else {
		count, users, err := h.service.SearchUsers(c, limit, offset, query)
		fmt.Print(err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
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
