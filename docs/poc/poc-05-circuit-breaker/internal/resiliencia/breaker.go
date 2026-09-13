package resiliencia

import (
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

// NewBreaker configura un circuit breaker para proteger llamadas a un
// servicio externo: se abre tras 2 fallos consecutivos, permanece abierto
// 2 segundos, y luego pasa a "medio abierto" para reintentar con una sola
// solicitud de prueba.
func NewBreaker(nombre string) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        nombre,
		MaxRequests: 1,
		Timeout:     2 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 2
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("[breaker] %s pasó de %s a %s\n", name, from, to)
		},
	}

	return gobreaker.NewCircuitBreaker(settings)
}
