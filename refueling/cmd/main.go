package main

import (
	"refueling/refueling/pkg/adding"
	"refueling/refueling/pkg/download"
	"refueling/refueling/pkg/listing"
	"refueling/refueling/pkg/server"
	"refueling/refueling/pkg/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.Use(CORSMiddleware())
	// engine.Use(cors.Default())
	storage := storage.NewStorage()
	adding := adding.NewService(storage)
	listing := listing.NewListingService(storage)
	download := download.NewService(storage)
	srvr := server.NewServer(engine, adding, listing, download)
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
