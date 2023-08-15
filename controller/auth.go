package controller

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/willywartono14/test-go-gin-gonic/config"
	"github.com/willywartono14/test-go-gin-gonic/model"
)

var ErrTokenExpired error = errors.New("token has expired")
var ErrTokenInvalid error = errors.New("token is invalid")
var ErrCredentialInvalid error = errors.New("credential is invalid")
var ErrUsernameExist error = errors.New("username already exist")

type (
	AuthController interface {
		Login(ctx *gin.Context, username, password string) (model.Token, error)
		Register(ctx *gin.Context, userRequest model.User) (model.Token, error)
	}
	key     string
	Payload struct {
		UserId int64 `json:"user_id"`
		jwt.StandardClaims
	}
)

func (c *controller) Login(ctx *gin.Context, username, password string) (model.Token, error) {

	user, err := model.User{}.FindUserByUsername(ctx.Request.Context(), c.db, username)
	if err != nil {
		return model.Token{}, ErrCredentialInvalid
	}

	if user.Password != password {
		return model.Token{}, ErrCredentialInvalid
	}

	expToken := 30 * time.Minute

	token, err := c.createToken(expToken, int64(user.ID))
	if err != nil {
		return model.Token{}, ErrCredentialInvalid
	}

	return model.Token{
		AccessToken: token,
		ExpiredTime: time.Now().Add(expToken),
	}, nil
}

func (c *controller) Register(ctx *gin.Context, userRequest model.User) (model.Token, error) {

	user, err := model.User{}.FindUserByUsername(ctx.Request.Context(), c.db, userRequest.Username)
	if err != nil {
		return model.Token{}, ErrCredentialInvalid
	}

	if user.Username == userRequest.Username {
		return model.Token{}, ErrUsernameExist
	}

	id, err := model.User{}.InsertUser(ctx.Request.Context(), c.db, userRequest)
	if err != nil {
		return model.Token{}, ErrCredentialInvalid
	}

	expToken := 30 * time.Minute

	token, err := c.createToken(expToken, int64(id))
	if err != nil {
		return model.Token{}, ErrCredentialInvalid
	}

	return model.Token{
		AccessToken: token,
		ExpiredTime: time.Now().Add(expToken),
	}, nil
}

func (*controller) createToken(duration time.Duration, userId int64) (string, error) {
	tokenPayload := Payload{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenPayload)
	token, err := jwtToken.SignedString([]byte(config.Get().Jwt.SecretKey))

	return token, err
}

func (*controller) verifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrTokenInvalid
		}
		return []byte(config.Get().Jwt.SecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && verr.Errors == jwt.ValidationErrorExpired {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrTokenInvalid
	}

	return payload, nil
}
