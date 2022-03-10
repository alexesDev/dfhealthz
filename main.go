package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
)

func main() {
	log.SetFlags(0)

	addr := getenv("ADDR", ":5655")
	endpoint := getenv("ENDPOINT", "/healthz")
	path := getenv("TARGET_PATH", "/")

	criticalLimitRaw := getenv("CRITICAL_LIMIT_RAW", "0.1")
	criticalLimit, err := strconv.ParseFloat(criticalLimitRaw, 64)
	if err != nil {
		panic(err)
	}

	http.HandleFunc(endpoint, func(w http.ResponseWriter, req *http.Request) {
		v, err := getFreeSpace(path)
		if err != nil {
			log.Println("failed to getFreeSpace: %w", err)
			http.Error(w, "internal error", 502)
			return
		}

		if v < criticalLimit {
			log.Println("not enough disk space", v)
			http.Error(w, "not enough disk space", 500)
			return
		}

		w.Write([]byte("ok"))
	})

	log.Printf("Start listen %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getFreeSpace(path string) (float64, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return 0, fmt.Errorf("failed to syscall.Statfs: %w", err)
	}

	return float64(fs.Bfree) / float64(fs.Blocks), nil
}

func getenv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
		log.Printf("Uses default %s: %s\n", key, defaultValue)
	}
	return value
}
