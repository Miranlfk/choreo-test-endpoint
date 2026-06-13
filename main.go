// Minimal configurable-delay endpoint for perf tests.
//
//	GET /delay?ms=300   -> sleeps 300ms, then 200 {"slept_ms":300}
//
// The delay is per-request via the ?ms= query param (default 200, clamped to
// the 200-400ms range), so you can tune latency from the adaptive script
// without redeploying.
//
// Run locally:   go run main.go            (listens on :8080, or $PORT)
// Deploy:        any container platform (Choreo, Cloud Run, etc.). Co-locate it in
//                the same region as the IS so you measure the injected delay, not WAN.
package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/delay", func(w http.ResponseWriter, r *http.Request) {
		ms := 200
		if v := r.URL.Query().Get("ms"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				ms = n
			}
		}
		if ms < 200 { // floor: never faster than 200ms
			ms = 200
		}
		if ms > 400 { // cap: never slower than 400ms
			ms = 400
		}
		time.Sleep(time.Duration(ms) * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"slept_ms":%d}`, ms)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("delay endpoint listening on :" + port)
	http.ListenAndServe(":"+port, nil)
}
