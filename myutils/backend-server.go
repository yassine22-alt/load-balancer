package myutils

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func backend(w http.ResponseWriter, r *http.Request) {
	PrintRequest(r)

	serverAddr := r.Context().Value(http.LocalAddrContextKey).(net.Addr).String()

	fmt.Fprintln(w, "Hello from backend server on port", serverAddr)
}

func StartBackendServer(port int) {

	mux := http.NewServeMux()
	mux.HandleFunc("/", backend)

	log.Println("Starting backend server on :", port)
	err := http.ListenAndServe(":"+fmt.Sprint(port), mux)
	if err != nil {
		log.Fatal(err)
	}
}
