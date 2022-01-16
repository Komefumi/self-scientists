package server

import (
	"encoding/json"
	"net/http"
	"self-scientists/config"
	"self-scientists/data"
)

const (
	invalid_body_problem = iota
	internal_server_error_problem
	request_fields_validation_problem
)

func handleProblem(w http.ResponseWriter, r *http.Request, problemInteger int, passedErrors []string) {
	var errors []string
	if len(passedErrors) > 0 {
		errors = passedErrors
	} else {
		errors = emptyErrors
	}
	var resp standardResponse
	switch problemInteger {
	case invalid_body_problem:
		{
			w.WriteHeader(400)
			resp = responseForInvalidRequestBody
		}
	case internal_server_error_problem:
		{
			w.WriteHeader(500)
			resp = responseForInternalServerError
		}
	case request_fields_validation_problem:
		{
			w.WriteHeader(400)
			resp = standardResponse{Status: 1, Message: "Error processing request, check errors field", Data: emptyData, Errors: errors}
		}
	default:
		{
			panic("Must provide problemInteger")
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resp standardResponse
	switch r.Method {
	case http.MethodPost:
		{
			var newUser data.User
			err := json.NewDecoder(r.Body).Decode(&newUser)
			if err != nil {
				handleProblem(w, r, invalid_body_problem, emptyErrors)
				return
			}
			errors, internalServerError := newUser.CreateUser()
			if internalServerError {
				handleProblem(w, r, internal_server_error_problem, emptyErrors)
				return
			}
			if len(errors) > 0 {
				handleProblem(w, r, invalid_body_problem, errors)
				return
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
				handleProblem(w, r, invalid_body_problem, emptyErrors)
				return
			}
			token, errors, internallyErrored := ag.AuthenticateAndCreateToken()
			if internallyErrored {
				handleProblem(w, r, internal_server_error_problem, emptyErrors)
				return
			}
			if len(errors) > 0 {
				handleProblem(w, r, request_fields_validation_problem, errors)
				return
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

func getThreadsHandler(w http.ResponseWriter, r *http.Request) {
	resp := standardResponse{Status: 0, Message: "Test Successful", Data: emptyData, Errors: emptyErrors}
	json.NewEncoder(w).Encode(resp)
	return
}

/*
func createThreadsHandler(w http.ResponseWriter, r *http.Request) {

}
*/

func threadsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	/*
		case http.MethodPost:
			{
				var newThread data.Thread
				err := json.NewDecoder(r.Body).Decode(&newThread)
			}
	*/
	default:
		{
			testHandler(w, r)
			return
		}
	}
}
