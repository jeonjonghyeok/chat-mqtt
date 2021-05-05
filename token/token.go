package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const mySigningKey = "32iazLZ3hD4aH4EKjRkEo3is"

type customClaims struct {
	jwt.StandardClaims
	UserID int    `json:"uid"`
	Ttoken string `json:"ttoken"`
}

// New 새 token 만듦
func New(userID int, ttoken string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims{
		UserID: userID,
		Ttoken: ttoken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
			Issuer:    "jjh",
		},
	})
	return token.SignedString([]byte(mySigningKey))
}

// Parse token을 parse하고 사용자 id를 리턴
func Parse(token string) (userID int, err error) {
	parsed, err := jwt.ParseWithClaims(token, &customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		})

	if err != nil {
		return 0, err
	}

	if !parsed.Valid {
		return 0, errors.New("token is invalid")
	}

	if c, ok := parsed.Claims.(*customClaims); ok {
		return c.UserID, nil
	}

	return 0, errors.New("token is invalid")
}

func Tparse(token string) (ttoken string, err error) {
	parsed, err := jwt.ParseWithClaims(token, &customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		})

	if err != nil {
		return "", err
	}

	if !parsed.Valid {
		return "", errors.New("token is invalid")
	}

	if c, ok := parsed.Claims.(*customClaims); ok {
		return c.Ttoken, nil
	}
	return "", errors.New("token is invalid")
}