package main

import (
	"fmt"
	"refueling/pkg/listing"
	"refueling/pkg/server"
	"refueling/pkg/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hello!")
	// gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Use(CORSMiddleware())
	strge := storage.NewStorage()
	listing := listing.NewListingService(strge)
	srvr := server.NewServer(engine, listing)
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
