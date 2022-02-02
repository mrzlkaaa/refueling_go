package main

import (
	"refueling/refueling/pkg/adding"
	"refueling/refueling/pkg/listing"
	"refueling/refueling/pkg/server"
	"refueling/refueling/pkg/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	// engine.Use(CORSMiddleware())
	engine.Use(cors.Default())
	storage := storage.NewStorage()
	adding := adding.NewService(storage)
	listing := listing.NewListingService(storage)
	srvr := server.NewServer(engine, adding, listing)
	srvr.Run()

}

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }
