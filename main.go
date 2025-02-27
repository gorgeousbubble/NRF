package main

import (
	"fmt"
	"net/http"
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
	fmt.Println("Hello World")
}
