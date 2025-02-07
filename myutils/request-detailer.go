package myutils

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func PrintRequest(r *http.Request) {
	sender, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Fatal(err)
	}

	method := r.Method

	protocolVersion := r.Proto

	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		log.Fatal(err)
	}

	userAgent := r.Header.Get("User-Agent")

	accept := r.Header.Get("Accept")

	fmt.Println("Received request from", sender)
	fmt.Println(method, "/", protocolVersion)
	fmt.Println("Host:", host)
	fmt.Println("User-Agent:", userAgent)
	fmt.Println("Accept:", accept)
}
