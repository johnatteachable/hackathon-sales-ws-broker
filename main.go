package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type (
	recentSaleData struct {
		ID              int     `json:"id"`
		SchoolID        int     `json:"school_id"`
		Product         product `json:"product"`
		Name            string  `json:"name"`
		CountryCode     string  `json:"country_code"`
		BillingAddress  address `json:"billing_address"`
		ShippingAddress address `json:"shipping_address"`
		Price           int     `json:"price"`
		User            user    `json:"user"`
	}

	product struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	user struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	address struct {
		ID         int    `json:"id"`
		Line1      string `json:"line1"`
		Line2      string `json:"line2"`
		City       string `json:"city"`
		PostalCode string `json:"postal_code"`
		Country    string `json:"country"`
		Region     string `json:"region"`
	}
)

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
			var req recentSaleData
			fmt.Printf("Handling sale\n")
			fmt.Printf("processing request: %+v\n", req)

			if err := c.BindJSON(&req); err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}

			fmt.Printf("received request: %+v\n", req)

			c.JSON(200, req)
		})
	}

	r.Run("localhost:8080")
}
