package resident

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	r := rg.Group("/residents")

	r.POST("/", h.RegisterResident)
}
