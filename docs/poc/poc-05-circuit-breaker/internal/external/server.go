package external

import (
	"net/http"
	"sync/atomic"
)

// NewFlakyServer arma un servidor HTTP de prueba que simula un servicio
// externo en tres fases: responde bien las primeras `healthyBefore`
// solicitudes que le llegan, falla las siguientes `failCount`, y luego se
// "recupera" y vuelve a responder bien indefinidamente.
func NewFlakyServer(addr string, healthyBefore, failCount int) *http.Server {
	var count int64

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt64(&count, 1)

		if n > int64(healthyBefore) && n <= int64(healthyBefore+failCount) {
			http.Error(w, "servicio externo caído", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
