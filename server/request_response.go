package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"self-scientists/config"
	"self-scientists/utils"
	"self-scientists/validation"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const noticeToUserForMalformedToken = "Malformed Token or Token not provided. Provide token in `Token: (token)` format as Authorization Header"
const noticeToUserForTokenValidationFailure = "Token validation failure. Retry with a new access token (login again)"
const noticeToUserForTokenExpiration = "Token has expired, login again to get a new token to try again"
const noticeToUserForAccountNotExistingDespiteValidToken = "Though token is valid, account does not seem to exist. If you recently deleted your account, try logging out or clearing your cache"

var headerAccessTokenRegexp *regexp.Regexp = regexp.MustCompile("^Token: ([a-zA-Z0-9-_.]+)$")

var signingKey = []byte(config.JWTSecret)

type standardResponse struct {
	Status  uint8       `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}

type authClaims struct {
	Email string `json:"email"`
	ID    uint   `json:"id"`
	jwt.StandardClaims
}

type authGate struct {
	Email    string `json:"email"`
	ID       uint   `json:"id"`
	Password string `json:"password"`
}

var emptyData = struct{}{}
var emptyErrors = []string{}

var responseForInvalidRequestBody = standardResponse{Status: 1, Message: "Invalid request body", Data: emptyData, Errors: emptyErrors}
var responseForInternalServerError = standardResponse{Status: 2, Message: "Internal Server Error", Data: emptyData, Errors: emptyErrors}

var defaultAuthFailureError = "User with does not exist or password invalid"

func (ag *authGate) validateForAuth() (errors []string, internallyErrored bool) {
	internallyErrored = false
	if !validation.IsEmail(ag.Email) {
		errors = append(errors, "Must provide a valid email")
	}
	if len(ag.Password) == 0 {
		errors = append(errors, "Must provide a password")
	}
	var userPasswordHash string
	var userId uint
	row := config.DB.QueryRow("SELECT id, password_hash FROM users WHERE email=$1", ag.Email)
	err := row.Scan(&userId, &userPasswordHash)
	fmt.Println("UserId below")
	fmt.Println(userId)
	switch err {
	case sql.ErrNoRows:
		{
			errors = append(errors, defaultAuthFailureError)
			return errors, internallyErrored
		}
	case nil:
		{
			ag.ID = userId
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
		ID:    ag.ID,
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

func verifyToken(tokenString string) *authClaims {
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

func getDecodedAuthClaims(r *http.Request) *authClaims {
	tokenString := r.Header.Get(config.READY_TOKEN_STRING_HEADER_NAME)
	parsedClaims := verifyToken(tokenString)
	return parsedClaims
}
