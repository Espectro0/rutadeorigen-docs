package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Espectro0/rutadeorigen-docs/internal/external"
	"github.com/Espectro0/rutadeorigen-docs/internal/resiliencia"
)

func main() {
	// Servidor simulado: responde bien las primeras 3 veces que le llegan
	// solicitudes, falla las siguientes 3, y luego se "recupera".
	server := external.NewFlakyServer(":9090", 3, 3)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("servidor de prueba detenido:", err)
		}
	}()
	time.Sleep(200 * time.Millisecond) // le da tiempo al servidor a levantar

	breaker := resiliencia.NewBreaker("servicio-externo-prueba")
	client := resiliencia.NewClient("http://localhost:9090/ping", breaker)

	ctx := context.Background()

	for i := 1; i <= 25; i++ {
		resultado, err := client.Ping(ctx)
		if err != nil {
			fmt.Printf("intento %02d -> error: %v (estado breaker: %s)\n", i, err, breaker.State())
		} else {
			fmt.Printf("intento %02d -> %s (estado breaker: %s)\n", i, resultado, breaker.State())
		}

		time.Sleep(300 * time.Millisecond)
	}
}
