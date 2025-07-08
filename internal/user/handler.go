package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	generated "github.com/iamsuudi/digital-id-server/db/generated"
	"context"
)

type Handler struct {
	q *generated.Queries
}

func NewHandler(q *generated.Queries) *Handler {
	return &Handler{q}
}

func (h *Handler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	user, err := h.q.GetUserByID(context.Background(), int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.q.CreateUser(context.Background(), generated.CreateUserParams{
		Name:  input.Name,
		Email: input.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}
