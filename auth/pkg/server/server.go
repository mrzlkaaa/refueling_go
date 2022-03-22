package server

import (
	"log"
	"refueling/auth/pkg/adding"
	"refueling/auth/pkg/listing"
	"refueling/auth/pkg/loggining"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine    *gin.Engine
	adding    adding.Service
	loggining loggining.Service
	listing   listing.Service
}

func NewServer(engine *gin.Engine, adding adding.Service, loggining loggining.Service, listing listing.Service) *Server {

	return &Server{engine: engine, adding: adding, loggining: loggining, listing: listing}
}

func (s *Server) Run() {
	router := s.Router()
	err := router.Run(":8892")
	if err != nil {
		log.Fatal(err)
	}
}
