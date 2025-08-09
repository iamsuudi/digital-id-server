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
	var input types.ResidentPayload
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	prettyJSON, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	} else {
		fmt.Println(string(prettyJSON))
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	docsUrls := make([]string, 3)
	docs := form.File["documents"]
	for _, doc := range docs {
		name := utils.MakeFileName(doc.Filename)
		dst := filepath.Join("uploads", name)
		c.SaveUploadedFile(doc, dst)
		docsUrls = append(docsUrls, name)
	}

	face, err := c.FormFile("face")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing face image"})
		return
	}
	faceUrl := utils.MakeFileName(face.Filename)
	dst := filepath.Join("uploads", faceUrl)
	c.SaveUploadedFile(face, dst)
	
	fmt.Println("Files are saved!")

	err = h.service.RegisterResident(c.Request.Context(), input, docsUrls, faceUrl)
	if err != nil {
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
