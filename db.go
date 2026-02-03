package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var client *mongo.Client
var userCollection *mongo.Collection

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"` // Hashed
}

func initDB() {
	// Kết nối đến MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Default local MongoDB
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error
	client, err = mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Không thể kết nối đến MongoDB:", err)
	}

	// Kiểm tra kết nối
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Không thể ping MongoDB:", err)
	}

	fmt.Println("Kết nối MongoDB thành công")

	// Chọn database và collection
	database := client.Database("myapp")
	userCollection = database.Collection("users")

	// Tạo index cho username (unique)
	indexModel := mongo.IndexModel{
		Keys:    map[string]int{"username": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err = userCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Println("Lỗi tạo index:", err)
	}
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
