package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"refueling/refueling/pkg/adding"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	statusOkAdd    = "Data recorded"
	statusOkDelete = "Data deleted"
	notSerialized  = "given JSON is not serializable"
	unauthorized   = "unauthorized"
	noRigths       = "you are not eligible for this service"
	expired        = "token is expired"
	unexpected     = "something unexpected happened"
)

type Response struct {
	Status  int
	Message error
}

func (s *Server) Router() *gin.Engine {
	router := s.engine
	// router.GET("/refuelings", s.Refuels())
	auth := router.Group("/")
	auth.Use(s.AuthRequired())
	{
		auth.GET("/refuelingsList", s.RefNames())
		auth.GET("/refuelings", s.Refuels())
		auth.GET("/refuelings/:id/details", s.RefuelDetails())
		auth.GET("/refuelings/:id/:actId/PDC", s.RefuelPDC())
		auth.GET("/refuelings/:id/:actId/download", s.Download())
		auth.POST("/add", s.Add())
		auth.POST("/add-act", s.AddAct())
		auth.POST("/refuelings/:id/:actId/delete", s.DeleteAct())
		auth.POST("/refuelings/:id/delete", s.Delete())
	}

	return router
}

func (s *Server) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := ValidateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		rights, err := s.FetchValue(claims["access_uuid"].(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, expired)
			c.Abort()
			return
		}
		moderator, admin, err := ParseRights(rights)
		if err != nil {
			c.JSON(http.StatusUnauthorized, noRigths)
			c.Abort()
			return
		}

		c.Set("access_uuid", claims["access_uuid"])
		c.Set("moderator", moderator)
		c.Set("admin", admin)
		c.Next()
	}
}

func FetchAuth(c *gin.Context) (map[string]interface{}, error) {
	claims, err := ValidateToken(c)
	if err != nil {
		// c.JSON(http.StatusUnauthorized, err)
		return map[string]interface{}{}, err
	}

	return claims, nil
}

func ValidateToken(c *gin.Context) (map[string]interface{}, error) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return map[string]interface{}{}, errors.New("bearer token is not given")
	}
	tokenString := strings.Split(auth, " ")[len(strings.Split(auth, " "))-1]
	token, err := VerifyToken(tokenString)
	fmt.Println(token, err)
	if err != nil {
		return map[string]interface{}{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid && !ok {
		return map[string]interface{}{}, errors.New("is not valid")
	}
	// uuid := claims["access_uuid"].(string)

	return claims, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) { //todo add extra argument to seperate refresh and access verifications
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil

	})
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if ok && ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return token, errors.New("signature is invalid")
		} else if ok && ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return token, errors.New(expired)
		} else {
			return token, errors.New(unexpected)
		}
	}
	return token, nil
}

func (s *Server) FetchValue(key string) (string, error) {
	v, err := s.client.Get(key).Result()
	return v, err
}

func ParseRights(rights string) (bool, bool, error) {
	sliceRights := strings.Split(rights, ",")
	moderatorStr, adminStr := sliceRights[0], sliceRights[1]
	moderator, err := strconv.ParseBool(moderatorStr)
	if err != nil {
		return false, false, err
	}
	admin, _ := strconv.ParseBool(adminStr)
	if err != nil {
		return false, false, err
	}

	return moderator, admin, nil
}

func (s *Server) RefNames() gin.HandlerFunc {
	return func(c *gin.Context) {
		moderator, ok := c.Get("moderator")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}

		admin, ok := c.Get("admin")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}
		if !moderator.(bool) && !admin.(bool) {
			c.JSON(http.StatusUnauthorized, noRigths)
			return
		}

		data := s.listing.GetRefuelNames()
		c.IndentedJSON(http.StatusOK, data)
	}
}

func (s *Server) Refuels() gin.HandlerFunc {
	return func(c *gin.Context) {
		moderator, ok := c.Get("moderator")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}

		admin, ok := c.Get("admin")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}
		if !moderator.(bool) && !admin.(bool) {
			c.JSON(http.StatusUnauthorized, noRigths)
			return
		}

		data := s.listing.Refuels()
		c.IndentedJSON(http.StatusOK, data)
	}
}

func (s *Server) RefuelDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		moderator, ok := c.Get("moderator")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}

		admin, ok := c.Get("admin")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}
		if !moderator.(bool) && !admin.(bool) {
			c.JSON(http.StatusUnauthorized, noRigths)
			return
		}

		idstr := c.Param("id")
		id, _ := strconv.Atoi(idstr)
		data := s.listing.RefuelDetails(id)
		c.IndentedJSON(http.StatusOK, data)
	}
}

func (s *Server) RefuelPDC() gin.HandlerFunc {
	return func(c *gin.Context) {
		moderator, ok := c.Get("moderator")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}

		admin, ok := c.Get("admin")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}
		if !moderator.(bool) && !admin.(bool) {
			c.JSON(http.StatusUnauthorized, noRigths)
			return
		}

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
		admin, ok := c.Get("admin")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}
		if !admin.(bool) {
			c.JSON(http.StatusUnauthorized, noRigths)
			return
		}

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
		admin, ok := c.Get("admin")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}
		if !admin.(bool) {
			c.JSON(http.StatusUnauthorized, noRigths)
			return
		}

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
		admin, ok := c.Get("admin")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}
		if !admin.(bool) {
			c.JSON(http.StatusUnauthorized, noRigths)
			return
		}

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
		admin, ok := c.Get("admin")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}
		if !admin.(bool) {
			c.JSON(http.StatusUnauthorized, noRigths)
			return
		}

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

func (s *Server) Download() gin.HandlerFunc {
	return func(c *gin.Context) {
		moderator, ok := c.Get("moderator")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}

		admin, ok := c.Get("admin")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}
		if !moderator.(bool) && !admin.(bool) {
			c.JSON(http.StatusUnauthorized, noRigths)
			return
		}

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			panic(err)
		}
		actIdParam := c.Param("actId")
		actId, err := strconv.Atoi(actIdParam)
		if err != nil {
			panic(err)
		}
		path := s.download.SavePDC(id, actId)
		c.IndentedJSON(http.StatusOK, path)
	}
}
