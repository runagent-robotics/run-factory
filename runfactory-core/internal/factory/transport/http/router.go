package httptransport

import "net/http"

func NewRouter(handler *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /factories", handler.CreateFactory)
	mux.HandleFunc("GET /factories", handler.ListFactories)
	mux.HandleFunc("GET /factories/{factoryID}", handler.GetFactory)
	mux.HandleFunc("PUT /factories/{factoryID}/map", handler.UpdateFactoryMap)
	mux.HandleFunc("POST /factories/{factoryID}/robots", handler.AddRobot)
	mux.HandleFunc("DELETE /factories/{factoryID}/robots/{robotID}", handler.RemoveRobot)
	return mux
}
