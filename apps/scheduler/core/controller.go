package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Job struct {
	ID          uint
	JobName     string
	IsRecurring bool
	NextRunTime int64
	UserID      uint
	Cron        string
	IsCompleted bool
	IsRunning   bool
	IsDisabled  bool
	RetryCount  uint16
}

func getActiveJobs() ([]Job, error) {
	fmt.Println("Reading active jobs...")

	API_BASE := os.Getenv("API_BASE_URL")
	endpoint := API_BASE + "/scheduler/jobs?active=true"
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("x-api-key", os.Getenv("ADMIN_API_KEY"))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		fmt.Println("Failed. Status Code", resp.StatusCode)
		return nil, errors.New("failed to get active jobs")
	}
	fmt.Println("Received all active jobs")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var jobs []Job
	err = json.Unmarshal(body, &jobs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jobs, nil
}

func PublishActiveJobs() {
	jobs, err := getActiveJobs()
	if err != nil {
		fmt.Println("Failed to get active jobs")
		return
	}

	if len(jobs) == 0 {
		fmt.Println("No active jobs.")
		return
	}

	for _, job := range jobs {
		publishJobToQueue(job)
	}
}

func publishJobToQueue(job Job) {
	HOST := os.Getenv("RABBIT_MQ_HOST")
	QUEUE := os.Getenv("JOB_QUEUE_NAME")

	conn, err := amqp.Dial(HOST)
	if err != nil {
		fmt.Printf("%s %s", err, "Failed to connect to RabbitMQ")
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("%s %s", err, "Failed to open a channel")
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		QUEUE, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		fmt.Printf("%s %s", err, "Failed to declare a queue")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jobString, err := json.Marshal(&job)
	if err != nil {
		fmt.Printf("%s %s", err, "Failed to marshal job")
		return
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         jobString,
		})
	if err != nil {
		fmt.Printf("%s %s", err, "Failed to publish a message")
		return
	}
	log.Printf(" [x] Sent %v", job.ID)
}
