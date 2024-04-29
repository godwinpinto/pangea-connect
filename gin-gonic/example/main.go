package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pangea "github.com/godwinpinto/gin-gonic/middleware"
)

// A Sample gin gonic rest api
func main() {
	r := gin.Default()

	//A middleware using Pangea IP Intel
	r.Use(pangea.PangeaIpIntel())

	//example endpoint
	r.GET("/get", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	r.Run()
}
