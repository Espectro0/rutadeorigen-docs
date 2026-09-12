package task

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Espectro0/rutadeorigen-docs/internal/queue"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	NormalCount = 20
	PoisonCount = 3
	MaxRetries  = 3
)

type Task struct {
	ID        string `json:"id"`
	ForceFail bool   `json:"force_fail"`
}

type Result struct {
	TaskID   string `json:"task_id"`
	Instance string `json:"instance"`
}

const (
	mainExchange   = "tasks"
	mainQueue      = "task_queue"
	mainRoutingKey = "task"

	retryExchange   = "tasks.retry"
	retryQueue      = "task_queue.retry"
	retryRoutingKey = "task.retry"
	retryTTLMs      = 3000

	poisonExchange   = "tasks.poison"
	poisonQueue      = "task_queue.poison"
	poisonRoutingKey = "task.poison"

	resultsQueue = "results_queue"

	retryHeader = "x-retry-count"
)

func SetupTopology(conn *queue.Connection) error {
	if err := conn.DeclareExchange(mainExchange, "direct"); err != nil {
		return fmt.Errorf("Failed to declare main exchange: %v", err)
	}

	if err := conn.DeclareExchange(retryExchange, "direct"); err != nil {
		return fmt.Errorf("Failed to declare retry exchange: %v", err)
	}

	if err := conn.DeclareExchange(poisonExchange, "direct"); err != nil {
		return fmt.Errorf("Failed to declare poison exchange: %v", err)
	}

	if _, err := conn.DeclareQueue(mainQueue, nil); err != nil {
		return fmt.Errorf("Failed to declare main queue: %v", err)
	}

	if err := conn.BindQueue(mainQueue, mainRoutingKey, mainExchange); err != nil {
		return fmt.Errorf("Failed to bind main queue: %v", err)
	}

	retryArgs := amqp.Table{
		"x-message-ttl":             int32(retryTTLMs),
		"x-dead-letter-exchange":    mainExchange,
		"x-dead-letter-routing-key": mainRoutingKey,
	}

	if _, err := conn.DeclareQueue(retryQueue, retryArgs); err != nil {
		return fmt.Errorf("Failed to declare retry queue: %v", err)
	}

	if err := conn.BindQueue(retryQueue, retryRoutingKey, retryExchange); err != nil {
		return fmt.Errorf("Failed to bind retry queue: %v", err)
	}

	if _, err := conn.DeclareQueue(poisonQueue, nil); err != nil {
		return fmt.Errorf("Failed to declare poison queue: %v", err)
	}

	if err := conn.BindQueue(poisonQueue, poisonRoutingKey, poisonExchange); err != nil {
		return fmt.Errorf("Failed to bind poison queue: %v", err)
	}

	if _, err := conn.DeclareQueue(resultsQueue, nil); err != nil {
		return fmt.Errorf("Failed to declare results queue: %v", err)
	}

	return nil
}

func PublishSampleTasks(conn *queue.Connection) error {
	tasks := make([]Task, 0, NormalCount+PoisonCount)
	for i := 1; i <= NormalCount; i++ {
		tasks = append(tasks, Task{
			ID:        fmt.Sprintf("normal-%d", i),
			ForceFail: false,
		})
	}

	for i := 1; i <= PoisonCount; i++ {
		tasks = append(tasks, Task{
			ID:        fmt.Sprintf("poison-%d", i),
			ForceFail: true,
		})
	}

	rand.Shuffle(len(tasks), func(i, j int) {
		tasks[i], tasks[j] = tasks[j], tasks[i]
	})

	for _, t := range tasks {
		body, err := json.Marshal(t)
		if err != nil {
			return fmt.Errorf("Failed to marshal task: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = conn.Publish(ctx, mainExchange, mainRoutingKey, amqp.Table{
			retryHeader: int32(0),
		}, body)

		cancel()
		if err != nil {
			return fmt.Errorf("Failed to publish task: %v", err)
		}

		log.Printf("[x] Published task: %s [force fail: %v]", t.ID, t.ForceFail)
	}
	log.Printf("Total Sended: %d Normal Tasks & %d Poison Tasks", NormalCount, PoisonCount)
	return nil
}

func StartWorker(conn *queue.Connection, instance string) error {
	if err := conn.Channel.Qos(1, 0, false); err != nil {
		return fmt.Errorf("Failed to set QoS: %v", err)
	}

	msgs, err := conn.Channel.Consume(mainQueue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Failed to register a consumer: %v", err)
	}

	log.Printf("Worker %s started. Waiting for messages...", instance)
	for d := range msgs {
		var t Task
		if err := json.Unmarshal(d.Body, &t); err != nil {
			log.Printf("Failed to unmarshal task: %v", err)
			d.Ack(false)
			continue
		}

		retryCount := headerInt32(d.Headers, retryHeader)
		log.Printf("[%s] Received: %s [Try %d]", instance, t.ID, retryCount+1)

		time.Sleep(300 * time.Millisecond)
		if !t.ForceFail {
			log.Printf("[%s] Processed successfully: %s", instance, t.ID)
			publishResult(conn, t.ID, instance)
			d.Ack(false)
			continue
		}

		if retryCount < MaxRetries {
			log.Printf("[%s] Failed: %s -> To Retries [%d/%d]", instance, t.ID, retryCount+1, MaxRetries)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := conn.Publish(ctx, retryExchange, retryRoutingKey, amqp.Table{retryHeader: retryCount + 1}, d.Body)
			cancel()
			if err != nil {
				log.Printf("Failed to re-queue %s: %v", t.ID, err)
				d.Nack(false, true)
				continue
			}
			d.Ack(false)
			continue
		}

		log.Printf("[%s] Definitive Failure %s -> Poison Queue", instance, t.ID)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = conn.Publish(ctx, poisonExchange, poisonRoutingKey, nil, d.Body)
		cancel()
		if err != nil {
			log.Printf("Failed to re-queue to Poison Queue %s: %v", t.ID, err)
			d.Nack(false, true)
			continue
		}
		publishResult(conn, t.ID, instance)
		d.Ack(false)
	}
	return nil
}

func publishResult(conn *queue.Connection, taskID, instance string) {
	r := Result{TaskID: taskID, Instance: instance}
	body, err := json.Marshal(r)
	if err != nil {
		log.Printf("Failed to marshal result for %s: %v", taskID, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.Publish(ctx, "", resultsQueue, nil, body); err != nil {
		log.Printf("Failed to publish result for %s: %v", taskID, err)
	}
}

func headerInt32(headers amqp.Table, key string) int32 {
	if headers == nil {
		return 0
	}
	switch v := headers[key].(type) {
	case int32:
		return v
	case int64:
		return int32(v)
	case int:
		return int32(v)
	}
	return 0
}
