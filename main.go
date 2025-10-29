package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	uploadDirPath string
	secretKey     string
)

func init() {
	uploadDirPath = getEnv("UPLOAD_DIR_PATH", "uploads")
	secretKey = getEnv("SECRET_KEY", "")

	if secretKey == "" {
		panic("SECRET_KEY environment variable is required")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getMimeType(filename string) string {
	ext := filepath.Ext(filename)
	return mime.TypeByExtension(ext)
}

func validateUrl(c *gin.Context) bool {
	filename := c.Param("filename")
	expireStr := c.Query("expires")
	signature := c.Query("signature")

	if expireStr == "" || signature == "" {
		return false
	}

	expires, err := strconv.ParseInt(expireStr, 10, 64)
	if err != nil || time.Now().Unix() > expires {
		return false
	}

	method := c.Request.Method

	// For POST requests without filename, use empty string
	data := fmt.Sprintf("%s:%s:%d", method, filename, expires)
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	expectedsignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedsignature))
}

func SignedURLMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !validateUrl(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid or expired URL"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Server is running"})
	})

	router.GET("/images/:filename", SignedURLMiddleware(), func(c *gin.Context) {
		filename := c.Param("filename")
		path := filepath.Join(uploadDirPath, filename)

		file, err := os.Open(path)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "File not found"})
			return
		}
		defer file.Close()

		c.Header("Content-Disposition", "inline; filename="+filename)
		c.Header("Content-Type", getMimeType(filename))
		c.File(path)
	})

	router.POST("/images", SignedURLMiddleware(), func(c *gin.Context) {
		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "File not found in the request"})
			return
		}
		defer file.Close()

		if _, err := os.Stat(uploadDirPath); os.IsNotExist(err) {
			os.MkdirAll(uploadDirPath, 0755)
		}

		extentionName := filepath.Ext(fileHeader.Filename)
		newFileName := uuid.New().String() + extentionName

		destinationPath := filepath.Join(uploadDirPath, newFileName)
		destinationFile, err := os.Create(destinationPath)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create file."})
			return
		}
		defer destinationFile.Close()

		if _, err := io.Copy(destinationFile, file); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file."})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{
			"message":           "File uploaded",
			"filename":          newFileName,
			"original_filename": fileHeader.Filename,
			"size":              fileHeader.Size,
		})
	})

	router.PUT("/images/:filename", SignedURLMiddleware(), func(c *gin.Context) {
		path := filepath.Join(uploadDirPath, c.Param("filename"))

		if _, err := os.Stat(path); os.IsNotExist(err) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "File Not found."})
			return
		}

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "File not found in the request"})
			return
		}
		defer file.Close()

		existingFile, err := os.Create(path)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create file."})
			return
		}
		defer existingFile.Close()

		if _, err := io.Copy(existingFile, file); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file."})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "File updated",
			"size":    header.Size,
		})
	})

	router.DELETE("/images/:filename", SignedURLMiddleware(), func(c *gin.Context) {
		path := filepath.Join(uploadDirPath, c.Param("filename"))

		if err := os.Remove(path); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to remove file."})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "File removed"})
	})

	port := getEnv("SERVER_PORT", ":8000")
	if port[0] != ':' {
		port = ":" + port
	}
	router.Run(port)
}
