package main

import (
	"log"
	"qwixx/internal/server"
)

func main() {
	settings := server.Settings{
		Endpoint: "localhost:8080",
	}
	s := server.New()
	err := s.Start(settings)
	if err != nil {
		log.Fatal(err)
	}
}
