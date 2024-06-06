package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

const size = 1024
const interval = 5 * time.Second

func CPUHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		end := time.Now().Add(interval)
		for time.Now().Before(end) {
			for i := 0; i < 1000000; i++ {
				_ = i * i
			}
		}
		fmt.Print("finish cpu compute")
	}()
	w.WriteHeader(http.StatusOK)
}

func MemHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		memory := make([]byte, size*1024*1024)
		time.Sleep(interval)
		_ = memory
		fmt.Print("finish mem alloc\n")
	}()
	w.WriteHeader(http.StatusOK)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/cpu", CPUHandler)
	http.HandleFunc("/memory", MemHandler)
	err := http.ListenAndServe(":7070", nil)
	if err != nil {
		fmt.Print(err.Error())
	}
}
