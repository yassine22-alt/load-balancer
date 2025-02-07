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
		myutils.StartBackendServer()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		myutils.StartLoadBalancer()
	}()

	wg.Wait()
}
