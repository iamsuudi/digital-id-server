package city

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

func (h *Handler) CreateCity(c *gin.Context) {
	var input types.CityInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	city, err := h.service.CreateCity(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create city: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, city)
}

func (h *Handler) UpdateCityInfo(c *gin.Context) {
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

	err = h.service.UpdateCityInfo(c.Request.Context(), id, input.Name, input.Lat, input.Lon)
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

func (h *Handler) DeleteCity(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city ID"})
		return
	}

	err = h.service.DeleteCity(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "City not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete city"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetSubCities(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city ID"})
		return
	}

	subcities, err := h.service.GetSubCities(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "City not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subcities"})
		}
		return
	}

	c.JSON(http.StatusOK, subcities)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city ID"})
		return
	}

	err = h.service.AssignAdmin(c.Request.Context(), id, input.StaffID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "city not found"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city ID"})
		return
	}

	err = h.service.RemoveStaff(c.Request.Context(), input.StaffID)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "city not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove staff"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetCity(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city ID"})
		return
	}

	city, err := h.service.GetCity(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "City not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch city"})
		}
		return
	}

	c.JSON(http.StatusOK, city)
}

func (h *Handler) GetCities(c *gin.Context) {
	limit, offset, query := utils.PaginationHelper(c)

	if strings.TrimSpace(query) == "" {
		count, cities, err := h.service.GetCities(c.Request.Context(), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cities"})
			return
		}
		if cities == nil {
			cities = []repository.ListCitiesRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"cities": cities,
			"count":  count,
		})
	} else {
		count, cities, err := h.service.SearchCities(c, limit, offset, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search cities"})
			return
		}
		if cities == nil {
			cities = []repository.SearchCitiesRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"cities": cities,
			"count":  count,
		})
	}
}
