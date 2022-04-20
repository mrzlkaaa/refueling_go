package loggining

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	wrongCredentials = "Wrong credentials"
	persistFailed    = "Persistenting credentials failed"
)

type service struct {
	storage Storage
	client  *redis.Client
}

type Service interface {
	Login(User) (map[string]string, error)
	FetchValue(string) (string, error)
	RefreshToken(uint) (map[string]string, error)
	DeleteToken(string) error
}

type Storage interface {
	FindUser(string) (UserData, error)
	FindUserID(uint) (UserData, error)
}

func NewService(storage Storage) Service {
	// redisRun()
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", os.Getenv("HOST"), "49153")})
	return &service{storage: storage, client: client}
}

// var Client *redis.Client

// func redisRun() {
// 	Client = redis.NewClient(&redis.Options{
// 		Addr: fmt.Sprintf("%v:%v", os.Getenv("HOST"), "49153")})
// }

func (s *service) Login(user User) (map[string]string, error) {
	var userData UserData
	token := map[string]string{}
	userData, err := s.storage.FindUser(user.Username)
	if err != nil {
		return token, err
	}

	if ok := CheckPass(userData.PswdHash, user.Password); userData.Username != user.Username || !ok {
		return token, errors.New(wrongCredentials)
	}

	var td TokenDetails
	CreateRefreshToken(&td, userData.Name, userData.Surname, userData.Email)
	CreateAccessToken(&td, userData.ID, userData.Moderator, userData.Admin)

	err = s.PersistTokenDetails(&td, userData.Moderator, userData.Admin, userData.ID)
	if err != nil {
		fmt.Println(err)
		return token, errors.New(persistFailed)
	}

	token["accessToken"] = td.AccessToken
	token["refreshToken"] = td.RefreshToken

	return token, err
}

func (s *service) FetchValue(key string) (string, error) {
	v, err := s.client.Get(key).Result()
	return v, err
}

func (s *service) PersistTokenDetails(td *TokenDetails, moderator, admin bool, id uint) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	accValue := fmt.Sprintf("%v,%v", moderator, admin)

	errAccess := s.client.Set(td.AccessUuid, accValue, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := s.client.Set(td.RefreshUuid, id, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errAccess
	}

	return nil
}

func (s *service) RefreshToken(userID uint) (map[string]string, error) {
	var userData UserData
	token := map[string]string{}
	userData, err := s.storage.FindUserID(userID)
	if err != nil {
		return token, err
	}
	var td TokenDetails
	CreateRefreshToken(&td, userData.Name, userData.Surname, userData.Email)
	CreateAccessToken(&td, userData.ID, userData.Moderator, userData.Admin)

	err = s.PersistTokenDetails(&td, userData.Moderator, userData.Admin, userData.ID)
	if err != nil {
		fmt.Println(err)
		return token, errors.New(persistFailed)
	}

	token["accessToken"] = td.AccessToken
	token["refreshToken"] = td.RefreshToken
	return token, nil
}

func (s *service) DeleteToken(key string) error {
	_, err := s.client.Del(key).Result()
	return err
}

func CheckPass(hash []byte, pswd string) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(pswd))
	return err == nil
}

func CreateRefreshToken(token *TokenDetails, name, surname, email string) {
	refExp := time.Now().Add(time.Hour * 24 * 2).Unix()

	token.RtExpires = refExp
	token.RefreshUuid = uuid.NewString()

	var err error

	//*identfy who the user is
	rtClaims := jwt.MapClaims{}
	rtClaims["name"] = name
	rtClaims["surname"] = surname
	rtClaims["email"] = email
	rtClaims["refresh_uuid"] = token.RefreshUuid
	rtClaims["exp"] = token.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims) //sign the token
	token.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		panic(err)
	}
}

func CreateAccessToken(token *TokenDetails, id uint, moderator, admin bool) {
	atExp := time.Now().Add(time.Minute * 60 * 1).Unix()
	// atExp := time.Now().Add(time.Minute * 1 * 1).Unix()

	token.AtExpires = atExp
	token.AccessUuid = uuid.NewString()

	var err error

	//*identfy what is allowed the user to do
	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = id
	atClaims["moderator"] = moderator
	atClaims["access_uuid"] = token.AccessUuid
	atClaims["admin"] = admin
	atClaims["exp"] = token.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims) //sign the token
	token.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		panic(err)
	}
}
