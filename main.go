package main

import (
	"fmt"
	"net/http"
	. "nrf/app"
	"runtime"
)

func init() {
	// start multi-cpu
	core := runtime.NumCPU()
	runtime.GOMAXPROCS(core)
	// start debug pprof
	go func() {
		_ = http.ListenAndServe(":10514", nil)
	}()
}

func main() {
	fmt.Println("The 5G System Network Function Repository Services.")
	NRFService = New()
	NRFService.Init()
	NRFService.Start()
	fmt.Println("The NRF is running...")
}
