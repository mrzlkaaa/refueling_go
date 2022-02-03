package server

import (
	"fmt"
	"net/http"
	"refueling/refueling/pkg/adding"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int
	Message error
}

func (s *Server) Router() *gin.Engine {
	router := s.engine
	router.GET("/refuelingsList", s.RefNames())
	router.GET("/refuelings", s.Refuels())
	router.GET("/refuelings/:id/details", s.RefuelDetails())
	router.GET("/refuelings/:id/PDC", s.RefuelPDC())
	router.POST("/add", s.Add())
	return router
}

func (s *Server) RefNames() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := s.listing.GetRefuelNames()
		c.IndentedJSON(http.StatusOK, data)
	}
}

func (s *Server) Refuels() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := s.listing.Refuels()
		c.IndentedJSON(http.StatusOK, data)
	}
}

func (s *Server) RefuelDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, _ := strconv.Atoi(idstr)
		data := s.listing.RefuelDetails(id)
		c.IndentedJSON(http.StatusOK, data)
	}
}

func (s *Server) RefuelPDC() gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, _ := strconv.Atoi(idstr)
		data := s.listing.RefuelPDC(id)
		c.JSON(http.StatusOK, data)
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
