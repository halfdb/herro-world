package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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
