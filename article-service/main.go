package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("Set DSN env var, contoh: export DSN=\"root:secret@tcp(127.0.0.1:3306)/article?parseTime=true\"")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed connect db: %v", err)
	}

	// Auto-migrate (migrasi otomatis)
	if err := db.AutoMigrate(&Post{}); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	r := gin.Default()

	// Routes
	r.POST("/article/", createArticle(db))
	r.GET("/article/:limit/:offset", listArticles(db))
	r.GET("/article/:id", getArticle(db))
	r.PUT("/article/:id", updateArticle(db))    // update
	r.PATCH("/article/:id", updateArticle(db))  // also accept PATCH
	r.DELETE("/article/:id", deleteArticle(db))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
