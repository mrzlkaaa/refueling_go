package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Router() *gin.Engine {
	router := s.engine
	router.GET("/", s.test())
	return router
}

func (s *Server) test() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := s.listing.GetRefuelNames()
		c.IndentedJSON(http.StatusOK, data)
	}
}
