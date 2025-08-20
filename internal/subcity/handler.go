package subcity

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"digital-id-server/internal/repository"
	"digital-id-server/shared/types"
	"digital-id-server/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateSubCity(c *gin.Context) {
	var input types.SubCityInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	raw, _ := c.Get("user_id")
	actorID, _ := raw.(uuid.UUID)

	city, err := h.service.CreateSubCity(c.Request.Context(), actorID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create city: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, city)
}

func (h *Handler) UpdateSubCityInfo(c *gin.Context) {
	var input struct {
		Name string   `json:"name" binding:"required"`
		Lat  *float64 `json:"lat"`
		Lon  *float64 `json:"lon"`
	}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city ID"})
		return
	}

	raw, _ := c.Get("user_id")
	actorID, _ := raw.(uuid.UUID)

	err = h.service.UpdateSubCity(c.Request.Context(), actorID, id, input.Name, input.Lat, input.Lon)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "SubCity not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subcity"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) DeleteSubCity(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subcity ID"})
		return
	}

	raw, _ := c.Get("user_id")
	actorID, _ := raw.(uuid.UUID)

	err = h.service.DeleteSubCity(c.Request.Context(), actorID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "SubCity not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subcity"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetKebeles(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city ID"})
		return
	}

	kebeles, err := h.service.GetKebeles(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "SubCity not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch kebeles"})
		}
		return
	}

	c.JSON(http.StatusOK, kebeles)
}

func (h *Handler) AddStaff(c *gin.Context) {
	var input struct {
		StaffID uuid.UUID `json:"staff_id" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subcity ID"})
		return
	}

	raw, _ := c.Get("user_id")
	actorID, _ := raw.(uuid.UUID)

	err = h.service.AssignManager(c.Request.Context(), actorID, id, input.StaffID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subcity not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add staff"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) RemoveStaff(c *gin.Context) {
	var input struct {
		StaffID uuid.UUID `json:"staff_id"`
	}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	idParam := c.Param("id")
	_, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subcity ID"})
		return
	}

	raw, _ := c.Get("user_id")
	actorID, _ := raw.(uuid.UUID)

	err = h.service.RemoveStaff(c.Request.Context(), actorID, input.StaffID)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subcity not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove staff"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetSubCity(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city ID"})
		return
	}

	subcity, err := h.service.GetSubCity(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "SubCity not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch city"})
		}
		return
	}

	c.JSON(http.StatusOK, subcity)
}

func (h *Handler) GetSubCities(c *gin.Context) {
	limit, offset, query := utils.PaginationHelper(c)

	if strings.TrimSpace(query) == "" {
		count, subcities, err := h.service.GetSubCities(c.Request.Context(), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subcities"})
			return
		}
		if subcities == nil {
			subcities = []repository.ListSubCitiesRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"subcities": subcities,
			"count":     count,
		})
	} else {
		count, subcities, err := h.service.SearchSubCities(c, limit, offset, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search subcities"})
			return
		}
		if subcities == nil {
			subcities = []repository.SearchSubCitiesRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"subcities": subcities,
			"count":     count,
		})
	}
}

func (h *Handler) GetSubCityAudit(c *gin.Context) {
	raw := c.Param("id")
	id, err := uuid.Parse(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subcity ID"})
		return
	}

	limit, offset, _ := utils.PaginationHelper(c)

	count, logs, err := h.service.ListAuditLogs(c.Request.Context(), limit, offset, id)

	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}
	if logs == nil {
		logs = []repository.ListSubCityAuditLogsRow{}
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"count": count,
	})
}
