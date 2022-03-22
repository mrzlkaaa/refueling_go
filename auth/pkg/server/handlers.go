package server

import (
	"net/http"
	"refueling/auth/pkg/adding"

	"github.com/gin-gonic/gin"
)

func (s *Server) Router() *gin.Engine {
	router := s.engine
	router.POST("/add", s.AddUser())
	return router
}

func (s *Server) AddUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userForm *adding.User

		err := c.BindJSON(&userForm)
		if err != nil {
			panic(err)
		}

		err = s.adding.AddUser(&userForm)
		if err != nil {
			c.JSON(500, err)
		}
		c.JSON(http.StatusOK, "User successfully registered")
	}
}
