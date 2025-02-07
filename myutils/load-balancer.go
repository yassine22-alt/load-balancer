package myutils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	PrintRequest(r)

	target, _ := url.Parse("http://localhost:8081")

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ModifyResponse = func(resp *http.Response) error {
		fmt.Println("\nResponse from server:", resp.Proto, resp.Status)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		fmt.Println("\n", string(body))

		resp.Body = io.NopCloser(bytes.NewReader(body))
		return nil
	}

	proxy.ServeHTTP(w, r)
}

func StartLoadBalancer() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Println("Starting load balancer on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
