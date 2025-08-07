package resident

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

func (h *Handler) RegisterResident(c *gin.Context) {
	var input types.ResidentPayload
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}
	
    file, _ := c.FormFile("document")
    if file != nil {
        c.SaveUploadedFile(file, "uploads/"+file.Filename)
    }

    c.JSON(200, gin.H{"ok": true})

	// faceFile, err := c.FormFile("face")
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "missing face image"})
	// 	return
	// }

	// docFile, err := c.FormFile("document")
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "missing document"})
	// 	return
	// }

	// fmt.Println(faceFile, docFile)

	// // TODO: Save files to storage (or mock)
	// faceURL := "/mock/face.jpg"
	// docURL := "/mock/document.pdf"

	// err = h.service.RegisterResident(c.Request.Context(), input, faceURL, docURL)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register resident: " + err.Error()})
	// 	return
	// }

	// c.JSON(http.StatusCreated, gin.H{"message": "resident registered successfully"})
}

func (h *Handler) GetResident(c *gin.Context) {
	raw := c.Param("id")
	id, err := uuid.Parse(raw)
	fmt.Println(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resident ID"})
		return
	}
	fmt.Println(id)

	resident, err := h.service.GetResident(c.Request.Context(), id)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resident not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resident"})
		}
		return
	}

	c.JSON(http.StatusOK, resident)
}

func (h *Handler) GetResidents(c *gin.Context) {
	limit, offset, query := utils.PaginationHelper(c)

	if strings.TrimSpace(query) == "" {
		count, residents, err := h.service.GetResidents(c.Request.Context(), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch residents"})
			return
		}
		if residents == nil {
			residents = []repository.ListResidentsRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"residents": residents,
			"count":  count,
		})
	} else {
		count, residents, err := h.service.SearchResidents(c, limit, offset, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search residents"})
			return
		}
		if residents == nil {
			residents = []repository.SearchResidentsRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"residents": residents,
			"count":  count,
		})
	}
}
