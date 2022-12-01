package main

import (
	"fmt"
	"net/http"
)

type RecentSaleData struct {
}

// Creat API to consume RecentSaleData
// Establish connection to websocket and forward data from API request
func main() {
	fmt.Println("Starting sales websocket broker")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Setting up the server!")
	})

	http.HandleFunc("/sale", handleSale)

	http.ListenAndServe(":8080", nil)
}

func handleSale(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling sale")
}
