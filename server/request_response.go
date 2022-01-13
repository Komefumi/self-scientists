package server

type standardResponse struct {
	Status  uint8       `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}

var emptyData = struct{}{}
var emptyErrors = []string{}
