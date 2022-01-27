package main

import (
	"fmt"
	"refueling/diary/pkg/adding"
	"refueling/diary/pkg/listing"
	"refueling/diary/pkg/server"
	"refueling/diary/pkg/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hello!")
	// gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Use(CORSMiddleware())
	storage := storage.NewStorage()
	adding := adding.NewService(storage)
	listing := listing.NewListingService(storage)
	srvr := server.NewServer(engine, listing, adding)
	srvr.Run()

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
