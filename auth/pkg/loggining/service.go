package loggining

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type service struct {
	storage Storage
}

type Service interface {
}

type Storage interface {
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

func (s *service) Login(user User) {
	var token TokenDetails
	CreateToken(&token, user.ID)

}

func CreateToken(token *TokenDetails, id uint) error {
	atExp := time.Now().Add(time.Minute * 15).Unix()
	refExp := time.Now().Add(time.Hour * 24 * 1).Unix()
	token.AtExpires = atExp
	token.AccessUuid = uuid.NewV4().String()

	token.RtExpires = refExp
	token.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "q1w2E#R$Q!W@e3r4") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = id
	atClaims["exp"] = token.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims) //sign the token
	token.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	//TODO create func for MapClaims and token registering to avoid hardcoding
}
