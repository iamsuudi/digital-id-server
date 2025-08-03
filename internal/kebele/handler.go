package kebele

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

func (h *Handler) CreateKebele(c *gin.Context) {
	var input types.KebeleInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	city, err := h.service.CreateKebele(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create city: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, city)
}

func (h *Handler) UpdateKebele(c *gin.Context) {
	var input types.KebeleInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kebele ID"})
		return
	}

	err = h.service.UpdateKebele(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kebele not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update kebele"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetKebele(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kebele ID"})
		return
	}

	kebele,err := h.service.GetKebele(c, id)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kebele not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch kebele"})
		}
		return
	}

	c.JSON(http.StatusOK, kebele)
}

func(h *Handler) GetKebeleCashiers(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kebele ID"})
		return
	}

	cashiers,err := h.service.GetCashiers(c, id)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kebele not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cashiers"})
		}
		return
	}

	c.JSON(http.StatusOK, cashiers)
}

func(h *Handler) GetKebeleEncoders(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kebele ID"})
		return
	}

	encoders,err := h.service.GetCashiers(c, id)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kebele not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch encoders"})
		}
		return
	}

	c.JSON(http.StatusOK, encoders)
}

func (h *Handler) GetKebeles(c *gin.Context) {
	limit, offset, query := utils.PaginationHelper(c)

	if strings.TrimSpace(query) == "" {
		count, kebeles, err := h.service.GetKebeles(c.Request.Context(), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch kebeles"})
			return
		}
		if kebeles == nil {
			kebeles = []repository.ListKebelesRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"kebeles": kebeles,
			"count":   count,
		})
	} else {
		count, kebeles, err := h.service.SearchKebeles(c, limit, offset, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search kebeles"})
			return
		}
		if kebeles == nil {
			kebeles = []repository.SearchKebelesRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"kebeles": kebeles,
			"count":   count,
		})
	}
}
