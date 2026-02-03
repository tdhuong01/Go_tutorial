package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var jwtKey = []byte("your_secret_key") // Đổi thành key bí mật của bạn

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
	Username string `json:"username"`
	Password string `json:"password"`
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
	http.ServeFile(w, r, "login.html")
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
	err = userCollection.FindOne(context.TODO(), bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
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

	if req.Username == "" || req.Password == "" {
		resp := RegisterResponse{Success: false, Message: "Tên đăng nhập và mật khẩu không được để trống"}
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
		Username: req.Username,
		Password: hashedPassword,
	}

	// Chèn vào database
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			resp := RegisterResponse{Success: false, Message: "Tên đăng nhập đã tồn tại"}
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
