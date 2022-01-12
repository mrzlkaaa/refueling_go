package server

import (
	"fmt"
	"log"
	"net/http"
	"refueling/pkg/adding"

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
		var detailWeek []adding.DetailWeek
		err := c.BindJSON(&detailWeek)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(detailWeek)
	}
}
