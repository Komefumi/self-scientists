package server

import (
	"database/sql"
	"fmt"
	"self-scientists/config"
	"self-scientists/utils"
	"self-scientists/validation"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var signingKey = []byte(config.JWTSecret)

type standardResponse struct {
	Status  uint8       `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}

type authClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type authGate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var emptyData = struct{}{}
var emptyErrors = []string{}

var responseForInvalidRequestBody = standardResponse{Status: 1, Message: "Invalid request body", Data: emptyData, Errors: emptyErrors}
var responseForInternalServerError = standardResponse{Status: 2, Message: "Internal Server Error", Data: emptyData, Errors: emptyErrors}

var defaultAuthFailureError = "User with does not exist or password invalid"

func (ag authGate) validateForAuth() (errors []string, internallyErrored bool) {
	internallyErrored = false
	if !validation.IsEmail(ag.Email) {
		errors = append(errors, "Must provide a valid email")
	}
	if len(ag.Password) == 0 {
		errors = append(errors, "Must provide a password")
	}
	var userPasswordHash string
	row := config.DB.QueryRow("SELECT password_hash FROM users WHERE email=$1", ag.Email)
	err := row.Scan(&userPasswordHash)
	switch err {
	case sql.ErrNoRows:
		{
			errors = append(errors, defaultAuthFailureError)
			return errors, internallyErrored
		}
	case nil:
		{
			break
		}
	default:
		{
			fmt.Println(err)
			internallyErrored = true
			return errors, internallyErrored
		}
	}
	if matched := utils.VerifyPassword(userPasswordHash, ag.Password); !matched {
		errors = append(errors, defaultAuthFailureError)
	}

	return errors, internallyErrored
}

func (ag authGate) AuthenticateAndCreateToken() (tokenString string, errors []string, internallyErrored bool) {
	errors, internallyErrored = ag.validateForAuth()
	if internallyErrored || len(errors) > 0 {
		return tokenString, errors, internallyErrored
	}

	ttl := time.Second * 60 * 60 * 24 * 10
	timeOfExpiry := time.Now().UTC().Add(ttl).Unix()
	claims := authClaims{
		Email: ag.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: timeOfExpiry,
			Issuer:    "self-scientists.github.io",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, errInternal := token.SignedString(signingKey)
	if errInternal != nil {
		return tokenString, errors, true
	}

	return tokenString, errors, false
}

func VerifyToken(tokenString string) *authClaims {
	token, _ := jwt.ParseWithClaims(tokenString, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if claims, ok := token.Claims.(*authClaims); ok && token.Valid {
		fmt.Println(claims)
		return claims
	} else {
		return nil
	}
}
