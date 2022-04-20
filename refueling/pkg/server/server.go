package server

import (
	"fmt"
	"os"
	"refueling/refueling/pkg/adding"
	"refueling/refueling/pkg/download"
	"refueling/refueling/pkg/listing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Server struct {
	engine   *gin.Engine
	listing  listing.ListingService
	adding   adding.AddingService
	download download.DownloadService
	client   *redis.Client
	// refs to listing, adding and so on
}

func NewServer(engine *gin.Engine, adding adding.AddingService,
	listing listing.ListingService, download download.DownloadService) *Server {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", os.Getenv("HOST"), "49153")})
	return &Server{engine: engine, adding: adding, listing: listing, download: download, client: client}
}

func (s *Server) Run() {
	router := s.Router()
	err := router.Run(":8888")
	if err != nil {
		panic(err)
	}
}
