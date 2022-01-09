package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	switch r.Method {
	case http.MethodPost:
		{
			resp["message"] = "Post Success"
		}
	default:
		{
			resp["message"] = "Invalid Method"
		}
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
