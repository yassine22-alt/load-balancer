package myutils

import (
	"fmt"
	"log"
	"net/http"
)

func backend(w http.ResponseWriter, r *http.Request) {
	PrintRequest(r)
	fmt.Fprintln(w, "Hello from backend server!")
}

func StartBackendServer() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", backend)

	log.Println("Starting backend server on :8081")
	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatal(err)
	}
}
