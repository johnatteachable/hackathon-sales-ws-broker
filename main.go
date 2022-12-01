package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type RecentSaleData struct {
}

// Creat API to consume RecentSaleData
// Establish connection to websocket and forward data from API request
func main() {
	fmt.Println("Starting sales websocket broker")

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Gin here :p")
	})

	sales := r.Group("sale")
	{

		sales.POST("", func(c *gin.Context) {
			fmt.Printf("Handling sale")
			c.JSON(200, nil)
		})
	}

	r.Run("localhost:8080")
}
