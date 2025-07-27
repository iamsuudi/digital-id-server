package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func PaginationHelper(c *gin.Context) (limit, offset int, query string, limitErr, pageErr error) {
	limit, limitErr = strconv.Atoi(c.Query("rows"))
	page, pageErr := strconv.Atoi(c.Query("page"))
	query = c.Query("query")
	
	if limitErr != nil || pageErr != nil {
		return limit, offset, query, limitErr, pageErr
	}
	
	if limit < 1 || limit > 100 {
		limit = 10
	}
	
	if page < 1 {
		page = 1
	}
	
	offset = (page - 1) * limit
	return limit, offset, query,limitErr, pageErr
}