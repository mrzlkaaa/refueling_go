package main

import (
	"fmt"
	"refueling/pkg/adding"
	"refueling/pkg/listing"
	"refueling/pkg/server"
	"refueling/pkg/storage/NoSQL"
	"refueling/pkg/storage/SQL"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hello!")
	// gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Use(CORSMiddleware())
	strgNoSQL := NoSQL.NewStorage()
	strgeSQL := SQL.NewStorage()
	adding := adding.NewService(strgNoSQL)
	listingR := listing.NewListingService(strgeSQL)
	srvr := server.NewServer(engine, listingR, adding)
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
