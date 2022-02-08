package server

import (
	"fmt"
	"net/http"
	"refueling/refueling/pkg/adding"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	statusOkAdd    = "Data recorded"
	statusOkDelete = "Data deleted"
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
	router.GET("/refuelings/:id/:actId/PDC", s.RefuelPDC())
	router.POST("/add", s.Add())
	router.POST("/add-act", s.AddAct())
	router.POST("/refuelings/:id/:actId/delete", s.DeleteAct())
	router.POST("/refuelings/:id/delete", s.Delete())
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
		idParam := c.Param("actId")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			panic(err)
		}
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
		c.IndentedJSON(http.StatusOK, fmt.Sprintf("%v. You will be redirect in few seconds", statusOkAdd))
	}
}

func (s *Server) AddAct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var act adding.Act
		// var refuel interface{}
		if err := c.BindJSON(&act); err != nil {
			panic(err)
		}
		err, id := s.adding.AddingAct(act)
		var obj map[string]interface{} = map[string]interface{}{"msg": statusOkAdd, "id": id}
		if err != nil {
			errText := fmt.Sprintf("%v", err)
			obj["msg"] = errText
			c.JSON(http.StatusBadRequest, obj)
		}
		c.IndentedJSON(http.StatusOK, obj)
	}
}

func (s *Server) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var refuel interface{}
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			panic(err)
		}
		err = s.adding.Deleting(id)
		if err != nil {
			errText := fmt.Sprintf("%v", err)
			fmt.Println(errText)
			c.JSON(http.StatusBadRequest, errText)
			return
		}
		c.IndentedJSON(http.StatusOK, statusOkDelete)
	}
}

func (s *Server) DeleteAct() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var refuel interface{}
		idParam := c.Param("actId")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			panic(err)
		}
		err = s.adding.DeletingAct(id)
		if err != nil {
			errText := fmt.Sprintf("%v", err)
			fmt.Println(errText)
			c.JSON(http.StatusBadRequest, errText)
			return
		}
		c.IndentedJSON(http.StatusOK, statusOkDelete)
	}
}
