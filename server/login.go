package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type RegisterRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Sinh JWT token
func generateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Handler trả về trang login.html
func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:5174", http.StatusSeeOther)
}

// Handler xử lý API login
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

	// Tìm user trong database
	var user User
	row := db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", req.Username)
	err = row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			resp := LoginResponse{Success: false, Message: "Sai tài khoản hoặc mật khẩu"}
			json.NewEncoder(w).Encode(resp)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Kiểm tra mật khẩu
	if !checkPassword(user.Password, req.Password) {
		resp := LoginResponse{Success: false, Message: "Sai tài khoản hoặc mật khẩu"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	token, err := generateJWT(req.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	// Set token vào cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Đặt true nếu dùng HTTPS
		SameSite: http.SameSiteLaxMode,
	})
	resp := LoginResponse{Success: true, Message: "Đăng nhập thành công", Token: token}
	json.NewEncoder(w).Encode(resp)
}

// Handler xử lý API đăng ký
func registerAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" || req.Email == "" {
		resp := RegisterResponse{Success: false, Message: "Tên đăng nhập, email và mật khẩu không được để trống"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if req.Password != req.ConfirmPassword {
		resp := RegisterResponse{Success: false, Message: "Mật khẩu xác nhận không khớp"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Hash mật khẩu
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	// Tạo user
	user := User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  hashedPassword,
	}

	// Chèn vào database
	_, err = db.Exec("INSERT INTO users (first_name, last_name, username, email, phone, password) VALUES ($1, $2, $3, $4, $5, $6)",
		user.FirstName, user.LastName, user.Username, user.Email, user.Phone, user.Password)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			resp := RegisterResponse{Success: false, Message: "Tên đăng nhập đã tồn tại"}
			json.NewEncoder(w).Encode(resp)
			return
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			resp := RegisterResponse{Success: false, Message: "Email đã được sử dụng"}
			json.NewEncoder(w).Encode(resp)
			return
		}
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	resp := RegisterResponse{Success: true, Message: "Đăng ký thành công"}
	json.NewEncoder(w).Encode(resp)
}

// Handler xử lý logout
func logoutAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Xóa cookie token
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Đã đăng xuất",
	})
}

// Middleware kiểm tra JWT
func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
