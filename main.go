package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "login.html")
}

func loginAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "admin" && req.Password == "123456" {
		resp := LoginResponse{Success: true, Message: "Đăng nhập thành công"}
		json.NewEncoder(w).Encode(resp)
		return
	} else {
		resp := LoginResponse{Success: false, Message: "Sai tài khoản hoặc mật khẩu"}
		json.NewEncoder(w).Encode(resp)
	}
}

func main() {
	http.HandleFunc("/login", loginPageHandler)
	http.HandleFunc("/api/login", loginAPIHandler)
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("Server chạy tại http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
