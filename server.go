package main

import (
	"self-scientists/config"
	"self-scientists/server"
)

func main() {
	defer config.CloseDB()
	server.SetupServer()
}
