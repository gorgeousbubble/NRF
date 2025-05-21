package main

import (
	"fmt"
	"net/http"
	. "nrf/app"
	. "nrf/logs"
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
	nrf := New()
	err := nrf.Init()
	if err != nil {
		fmt.Println("The NRF initialization failed:", err.Error())
		L.Error("The NRF initialization failed:", err.Error())
		os.Exit(1)
	}
	nrf.Start()
}
