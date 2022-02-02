package server

import (
	"fmt"
	"net/http"
	"refueling/refueling/pkg/adding"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int
	Message error
}

func (s *Server) Router() *gin.Engine {
	router := s.engine
	router.GET("/refuelingsList", s.RefNames())
	router.POST("/add", s.Add())
	return router
}

func (s *Server) RefNames() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := s.listing.GetRefuelNames()
		c.IndentedJSON(http.StatusOK, data)
	}
}

func (s *Server) Add() gin.HandlerFunc {
	return func(c *gin.Context) {
		var refuel adding.Refuel
		// var refuel interface{}
		if err := c.BindJSON(&refuel); err != nil {
			panic(err)
		}
		if err := s.adding.Adding(refuel); err != nil {
			errText := fmt.Sprintf("%v", err)
			c.JSON(http.StatusBadRequest, errText)
			return
		}
		c.IndentedJSON(http.StatusOK, "Data recorded")
	}
}
