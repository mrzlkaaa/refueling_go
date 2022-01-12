package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Router() *gin.Engine {
	router := s.engine
	router.GET("/", s.RefNames())
	router.POST("/post", s.Post())
	return router
}

func (s *Server) RefNames() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := s.listing.GetRefuelNames()
		c.IndentedJSON(http.StatusOK, data)
	}
}

func (s *Server) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		msg := c.PostForm("title")
		c.IndentedJSON(http.StatusOK, msg)
	}
}
