package server

import (
	"encoding/json"
	"net/http"
	"self-scientists/config"
	"self-scientists/data"
	"time"
)

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resp standardResponse
	switch r.Method {
	case http.MethodPost:
		{
			var newUser data.User
			err := json.NewDecoder(r.Body).Decode(&newUser)
			if err != nil {
				w.WriteHeader(400)
				resp = responseForInvalidRequestBody
				break
			}
			errors, internalServerError := newUser.CreateUser()
			if internalServerError {
				w.WriteHeader(500)
				resp = responseForInternalServerError
				break
			}
			if len(errors) > 0 {
				w.WriteHeader(400)
				resp = standardResponse{Status: 1, Message: "Error In User Registration, Check errors field", Data: emptyData, Errors: errors}
				break
			}
			resp = standardResponse{Status: 0, Message: "User Registration Success!", Data: emptyData, Errors: emptyErrors}
		}
	default:
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
			resp = standardResponse{Status: 1, Message: "Invalid Method", Data: emptyData, Errors: emptyErrors}
		}
	}
	json.NewEncoder(w).Encode(resp)
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resp standardResponse
	switch r.Method {
	case http.MethodPost:
		{
			var ag authGate
			err := json.NewDecoder(r.Body).Decode(&ag)
			if err != nil {
				w.WriteHeader(400)
				resp = responseForInvalidRequestBody
				break
			}
			token, errors, internallyErrored := ag.AuthenticateAndCreateToken()
			if internallyErrored {
				w.WriteHeader(500)
				resp = responseForInternalServerError
				break
			}
			if len(errors) > 0 {
				w.WriteHeader(400)
				resp = standardResponse{Status: 1, Message: "Error In User Authentication, Check errors field", Data: emptyData, Errors: errors}
				break
			}
			returnableData := make(map[string]interface{})
			returnableData["token"] = token
			resp = standardResponse{Status: 0, Message: "Successfully logged in and created AUTH token", Data: returnableData, Errors: emptyErrors}
		}
	default:
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
	json.NewEncoder(w).Encode(resp)
	return
}

func tokenCheckMiddleware(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp standardResponse

		w.Header().Set("Content-Type", "application/json")
		headerTokenString := r.Header.Get(config.AUTH_HEADER_NAME)
		tokenFormatSatisfied := headerAccessTokenRegexp.MatchString(headerTokenString)
		if !tokenFormatSatisfied {
			w.WriteHeader(400)
			resp = standardResponse{Status: 1, Message: noticeToUserForMalformedToken, Data: emptyData, Errors: emptyErrors}
			json.NewEncoder(w).Encode(resp)
			return
		}
		extractedTokenString := headerAccessTokenRegexp.FindAllStringSubmatch(headerTokenString, -1)[0][1]
		extractedClaims := VerifyToken(extractedTokenString)
		if extractedClaims == nil {
			w.WriteHeader(400)
			resp = standardResponse{Status: 1, Message: noticeToUserForTokenValidationFailure, Data: emptyData, Errors: emptyErrors}
			json.NewEncoder(w).Encode(resp)
			return
		}
		tExpiry := time.Unix(extractedClaims.ExpiresAt, 0)
		timeOfExpiryAlreadyReached := time.Now().After(tExpiry)
		if timeOfExpiryAlreadyReached {
			w.WriteHeader(400)
			resp = standardResponse{Status: 1, Message: noticeToUserForTokenExpiration, Data: emptyData, Errors: emptyErrors}
			json.NewEncoder(w).Encode(resp)
			return
		}

		userExists := data.DoesUserOfIDExist(extractedClaims.ID)
		if !userExists {
			w.WriteHeader(400)
			resp = standardResponse{Status: 1, Message: noticeToUserForAccountNotExistingDespiteValidToken, Data: emptyData, Errors: emptyErrors}
			json.NewEncoder(w).Encode(resp)
			return
		}

		r.Header.Set(config.READY_TOKEN_STRING_HEADER_NAME, extractedTokenString)
		next(w, r)
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	var resp standardResponse
	respBody := make(map[string]interface{})
	tokenString := r.Header.Get(config.READY_TOKEN_STRING_HEADER_NAME)
	parsedClaims := VerifyToken(tokenString)
	respBody["extractedClaims"] = parsedClaims
	resp = standardResponse{Status: 0, Message: "Test Successful", Data: respBody, Errors: emptyErrors}
	json.NewEncoder(w).Encode(resp)
	return
}
