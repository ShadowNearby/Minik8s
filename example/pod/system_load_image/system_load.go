package main

import (
	"encoding/json"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type LoadRequest struct {
	CPUDuration    int `json:"cpu_duration"`
	CPULoad        int `json:"cpu_load"`
	MemorySize     int `json:"memory_size"`
	MemoryDuration int `json:"memory_duration"`
}

func cpuLoad(duration, load int, wg *sync.WaitGroup) {
	defer wg.Done()
	end := time.Now().Add(time.Duration(duration) * time.Second)
	for time.Now().Before(end) {
		for i := 0; i < load*1000000; i++ {
			_ = i * i
		}
	}
}

func memoryLoad(size, duration int, wg *sync.WaitGroup) {
	defer wg.Done()
	memory := make([]byte, size*1024*1024)
	time.Sleep(time.Duration(duration) * time.Second)
	_ = memory // 使用 memory 以防止编译器优化
}

func loadHandler(w http.ResponseWriter, r *http.Request) {
	var loadReq LoadRequest
	if err := json.NewDecoder(r.Body).Decode(&loadReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup

	if loadReq.CPUDuration > 0 && loadReq.CPULoad > 0 {
		wg.Add(1)
		go cpuLoad(loadReq.CPUDuration, loadReq.CPULoad, &wg)
	}

	if loadReq.MemorySize > 0 && loadReq.MemoryDuration > 0 {
		wg.Add(1)
		go memoryLoad(loadReq.MemorySize, loadReq.MemoryDuration, &wg)
	}

	wg.Wait()
	w.Write([]byte(`{"status": "load completed"}`))
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // 使用所有可用的CPU核心
	http.HandleFunc("/load", loadHandler)
	http.ListenAndServe(":8080", nil)
}
