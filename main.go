package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var ws = websocket.Upgrader{}

var recentSales = make(map[string](chan recentSaleData))

type recentSaleData struct {
	SchoolID string `json:"school_id"`
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
			fmt.Println("Handling sale")

			test := recentSaleData{SchoolID: "1232123"}

			if recentSales[test.SchoolID] == nil {
				//init channel
				recentSales[test.SchoolID] = make(chan recentSaleData, 10)
			}

			fmt.Println("Pushing to channel")
			recentSales[test.SchoolID] <- test

			fmt.Println("Success")
			c.JSON(200, nil)
			return
		})
	}

	socket := r.Group("/school")
	{
		go socket.GET("/:id", func(c *gin.Context) {
			//wsHandler(c.Writer, c.Request)
			id := c.Param("id")
			fmt.Printf("Received request for school/%s\n", id)
			fmt.Println("Upgrading endpoint to websocket")
			conn, err := ws.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				fmt.Printf("Failed to upgrade websocket: %+v\n", err)
				return
			}

			fmt.Println("Success, starting for loop")
			fmt.Printf("looking for school channel ID: %s\n", id)

			defer conn.Close()

			for {
				time.Sleep(2 * time.Second)
				schoolChannel := recentSales[id]
				if schoolChannel == nil {
					fmt.Printf("school channel for ID:%s is nil\n", id)
					continue
				}
				msg := <-schoolChannel
				bytes, err := json.Marshal(msg)
				if err != nil {
					fmt.Println("json.Marshal had an error while Marshaling %+w", err)
					continue
				}
				fmt.Printf("Writing message in bytes%v\n", msg)
				conn.WriteMessage(1, bytes)
			}
		})
	}

	r.Run("localhost:8080")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

}
