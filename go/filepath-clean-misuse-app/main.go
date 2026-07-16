package main

import "log"

// Single entry point for this application.
func main() {
	log.Println("listening on :8080")
	log.Fatal(StartServer())
}
