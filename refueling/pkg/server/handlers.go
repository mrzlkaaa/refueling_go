package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Router() *gin.Engine {
	router := s.engine
	router.GET("/refuelingsList", s.RefNames())
	return router
}

func (s *Server) RefNames() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := s.listing.GetRefuelNames()
		c.IndentedJSON(http.StatusOK, data)
	}
}
