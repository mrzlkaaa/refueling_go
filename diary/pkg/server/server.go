package server

import (
	"refueling/diary/pkg/adding"
	"refueling/diary/pkg/listing"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine  *gin.Engine
	listing listing.ListingService
	adding  adding.AddingService
	// refs to listing, adding and so on
}

func NewServer(engine *gin.Engine, listing listing.ListingService,
	adding adding.AddingService) *Server {
	return &Server{engine: engine, listing: listing, adding: adding}
}

func (s *Server) Run() {
	router := s.Router()
	err := router.Run(":8888")
	if err != nil {
		panic(err)
	}
}
