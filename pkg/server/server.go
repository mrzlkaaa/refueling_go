package server

import (
	"refueling/pkg/listing"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine  *gin.Engine
	listing listing.ListingService
	// refs to listing, adding and so on
}

func NewServer(engine *gin.Engine, listing listing.ListingService) *Server {
	return &Server{engine: engine, listing: listing}
}

func (s *Server) Run() {
	router := s.Router()
	err := router.Run(":8081")
	if err != nil {
		panic(err)
	}
}
