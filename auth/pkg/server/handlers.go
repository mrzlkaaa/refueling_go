package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"refueling/auth/pkg/adding"
	"refueling/auth/pkg/loggining"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	notSerialized = "given JSON is not serializable"
	unauthorized  = "unauthorized"
	noRigths      = "you are not eligible for this service"
	expired       = "token is expired"
)

func (s *Server) Router() *gin.Engine {
	router := s.engine
	auth := router.Group("/")
	auth.Use(AuthRequired())
	{
		auth.POST("/add", s.AddUser())
	}
	router.POST("/login", s.Login())
	router.POST("/refreshToken", s.RefreshToken())
	return router
}

func (s *Server) AddUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := FetchAuth(c)
		rights, err := s.loggining.FetchValue(claims["access_uuid"].(string)) //* may be should ommit here and use only to logOut
		if err != nil {
			c.JSON(http.StatusUnauthorized, expired)
			return
		}
		//!checking for rigths of user
		err = ParseRights(rights)
		if err != nil {
			c.JSON(http.StatusUnauthorized, noRigths)
			return
		}

		var userForm adding.User
		err = c.BindJSON(&userForm)
		if err != nil {
			c.JSON(http.StatusBadRequest, notSerialized)
			return
		}

		err = s.adding.AddUser(userForm)
		if err != nil {
			errText := fmt.Sprintf("%v", err)
			c.JSON(http.StatusBadRequest, errText)
			return
		}
		c.JSON(http.StatusOK, "User successfully registered")
	}
}

func (s *Server) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user loggining.User

		err := c.BindJSON(&user)
		if err != nil {
			panic(err)
		}
		token, err := s.loggining.Login(user)
		if err != nil {
			errText := fmt.Sprintf("%v", err)
			c.IndentedJSON(http.StatusOK, errText)
			return
		}
		c.IndentedJSON(http.StatusOK, token)

	}
}

func (s *Server) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		mapToken := map[string]string{}
		err := c.BindJSON(&mapToken) //* also mapping username
		if err != nil || len(mapToken) == 0 {
			c.JSON(http.StatusUnprocessableEntity, notSerialized)
			return
		}
		token, err := VerifyToken(mapToken["token"], false)
		if err != nil {
			c.JSON(http.StatusUnauthorized, unauthorized)
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !token.Valid && !ok {
			c.JSON(http.StatusUnauthorized, errors.New("token is not valid"))
			return
		}
		id, err := s.loggining.FetchValue(claims["refresh_uuid"].(string)) //* may be should ommit here and use only to logOut
		if err != nil {
			c.JSON(http.StatusUnauthorized, expired) //* relogging required
			return
		}
		err = s.loggining.DeleteToken(claims["refresh_uuid"].(string))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "something unexpectable happened")
		}
		idInt, _ := strconv.Atoi(id)
		idUint := uint(idInt)
		newToken, err := s.loggining.RefreshToken(idUint)
		if err != nil {
			errText := fmt.Sprintf("%v", err)
			c.IndentedJSON(http.StatusOK, errText)
			return
		}
		c.IndentedJSON(http.StatusOK, newToken)

	}
}

func FetchAuth(c *gin.Context) map[string]interface{} {
	claims, err := ValidateToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, unauthorized)
		return map[string]interface{}{}
	}

	return claims
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := ValidateToken(c); err != nil {
			//! what if token expired? need to call refresh function
			c.JSON(http.StatusUnauthorized, expired)
			c.Abort()
			return
		}
		c.Next()
	}
}

func ValidateToken(c *gin.Context) (map[string]interface{}, error) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return map[string]interface{}{}, errors.New(unauthorized)
	}
	tokenString := strings.Split(auth, " ")[len(strings.Split(auth, " "))-1]
	token, err := VerifyToken(tokenString, true)
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

func VerifyToken(tokenString string, access bool) (*jwt.Token, error) { //todo add extra argument to seperate refresh and access verifications
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if access {
			return []byte(os.Getenv("ACCESS_SECRET")), nil
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil

	})
	return token, err
}

func ParseRights(rights string) error {
	sliceRights := strings.Split(rights, ",")
	moderatorStr, adminStr := sliceRights[0], sliceRights[1]
	moderator, _ := strconv.ParseBool(moderatorStr)
	admin, _ := strconv.ParseBool(adminStr)
	if moderator || admin {
		return nil
	}

	return errors.New(noRigths)
}
