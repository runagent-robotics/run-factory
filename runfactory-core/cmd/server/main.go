package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	addr := ":8080"
	fmt.Printf("runfactory-core listening on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
