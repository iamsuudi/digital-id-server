package analytics

import (
	"digital-id-server/internal/auth"
	"digital-id-server/internal/cache"
	"digital-id-server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries, cache *cache.Cache) {
	service := NewService(db, q)
	handler := NewHandler(service, cache)

	r := rg.Group("/analytics", auth.Authenticate())
	{
		r.GET("/age-group", handler.AgeGroup)
		r.GET("/gender", handler.Gender)
		r.GET("/gender-age-group", handler.GenderAgeGroup)
	}
}
