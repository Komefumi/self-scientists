package server

import (
	"encoding/json"
	"net/http"
	"self-scientists/config"
	"self-scientists/data"
	"strconv"
)

const (
	invalid_body_problem = iota
	internal_server_error_problem
	request_fields_validation_problem
)

func handleProblem(w http.ResponseWriter, r *http.Request, problemInteger int, passedErrors []string) {
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
			resp = standardResponse{Status: 1, Message: "Error processing request, check errors field", Data: emptyData, Errors: passedErrors}
		}
	default:
		{
			panic("Must provide problemInteger")
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func handleInvalidBodyProblem(w http.ResponseWriter, r *http.Request) {
	handleProblem(w, r, invalid_body_problem, emptyErrors)
}

func handleInternalServerErrorProblem(w http.ResponseWriter, r *http.Request) {
	handleProblem(w, r, internal_server_error_problem, emptyErrors)
}

func handleRequestFieldsValidationProblem(w http.ResponseWriter, r *http.Request, errors []string) {
	handleProblem(w, r, request_fields_validation_problem, errors)
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
				handleInvalidBodyProblem(w, r)
				return
			}
			errors, internalServerError := newUser.CreateUser()
			if internalServerError {
				handleInternalServerErrorProblem(w, r)
				return
			}
			if len(errors) > 0 {
				handleRequestFieldsValidationProblem(w, r, errors)
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
				handleInvalidBodyProblem(w, r)
				return
			}
			token, errors, internallyErrored := ag.AuthenticateAndCreateToken()
			if internallyErrored {
				handleInternalServerErrorProblem(w, r)
				return
			}
			if len(errors) > 0 {
				handleRequestFieldsValidationProblem(w, r, errors)
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
	parsedClaims := verifyToken(tokenString)
	respBody["extractedClaims"] = parsedClaims
	resp = standardResponse{Status: 0, Message: "Test Successful", Data: respBody, Errors: emptyErrors}
	json.NewEncoder(w).Encode(resp)
	return
}

func getThreadByIDHandler(w http.ResponseWriter, r *http.Request, id string) {
	threadId, idParseErr := strconv.ParseUint(id, 10, 64)
	if idParseErr != nil {
		w.WriteHeader(400)
		resp := standardResponse{Status: 0, Message: "Error: Must provide valid id to retrieve thread by id", Data: emptyData, Errors: emptyErrors}
		json.NewEncoder(w).Encode(resp)
		return
	}
	threadData, internallyErrored := data.GetThreadById(uint(threadId))
	if internallyErrored {
		handleInternalServerErrorProblem(w, r)
		return
	}
	if threadData == nil {
		w.WriteHeader(404)
		resp := standardResponse{Status: 0, Message: "Error: Seems like no thread with requested id was found", Data: emptyData, Errors: emptyErrors}
		json.NewEncoder(w).Encode(resp)
	} else {
		returnData := make(map[string]interface{})
		returnData["thread"] = threadData
		resp := standardResponse{Status: 0, Message: "Thread successfully retrieved", Data: returnData, Errors: emptyErrors}
		json.NewEncoder(w).Encode(resp)
	}

}

func getThreadListByPageHandler(w http.ResponseWriter, r *http.Request, pageNumberString string) {
	pageNumber, idParseErr := strconv.ParseUint(pageNumberString, 10, 64)
	if idParseErr != nil {
		w.WriteHeader(400)
		resp := standardResponse{Status: 0, Message: "Error: Must provide valid page number to retrieve threads by page", Data: emptyData, Errors: emptyErrors}
		json.NewEncoder(w).Encode(resp)
		return
	}

	threadDataList, internallyErrored := data.GetThreadListByPage(uint(pageNumber))

	if internallyErrored {
		handleInternalServerErrorProblem(w, r)
		return
	}

	returnData := make(map[string]interface{})
	returnData["pageSize"] = config.ThreadPaginationSize
	if len(threadDataList) == 0 {
		returnData["empty"] = true
	} else {
		returnData["empty"] = false
	}
	returnData["threads"] = threadDataList
	resp := standardResponse{Status: 0, Message: "Threads successfully retrieved", Data: returnData, Errors: emptyErrors}
	json.NewEncoder(w).Encode(resp)
}

func threadsRetrievalHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	threadId := query.Get("id")
	if len(threadId) != 0 {
		getThreadByIDHandler(w, r, threadId)
		return
	}

	pageNumber := query.Get("page")
	if len(pageNumber) != 0 {
		getThreadListByPageHandler(w, r, pageNumber)
		return
	}

	w.WriteHeader(400)
	errMessage := "Error: Must provide at least one of the following in query params: id (to get specific thread), page (to get a paginated set of threads)"
	resp := standardResponse{Status: 1, Message: errMessage, Data: emptyData, Errors: emptyErrors}

	json.NewEncoder(w).Encode(resp)
	return
}

func createThreadHandler(w http.ResponseWriter, r *http.Request) {
	authClaims := getDecodedAuthClaims(r)
	var newThread data.Thread
	decodeErr := json.NewDecoder(r.Body).Decode(&newThread)
	if decodeErr != nil {
		handleInvalidBodyProblem(w, r)
		return
	}
	errors, internallyErrored := newThread.CreateThread(authClaims.ID)
	if internallyErrored {
		handleInternalServerErrorProblem(w, r)
		return
	}
	if len(errors) > 0 {
		handleRequestFieldsValidationProblem(w, r, errors)
		return
	}
	resp := standardResponse{Status: 0, Message: "Thread Successfully Created", Data: emptyData, Errors: emptyErrors}
	json.NewEncoder(w).Encode(resp)
}

func threadsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			threadsRetrievalHandler(w, r)
			return
		}
	case http.MethodPost:
		{
			createThreadHandler(w, r)
			return
		}
	default:
		{
			testHandler(w, r)
			return
		}
	}
}
