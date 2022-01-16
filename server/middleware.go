package server

import (
	"encoding/json"
	"net/http"
	"self-scientists/config"
	"self-scientists/data"
	"time"
)

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
		extractedClaims := verifyToken(extractedTokenString)
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
