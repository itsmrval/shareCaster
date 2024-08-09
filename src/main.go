package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	initConfig()
	initDB()
	initRoutes()
	go cleanupBuffers()

	log.Println("[INFO] Serving API on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	defer closeDb()
}
