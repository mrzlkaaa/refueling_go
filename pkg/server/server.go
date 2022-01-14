package server

import (
	"refueling/pkg/adding"
	listingDiary "refueling/pkg/listing/diary"
	listingRefuels "refueling/pkg/listing/refuels"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine   *gin.Engine
	listingR listingRefuels.ListingService
	listingD listingDiary.ListingService
	adding   adding.AddingService
	// refs to listing, adding and so on
}

func NewServer(engine *gin.Engine, listingR listingRefuels.ListingService,
	listingD listingDiary.ListingService, adding adding.AddingService) *Server {
	return &Server{engine: engine, listingR: listingR, listingD: listingD, adding: adding}
}

func (s *Server) Run() {
	router := s.Router()
	err := router.Run(":8888")
	if err != nil {
		panic(err)
	}
}
