package resident

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterResident(c *gin.Context) {
	var input RegisterResidentInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	faceFile, err := c.FormFile("face")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing face image"})
		return
	}

	docFile, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing document"})
		return
	}

	fmt.Println(faceFile, docFile)

	// TODO: Save files to storage (or mock)
	faceURL := "/mock/face.jpg"
	docURL := "/mock/document.pdf"

	err = h.service.RegisterResident(c.Request.Context(), input, faceURL, docURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register resident: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "resident registered successfully"})
}
