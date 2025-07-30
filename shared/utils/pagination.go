package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PaginationHelper(c *gin.Context) (limit, offset int, query string) {
	limit, limitErr := strconv.Atoi(c.Query("rows"))
	page, pageErr := strconv.Atoi(c.Query("page"))
	query = c.Query("query")

	if limitErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rows per page"})
		return
	}
	if pageErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}

	offset = (page - 1) * limit

	return limit, offset, query
}
