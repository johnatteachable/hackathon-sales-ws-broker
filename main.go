package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var ws = websocket.Upgrader{}

var recentSales = make(map[int](chan recentSaleData))

type (
	recentSaleData struct {
		ID          int     `json:"id"`
		SchoolID    int     `json:"school_id"`
		Product     product `json:"product"`
		CountryCode string  `json:"country_code"`
		Price       string  `json:"price"`
		User        user    `json:"user"`
	}

	product struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		ThumbnailURL string `json:"thumbnail_url"`
		URL          string `json:"url"`
	}

	user struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

// CORSConfig configuration struct for the CORS middleware
type corsConfig struct {
	MaxAgeHours  int64
	AllowOrigins []string
}

// CORS Cross Origin Resource Sharing middleware
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		MaxAge:       time.Duration(12) * time.Hour,
		AllowOrigins: []string{"*"},
	})
}

func main() {
	fmt.Println("Starting sales websocket broker")

	r := gin.New()
	r.Use(CORS())

	hc := r.Group("healthcheck")
	{
		hc.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "up",
			})
		})
	}

	sales := r.Group("sale")
	{
		// API to consume RecentSaleData
		go sales.POST("", func(c *gin.Context) {
			fmt.Println("Handling sale")
			var req recentSaleData

			fmt.Printf("parsing request: %+v\n", req)
			if err := c.BindJSON(&req); err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
			fmt.Printf("received request: %+v\n", req)

			if recentSales[req.SchoolID] == nil {
				//init channel
				recentSales[req.SchoolID] = make(chan recentSaleData, 100)
			}

			fmt.Println("Pushing request to channel")
			recentSales[req.SchoolID] <- req

			fmt.Println("Success")
			c.JSON(200, nil)
			return
		})
	}

	socket := r.Group("/school")
	{
		// Establish connection to school websocket and forward data associated channel
		go socket.GET("/:id", func(c *gin.Context) {
			stringID := c.Param("id")
			fmt.Printf("Received request for school/%s\n", stringID)
			id, err := strconv.Atoi(stringID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "school id must be an int"})
			}

			ws.CheckOrigin = CheckOrigin
			fmt.Println("Upgrading endpoint to websocket")
			conn, err := ws.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				fmt.Printf("Failed to upgrade websocket: %+v\n", err)
				return
			}

			fmt.Println("Success, starting for loop")
			fmt.Printf("looking for school channel ID: %d\n", id)

			defer conn.Close()

			for {
				time.Sleep(500 * time.Millisecond)
				schoolChannel := recentSales[id]
				if schoolChannel == nil {
					//fmt.Printf("school channel for ID:%d is nil\n", id)
					continue
				}
				msg := <-schoolChannel
				bytes, err := json.Marshal(msg)
				if err != nil {
					//fmt.Println("json.Marshal had an error while Marshaling %+w", err)
					continue
				}
				fmt.Printf("Writing message in bytes%v\n", msg)
				conn.WriteMessage(1, bytes)
			}
		})
	}

	r.Run("localhost:8080")
}

func CheckOrigin(r *http.Request) bool {
	return true
}
