package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

func createArticle(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUpdatePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json", "detail": err.Error()})
			return
		}
		if err := validate.Struct(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "detail": err.Error()})
			return
		}
		post := Post{
			Title:    req.Title,
			Content:  req.Content,
			Category: req.Category,
			Status:   req.Status,
		}
		if err := db.Create(&post).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed create", "detail": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": post.ID})
	}
}

func listArticles(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.Param("limit")
		offsetStr := c.Param("offset")
		limit, err1 := strconv.Atoi(limitStr)
		offset, err2 := strconv.Atoi(offsetStr)
		if err1 != nil || err2 != nil || limit < 0 || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit and offset must be non-negative integers"})
			return
		}
		var posts []Post
		if err := db.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error", "detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, posts)
	}
}

func getArticle(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var post Post
		if err := db.First(&post, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "db error", "detail": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, post)
	}
}

func updateArticle(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var req CreateUpdatePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json", "detail": err.Error()})
			return
		}
		if err := validate.Struct(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "detail": err.Error()})
			return
		}
		var post Post
		if err := db.First(&post, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "db error", "detail": err.Error()})
			}
			return
		}
		post.Title = req.Title
		post.Content = req.Content
		post.Category = req.Category
		post.Status = req.Status
		if err := db.Save(&post).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed update", "detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}

func deleteArticle(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if err := db.Delete(&Post{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed delete", "detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}
