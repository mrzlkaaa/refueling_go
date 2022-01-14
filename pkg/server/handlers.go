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
	router.GET("/refuelingsList", s.RefNames())
	router.GET("/getNewWeekNum/:fcName", s.NewWeekNum())
	router.POST("/SubmitWeekData", s.SubmitWeekData())
	return router
}

func (s *Server) RefNames() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := s.listing.GetRefuelNames()
		c.IndentedJSON(http.StatusOK, data)
	}
}

func (s *Server) NewWeekNum() gin.HandlerFunc {
	return func(c *gin.Context) {
		fcName := c.Param("fcName")
		s.listing.GetNewWeekNum(fcName)

	}
}

func (s *Server) SubmitWeekData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formsData adding.FormsData
		err := c.BindJSON(&formsData)
		fmt.Println("i got shit like", err)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(formsData)
		s.adding.AddWeeklyData(&formsData)
	}
}
