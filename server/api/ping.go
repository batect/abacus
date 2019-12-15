package api

import (
	"fmt"
	"net/http"
)

func Ping(w http.ResponseWriter, req *http.Request) {
	if _, err := fmt.Fprint(w, "pong"); err != nil {
		panic(err)
	}
}
