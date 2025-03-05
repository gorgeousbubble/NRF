package main

import (
	"fmt"
	"net/http"
	. "nrf/app"
	"os"
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
	err := NRFService.Init()
	if err != nil {
		fmt.Println("The NRF initialization failed:", err.Error())
		os.Exit(1)
	}
	NRFService.Start()
	fmt.Println("The NRF is running...")
}
