package subcity

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"digital-id-server/shared/types"
	"digital-id-server/internal/repository"
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
	limit, err := strconv.Atoi(c.Query("rows"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rows per page"})
		return
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	
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
		"count": count,
	})
}
