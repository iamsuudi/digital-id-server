package analytics

import (
	"net/http"

	"digital-id-server/internal/cache"
	"digital-id-server/internal/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	cache   *cache.Cache
}

func NewHandler(s *Service, c *cache.Cache) *Handler {
	return &Handler{service: s, cache: c}
}

func (h *Handler) AgeGroup(c *gin.Context) {
	data, err := h.service.GetAgeGroupDistribution(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if data == nil {
		data = []repository.GetAgeGroupDistributionRow{}
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) Gender(c *gin.Context) {
	data, err := h.service.GetGenderDistribution(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if data == nil {
		data = []repository.GetGenderDistributionRow{}
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GenderAgeGroup(c *gin.Context) {
	data, err := h.service.GetGenderAgeGroupDistribution(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if data == nil {
		data = []repository.GetGenderAgeGroupDistributionRow{}
	}

	c.JSON(http.StatusOK, data)
}
