package main

import (
	"refueling/auth/pkg/adding"
	"refueling/auth/pkg/listing"
	"refueling/auth/pkg/loggining"
	"refueling/auth/pkg/server"
	"refueling/auth/pkg/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.Use(CORSMiddleware())
	storage := storage.NewStorage() //*gorm pointer
	adding := adding.NewService(storage)
	loggining := loggining.NewService(storage)
	listing := listing.NewService(storage)
	server := server.NewServer(engine, adding,
		loggining, listing)
	server.Run()
	//todo services are add, list, storage
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
