package myutils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL   string
	Alive bool
	mu    sync.RWMutex
}

var backendServers = []*Backend{
	{URL: "http://localhost:8081", Alive: true},
	{URL: "http://localhost:8082", Alive: true},
	{URL: "http://localhost:8083", Alive: true},
}
var nextServerIndex = 0

func healthChecker() {
	for {
		for _, server := range backendServers {
			resp, err := http.Get(server.URL)
			server.mu.Lock()
			server.Alive = (err == nil && resp.StatusCode == http.StatusOK)
			server.mu.Unlock()

			if resp != nil {
				resp.Body.Close()
			}

		}
		time.Sleep(10 * time.Second)
	}
}

func getNextHealthyBackend() *Backend {
	for i, server := range backendServers {
		server.mu.RLock()
		if server.Alive && i >= nextServerIndex {
			server.mu.RUnlock()
			return server
		}
		server.mu.RUnlock()
	}

	for i, server := range backendServers {
		server.mu.RLock()
		if server.Alive {
			server.mu.RUnlock()
			nextServerIndex = i
			return server
		}
		server.mu.RUnlock()
	}

	return nil
}

func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	PrintRequest(r)

	server := getNextHealthyBackend()
	if server == nil {
		http.Error(w, "All servers are down", http.StatusServiceUnavailable)
		return
	}
	target, _ := url.Parse(server.URL)
	nextServerIndex += 1
	nextServerIndex %= len(backendServers)

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

	go healthChecker()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Println("Starting load balancer on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
