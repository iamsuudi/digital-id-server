package audit

import (
	"errors"
	"fmt"
	"net/http"

	"digital-id-server/internal/cache"
	"digital-id-server/internal/repository"
	"digital-id-server/shared/utils"

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

func (h *Handler) GetLogs(c *gin.Context) {
	limit, offset, _ := utils.PaginationHelper(c)

	str, _ := c.Get("user_id")
	id, _ := str.(uuid.UUID)
	user, _ := h.cache.GetUser(c, id)

	count, logs, err := h.service.ListAuditLogs(c.Request.Context(), limit, offset, user.CityID, user.SubcityID, user.KebeleID)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}
	if logs == nil {
		logs = []repository.ListAuditLogsRow{}
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"count": count,
	})
}

func (h *Handler) GetLog(c *gin.Context) {
	str, _ := c.Get("user_id")
	id, _ := str.(uuid.UUID)
	user, _ := h.cache.GetUser(c, id)

	raw := c.Param("id")
	logID, err := uuid.Parse(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audit log ID"})
		return
	}

	log, err := h.service.GetAuditLog(c.Request.Context(), logID, user.CityID, user.SubcityID, user.KebeleID)
	if err != nil {
		fmt.Print(err)
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Audit log not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit log"})
		}
		return
	}

	c.JSON(http.StatusOK, log)
}
