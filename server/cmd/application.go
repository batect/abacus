package main

import (
	"log"
	"net/http"

	"github.com/batect/abacus/server/api"
)

func main() {
	http.HandleFunc("/ping", api.Ping)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
