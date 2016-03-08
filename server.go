package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	//reading environment specific settings .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	//preparing mux and server
	conn := fmt.Sprint(host, ":", port)
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(http.Dir("./client")))
	router.Handle("/ws", wsHandler{})

	//serving
	log.Printf("serving game theory on %v", conn)
	log.Fatal(http.ListenAndServe(conn, router))
}
