package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL     = "http://localhost:3000/auditlog"
	projectID   = "dev"
	environment = "dev"
	apiToken    = "dev"
)

type actor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type grupo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type target struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type evento struct {
	Action   string `json:"action"`
	TeamID   string `json:"teamId"`
	Group    grupo  `json:"group"`
	Crud     string `json:"crud"`
	Created  string `json:"created"`
	SourceIP string `json:"source_ip"`
	Actor    actor  `json:"actor"`
	Target   target `json:"target"`
}

func main() {
	if err := publicarEvento(); err != nil {
		fmt.Println("error publicando el evento:", err)
		return
	}

	time.Sleep(1 * time.Second) // le da tiempo a Retraced a indexar el evento

	if err := consultarEventos(); err != nil {
		fmt.Println("error consultando eventos:", err)
	}
}

// publicarEvento simula lo que haría el backend real: emitir un evento de
// auditoría cuando se modifica un lote.
func publicarEvento() error {
	ev := evento{
		Action:   "lote.actualizado",
		TeamID:   "rutadeorigen",
		Group:    grupo{ID: environment, Name: "Ruta de Origen"},
		Crud:     "u",
		Created:  time.Now().UTC().Format(time.RFC3339),
		SourceIP: "127.0.0.1",
		Actor:    actor{ID: "juanes@rutadeorigen.co", Name: "Juanes"},
		Target:   target{ID: "lote-001", Name: "Lote 001", Type: "Lote"},
	}

	body, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/publisher/v1/project/%s/event", baseURL, projectID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token="+apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("evento publicado   -> %d: %s\n", resp.StatusCode, respBody)
	return nil
}

// consultarEventos trae de vuelta los últimos eventos registrados, tal como
// lo haría una pantalla de historial de cambios.
func consultarEventos() error {
	url := fmt.Sprintf("%s/admin/v1/project/%s/events/search?environment_id=%s", baseURL, projectID, environment)

	body, err := json.Marshal(map[string]any{
		"query": map[string]any{
			"length": 10,
			"offset": 0,
		},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token="+apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("eventos consultados -> %d: %s\n", resp.StatusCode, respBody)
	return nil
}
