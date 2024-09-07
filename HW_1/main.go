package main

import (
	"HW/server"
	"HW/storage"
	"flag"
	"log"
)

// @title			My API
// @version		1.0
// @description	This is a sample server.
// @host			localhost:8080
// @BasePath		/
func main() {
	addr := flag.String("addr", ":8080", "address for http server")

	s := storage.NewRaiStorage()

	log.Printf("Starting server on %s", *addr)
	if err := server.CreateAndRunServer(s, *addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
