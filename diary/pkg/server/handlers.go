package server

import (
	"fmt"
	"log"
	"net/http"
	"refueling/diary/pkg/adding"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) Router() *gin.Engine {
	router := s.engine
	router.GET("/WeeksNum/:fcName", s.WeeksNum())
	router.GET("/WeekDetails/:fcName/:weekName", s.WeeekDetails())
	router.POST("/submitWeekData", s.SubmitWeekData())
	return router
}

func (s *Server) WeeksNum() gin.HandlerFunc {
	return func(c *gin.Context) {
		fcName := c.Param("fcName")
		weekNum := s.listing.GetWeeksNum(fcName)
		c.IndentedJSON(http.StatusOK, weekNum)

	}
}

func (s *Server) WeeekDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		fcName := c.Param("fcName")
		weekName, _ := strconv.Atoi(c.Param("weekName"))
		fmt.Println(fcName, weekName)
		object := s.listing.WeekDetails(fcName, weekName)
		c.IndentedJSON(http.StatusOK, object)
	}
}

func (s *Server) SubmitWeekData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formsData adding.FuelCycle
		err := c.BindJSON(&formsData)
		fmt.Println("i got shit like", err)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(formsData)
		s.adding.AddWeeklyData(&formsData)
	}
}
