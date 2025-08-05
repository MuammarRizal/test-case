package main

import (
	"fmt"
	"log"
	"net/http"

	"test-case/internal/config"
	"test-case/internal/handlers"
)

func main() {
	config.InitDB()
	defer config.CloseDB()

	setupRoutes()

	fmt.Println("🚀 Server running on http://localhost:8080")
	fmt.Println("📚 Available endpoints:")
	fmt.Println("   - POST /penjualan")
	fmt.Println("   - GET|POST /hitung-pajak")
	fmt.Println("   - POST /hitung-diskon")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupRoutes() {
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/penjualan", handlers.PenjualanHandler)
	http.HandleFunc("/hitung-pajak", handlers.PajakHandler)
	http.HandleFunc("/hitung-diskon", handlers.DiskonHandler)

	http.HandleFunc("/pajak/hitung", handlers.PajakHandler)
	http.HandleFunc("/diskon/hitung", handlers.DiskonHandler)
}