package server

import (
	"refueling/refueling/pkg/adding"
	"refueling/refueling/pkg/listing"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine  *gin.Engine
	listing listing.ListingService
	adding  adding.AddingService
	// refs to listing, adding and so on
}

func NewServer(engine *gin.Engine, adding adding.AddingService, listing listing.ListingService) *Server {
	return &Server{engine: engine, adding: adding, listing: listing}
}

func (s *Server) Run() {
	router := s.Router()
	err := router.Run(":8888")
	if err != nil {
		panic(err)
	}
}
