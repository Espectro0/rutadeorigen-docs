package main

import (
	"log"
	"sync"

	"github.com/Espectro0/rutadeorigen-docs/internal/queue"
	"github.com/Espectro0/rutadeorigen-docs/internal/task"
)

func main() {
	const url = "amqp://guest:guest@localhost:5672/"
	instances := []string{"worker-1", "worker-2"}

	conn, err := queue.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := task.SetupTopology(conn); err != nil {
		log.Fatalf("Failed to setup topology: %v", err)
	}

	if err := task.PublishSampleTasks(conn); err != nil {
		log.Fatalf("Failed to publish sample tasks: %v", err)
	}

	var wg sync.WaitGroup
	for _, instance := range instances {
		wg.Add(1)
		go func() {
			defer wg.Done()

			workerConn, err := queue.Connect(url)
			if err != nil {
				log.Fatalf("Failed to connect to RabbitMQ for worker %s: %v", instance, err)
			}
			defer workerConn.Close()

			if err := task.StartWorker(workerConn, instance); err != nil {
				log.Fatalf("Worker %s encountered an error: %v", instance, err)
			}
		}()
	}

	wg.Wait()
}
