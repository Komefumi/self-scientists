package main

import (
	"encoding/json"
	"net/http"
	"self-scientists/data"
)

type standardResponse struct {
	Status  uint8       `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}

var emptyData = struct{}{}
var emptyErrors = []string{}

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
				resp = standardResponse{Status: 1, Message: "Invalid request body", Data: emptyData, Errors: emptyErrors}
				break
			}
			errors, internalServerError := newUser.CreateUser()
			if internalServerError {
				w.WriteHeader(500)
				resp = standardResponse{Status: 2, Message: "Internal Server Error", Data: emptyData, Errors: emptyErrors}
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
			w.WriteHeader(400)
			resp = standardResponse{Status: 1, Message: "Invalid Method", Data: emptyData, Errors: emptyErrors}
		}
	}
	json.NewEncoder(w).Encode(resp)
	return
}
