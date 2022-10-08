package main

import (
	"refueling/refueling/pkg/adding"
	"refueling/refueling/pkg/listing"
	"refueling/refueling/pkg/server"
	"refueling/refueling/pkg/storage"

	"github.com/gin-gonic/gin"
)

type Services struct {
	storage *storage.Storage
	adding  adding.AddingService
	listing listing.ListingService
}

func main() {
	engine := setUpEngine()
	engine.Use(CORSMiddleware())
	// engine.Use(cors.Default())

	services := setUpServices()
	srvr := server.NewServer(engine, services.adding, services.listing)
	srvr.Run()

}

func setUpEngine() *gin.Engine {
	engine := gin.Default()
	return engine
}

func setUpServices() *Services {
	storage := storage.NewStorage()
	adding := adding.NewService(storage)
	listing := listing.NewListingService(storage)
	return &Services{storage: storage, adding: adding, listing: listing}
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
