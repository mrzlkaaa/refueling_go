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
	storage := storage.NewStorage() //*gorm pointer
	adding := adding.NewService(storage)
	loggining := loggining.NewService(storage)
	listing := listing.NewService(storage)
	server := server.NewServer(engine, adding,
		loggining, listing)
	server.Run()
	//todo services are add, list, storage
}
