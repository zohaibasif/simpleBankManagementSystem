package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	MIN_SECRET_KEY_SIZE = 32
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {

	if len(secretKey) < MIN_SECRET_KEY_SIZE {
		return nil, fmt.Errorf("invalid key size : should be at least %d characterrs", MIN_SECRET_KEY_SIZE)
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

func (jwtm *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	JWTtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return JWTtoken.SignedString([]byte(jwtm.secretKey))
}

func (jwtm *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, Err_Invalid_Token
		}

		return []byte(jwtm.secretKey), nil
	}

	JWTtoken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, Err_Expired_Token) {
			return nil, Err_Expired_Token
		}
		return nil, Err_Invalid_Token
	}

	payload, ok := JWTtoken.Claims.(*Payload)
	if !ok {
		return nil, Err_Invalid_Token
	}

	return payload, nil
}
