package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"` // Hashed
}

func initDB() {
	// Kết nối đến PostgreSQL
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password" // Default, change as needed
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "myapp"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Không thể mở kết nối đến PostgreSQL:", err)
	}

	// Kiểm tra kết nối
	err = db.Ping()
	if err != nil {
		log.Fatal("Không thể ping PostgreSQL:", err)
	}

	fmt.Println("Kết nối PostgreSQL thành công")

	// Tạo bảng users nếu chưa tồn tại
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		username VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE,
		phone VARCHAR(20),
		password VARCHAR(255) NOT NULL
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Lỗi tạo bảng:", err)
	}

	fmt.Println("Bảng users đã được tạo hoặc đã tồn tại")
}

// Hash mật khẩu
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Kiểm tra mật khẩu
func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
