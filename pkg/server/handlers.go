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
	router.POST("/submitWeekData", s.SubmitWeekData())
	return router
}

func (s *Server) RefNames() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := s.listingR.GetRefuelNames()
		c.IndentedJSON(http.StatusOK, data)
	}
}

func (s *Server) NewWeekNum() gin.HandlerFunc {
	return func(c *gin.Context) {
		fcName := c.Param("fcName")
		weekNum := s.listingD.GetNewWeekNum(fcName)
		c.IndentedJSON(http.StatusOK, weekNum)

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
