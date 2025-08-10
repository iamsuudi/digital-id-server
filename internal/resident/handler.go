package resident

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
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
	// 1. Bind the entire multipart form at once.
	var input types.ResidentPayload
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	// 2. Optional pretty print.
	if b, err := json.MarshalIndent(input, "", "  "); err == nil {
		fmt.Println(string(b))
	}

	// 3. Collect uploaded documents.
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	docs := form.File["documents"]
	docsUrls := make([]string, 0, len(docs))
	for _, doc := range docs {
		name := utils.MakeFileName(doc.Filename)
		dst  := filepath.Join("uploads", name)
		if err := c.SaveUploadedFile(doc, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save document"})
			return
		}
		docsUrls = append(docsUrls, name)
	}

	// 4. Handle face image.
	face, err := c.FormFile("face")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing face image"})
		return
	}
	faceName := utils.MakeFileName(face.Filename)
	if err := c.SaveUploadedFile(face, filepath.Join("uploads", faceName)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save face image"})
		return
	}

	// 5. Call service layer.
	if err := h.service.RegisterResident(c.Request.Context(), input, docsUrls, faceName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "resident registered successfully"})
}

func (h *Handler) GetResident(c *gin.Context) {
	raw := c.Param("id")
	id, err := uuid.Parse(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resident ID"})
		return
	}

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

func (h *Handler) GetResidentDocuments(c *gin.Context) {
	raw := c.Param("id")
	id, err := uuid.Parse(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resident ID"})
		return
	}

	documents, err := h.service.GetResidentDocuments(c.Request.Context(), id)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resident not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		}
		return
	}

	c.JSON(http.StatusOK, documents)
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
			"count":     count,
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
			"count":     count,
		})
	}
}

func (h *Handler) GetUnpaidResidents(c *gin.Context) {
	limit, offset, query := utils.PaginationHelper(c)

	if strings.TrimSpace(query) == "" {
		count, data, err := h.service.GetUnpaidResidents(c.Request.Context(), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
			return
		}
		if data == nil {
			data = []repository.ListUnpaidResidentsRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"data": 	 data,
			"count":     count,
		})
	} else {
		count, data, err := h.service.SearchUnpaidResidents(c, limit, offset, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search data"})
			return
		}
		if data == nil {
			data = []repository.SearchUnpaidResidentsRow{}
		}

		c.JSON(http.StatusOK, gin.H{
			"data": 	 data,
			"count":     count,
		})
	}
}
