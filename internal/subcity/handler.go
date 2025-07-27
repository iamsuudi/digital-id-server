package subcity

import (
	"errors"
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

	city, err := h.service.CreateSubCity(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create city: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, city)
}

func (h *Handler) UpdateSubCity(c *gin.Context) {
	var input types.SubCityInput
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

	err = h.service.UpdateSubCity(c.Request.Context(), id, input)
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
		count, subcities, err := h.service.GetAll(c.Request.Context(), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subcities"})
			return
		}
		if subcities == nil {
			subcities = []repository.ListSubCitiesRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"subcities": subcities,
			"count":  count,
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
			"count":  count,
		})
	}
}
