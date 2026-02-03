package main

import (
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	initDB() // Khởi tạo database

	http.HandleFunc("/", requireAuth(homeHandler))
	http.HandleFunc("/login", loginPageHandler)
	http.HandleFunc("/api/login", loginAPIHandler)
	http.HandleFunc("/api/register", registerAPIHandler)
	http.HandleFunc("/api/logout", logoutAPIHandler)

	log.Println("Server chạy tại http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
