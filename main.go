package main

import (
	"load-balancer/network/myutils"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		myutils.StartBackendServer(8082)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		myutils.StartBackendServer(8081)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		myutils.StartBackendServer(8083)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		myutils.StartLoadBalancer()
	}()

	wg.Wait()
}
