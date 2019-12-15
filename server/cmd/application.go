package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/batect/abacus/server/api"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT environment variable is not set.")
	}

	address := fmt.Sprintf(":%s", port)

	http.HandleFunc("/ping", api.Ping)
	log.Fatal(http.ListenAndServe(address, nil))
}
