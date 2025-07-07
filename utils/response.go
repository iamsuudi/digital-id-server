package utils

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data any) {
    c.JSON(200, gin.H{"status": "success", "data": data})
}

func Error(c *gin.Context, code int, msg string) {
    c.JSON(code, gin.H{"status": "error", "message": msg})
}
