package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"net/http"
)

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Herro World!")
	})

	port := os.Getenv("PORT")
	r.Run(":" + port)
}
