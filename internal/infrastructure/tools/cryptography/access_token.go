package cryptography

import (
	"errors"
	"fmt"
	"github.com/KuYaki/waffler_server/config"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	AccessToken = iota
	RefreshToken
)

type TokenManager interface {
	CreateToken(userID string, ttl time.Duration, kind int) (string, error)
	ParseToken(inputToken string, kind int) (UserClaims, error)
	ParseTokenForHTTP(w http.ResponseWriter, r *http.Request) (*UserFromClaims, error)
}

type TokenJWT struct {
	AccessSecret  []byte
	RefreshSecret []byte
}

func NewTokenJWT(token *config.Token) TokenManager {
	return &TokenJWT{AccessSecret: []byte(token.AccessSecret), RefreshSecret: []byte(token.RefreshSecret)}
}

// UserClaims include custom claims on jwt.
type UserClaims struct {
	ID string `json:"uid"`
	jwt.RegisteredClaims
}

type UserFromClaims struct {
	ID int
}

// CreateToken create new token with parameters.
func (o *TokenJWT) CreateToken(userID string, ttl time.Duration, kind int) (string, error) {
	claims := UserClaims{
		ID:               userID,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl))},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var secret []byte
	switch kind {
	case AccessToken:
		secret = o.AccessSecret
	case RefreshToken:
		secret = o.RefreshSecret
	default:
		return "", errors.New("token type error")
	}

	return token.SignedString(secret)
}

// ParseToken parsing input token, and return email and role from token.
func (o *TokenJWT) ParseToken(inputToken string, kind int) (UserClaims, error) {
	token, err := jwt.Parse(inputToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		var secret []byte
		switch kind {
		case AccessToken:
			secret = o.AccessSecret
		case RefreshToken:
			secret = o.RefreshSecret
		default:
			return "", errors.New("token type error")
		}
		_ = secret

		return secret, nil
	})

	if err != nil {
		return UserClaims{}, err
	}

	if !token.Valid {
		return UserClaims{}, fmt.Errorf("not valid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return UserClaims{}, fmt.Errorf("error get user claims from token")
	}

	return UserClaims{
		ID:               claims["uid"].(string),
		RegisteredClaims: jwt.RegisteredClaims{},
	}, nil
}

func (o *TokenJWT) ParseTokenForHTTP(w http.ResponseWriter, r *http.Request) (*UserFromClaims, error) {
	res := &UserFromClaims{}

	tokenRaw := c.GetHeader("Authorization")
	log.Println(tokenRaw)
	tokenParts := strings.Split(tokenRaw, " ")
	if len(tokenParts) < 2 && tokenParts[0] != "Bearer" {
		return nil, errors.New("wrong input data")
	}
	u, err := o.ParseToken(tokenParts[1], AccessToken)
	if err != nil && err.Error() == "Token is expired" {
		return nil, err
	}

	idSting, err := strconv.Atoi(u.ID)
	if err != nil {
		return nil, err
	}

	res.ID = idSting
	return res, nil
}
