package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"refueling/auth/pkg/adding"
	"refueling/auth/pkg/listing"
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
	unexpected    = "something unexpected happened"
)

func (s *Server) Router() *gin.Engine {
	router := s.engine
	auth := router.Group("/")
	auth.Use(s.AuthRequired())
	{
		auth.POST("/logout", s.Logout())
		auth.POST("/users/:id/delete", s.DeleteUser()) //! admin rights
		auth.POST("/users/update", s.UpdateUsers())    //! admin rights
		auth.GET("/getusers", s.GetAllUsers())         //! admin rights
	}
	router.POST("/register", s.AddUser())
	router.POST("/login", s.Login())
	router.POST("/refreshToken", s.RefreshToken())
	return router
}

func (s *Server) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := ValidateToken(c)
		if err != nil {
			// c.JSON(http.StatusUnauthorized, err.Error())
			c.AbortWithError(http.StatusUnauthorized, err)
		}
		rights, err := s.loggining.FetchValue(claims["access_uuid"].(string))
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}
		moderator, admin, err := ParseRights(rights)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
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
		return map[string]interface{}{}, err
	}

	return claims, nil
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
			c.IndentedJSON(http.StatusUnauthorized, errText)
			return
		}
		c.IndentedJSON(http.StatusOK, token)

	}
}

func (s *Server) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessuuid, ok := c.Get("access_uuid")
		if !ok {
			c.JSON(http.StatusUnauthorized, unauthorized)
		}

		err := s.loggining.DeleteToken(accessuuid.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, unexpected)
			return
		}
		c.JSON(http.StatusOK, "logout successfully")
	}
}

func (s *Server) AddUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userForm adding.User
		err := c.BindJSON(&userForm)
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
		c.JSON(http.StatusOK, "user successfully registered")
	}
}

func (s *Server) UpdateUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		admin, ok := c.Get("admin")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}
		if !admin.(bool) {
			c.JSON(http.StatusUnauthorized, noRigths)
			return
		}

		var users []listing.User
		err := c.BindJSON(&users)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err = s.adding.UpdateUsers(&users)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		c.JSON(http.StatusOK, "successfully updated")
	}
}

func (s *Server) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		admin, ok := c.Get("admin")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}
		if !admin.(bool) {
			c.JSON(http.StatusUnauthorized, noRigths)
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, unexpected)
			return
		}
		uid := uint(id)
		err = s.adding.DeleteUser(uid)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, unexpected)
			return
		}
		c.JSON(http.StatusOK, "deleted successfully")
	}
}

func (s *Server) GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		admin, ok := c.Get("admin")
		if !ok {
			c.JSON(http.StatusUnauthorized, noRigths)
		}
		if !admin.(bool) {
			c.JSON(http.StatusUnauthorized, noRigths)
			return
		}

		users, err := s.listing.GetAllUsers()
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, users)
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
			c.JSON(http.StatusUnauthorized, err.Error())
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !token.Valid && !ok {
			c.JSON(http.StatusUnauthorized, "token is not valid")
			return
		}
		id, err := s.loggining.FetchValue(claims["refresh_uuid"].(string)) //* may be should ommit here and use only to logOut
		if err != nil {
			c.JSON(http.StatusUnauthorized, expired) //* relogging required
			return
		}
		err = s.loggining.DeleteToken(claims["refresh_uuid"].(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, unexpected)
			return
		}
		idInt, _ := strconv.Atoi(id)
		idUint := uint(idInt)
		newToken, err := s.loggining.RefreshToken(idUint)
		if err != nil {
			c.IndentedJSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		c.IndentedJSON(http.StatusOK, newToken)

	}
}

func ValidateToken(c *gin.Context) (map[string]interface{}, error) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return map[string]interface{}{}, errors.New("bearer token is not given")
	}
	tokenString := strings.Split(auth, " ")[len(strings.Split(auth, " "))-1]
	token, err := VerifyToken(tokenString, true)
	if err != nil {
		return map[string]interface{}{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid && !ok {
		return map[string]interface{}{}, errors.New("token is not valid")
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
