package service

import (
	"time"

	"test/conf"
	"test/models"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type JWTService interface {
	CreateLoginToken(user models.User) (string, error)
	VerifyLoginToken(tokenStr string) (*models.DecodeJWTClaims, error)
}

type jwtService struct {
	cacheService CacheService
}

func NewJWTService(cacheService CacheService) JWTService {
	return &jwtService{
		cacheService: cacheService,
	}
}

func (j *jwtService) CreateLoginToken(user models.User) (string, error) {
	claim := jwt.MapClaims{
		"exp":       time.Now().Add(time.Minute * 5).Unix(),
		"uid":       user.Id,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	secret := viper.Get(conf.JwtSecret).(string)

	j.cacheService.SetValue(user.Email, token, time.Minute*5)
	return token.SignedString([]byte(secret))
}

/*
* This method does following
* get jwt secret from conf, parse incoming token, get data stored in token, check if token has expired, check if uid exist,
* check active account id value, check account type
* @params token string
* @returns decoded jwt token, erro
 */
func (j jwtService) VerifyLoginToken(tokenStr string) (*models.DecodeJWTClaims, error) {
	secret := viper.Get(conf.JwtSecret).(string)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, getStandardError(err, models.TOKEN_EXPIRED, "VerifyLoginToken():56", err.Error())
	}
	if !token.Valid {
		return nil, getStandardError(err, models.INVALID_TOKEN, "VerifyLoginToken():59", models.INVALID_TOKEN_MESSAGE)
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, getStandardError(err, models.INVALID_TOKEN, "VerifyLoginToken():63", "cannot convert claim to MapClaims")
	}

	expiredAtVal, found := claim["exp"]
	if !found {
		return nil, getStandardError(err, models.INVALID_TOKEN, "VerifyLoginToken():69", models.INVALID_TOKEN_MESSAGE)
	}
	err = claim.Valid()
	if err != nil {
		return nil, getStandardError(err, models.TOKEN_EXPIRED, "VerifyLoginToken():73", err.Error())
	}

	uidVal, found := claim["uid"]
	if !found {
		return nil, getStandardError(err, models.INVALID_TOKEN, "VerifyLoginToken():77", models.INVALID_TOKEN_MESSAGE)
	}
	uid := uidVal.(float64)
	if uid <= 0 {
		return nil, getStandardError(err, models.INVALID_TOKEN, "VerifyLoginToken():82", models.INVALID_TOKEN_MESSAGE)
	}

	email, found := claim["email"]
	if !found {
		return nil, getStandardError(err, models.INVALID_TOKEN, "VerifyLoginToken():87", models.INVALID_TOKEN_MESSAGE)
	}

	return &models.DecodeJWTClaims{
		ExpiredAt: expiredAtVal.(float64),
		UID:       uid,
		Email:     email.(string),
	}, nil
}
