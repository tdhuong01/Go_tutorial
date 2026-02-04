package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	// Load biến môi trường từ .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Không thể load file .env:", err)
	}

	initDB() // Khởi tạo database

	mux := http.NewServeMux()
	mux.HandleFunc("/", requireAuth(homeHandler))
	mux.HandleFunc("/login", loginPageHandler)
	mux.HandleFunc("/api/login", loginAPIHandler)
	mux.HandleFunc("/api/register", registerAPIHandler)
	mux.HandleFunc("/api/logout", logoutAPIHandler)

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5174"}, // React dev server
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	log.Println("Server chạy tại http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", handler))
}
