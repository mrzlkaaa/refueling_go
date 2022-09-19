package server

import (
	"fmt"
	"os"
	"refueling/refueling/pkg/adding"
	"refueling/refueling/pkg/listing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Server struct {
	engine  *gin.Engine
	listing listing.ListingService
	adding  adding.AddingService
	// download download.DownloadService
	client *redis.Client
	// refs to listing, adding and so on
}

func NewServer(engine *gin.Engine, adding adding.AddingService,
	listing listing.ListingService) *Server {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", os.Getenv("HOST"), "49153")})
	return &Server{engine: engine, adding: adding, listing: listing, client: client}
}

func (s *Server) Run() {
	router := s.Router()
	err := router.Run(":1111") //* 8888 - prod, 1111 - dev
	if err != nil {
		panic(err)
	}
}
