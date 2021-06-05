package main

import (
	"log"

	"github.com/rezaAmiri123/proglog/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
