package main

import (
	"log"
	"os"
)

var port string

func initConfig() {
	log.Println("[INFO] Reading configuration")
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}
