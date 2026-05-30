package main

import (
	"fmt"
	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/repository/memory"
	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/service"
	httptransport "github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/transport/http"
	"github.com/runagent-robotics/run-factory/runfactory-core/internal/platform/config"
	"net/http"
)

func main() {
	cfg := config.Load()
	addr := cfg.Address()
	fmt.Printf("runfactory-core listening on %s\n", addr)
	if err := http.ListenAndServe(addr, newServer()); err != nil {
		panic(err)
	}
}

func newServer() http.Handler {
	repo := memory.NewFactoryStore()
	svc := service.NewFactoryService(repo)
	handler := httptransport.NewHandler(svc)
	return httptransport.NewRouter(handler)
}
