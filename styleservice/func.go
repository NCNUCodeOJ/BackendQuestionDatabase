package styleservice

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

// StyleTask is a template for a style task
type StyleTask struct {
	SourceCode   string `json:"source_code"`
	Language     string `json:"language"`
	SubmissionID uint   `json:"submission_id"`
}

var supportedLanguage = map[string]bool{
	"java":    true,
	"python3": true,
}
var conn *amqp.Connection
var channel *amqp.Channel
var queue amqp.Queue

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Setup connect rabbitmq
func Setup() {
	var err error

	if gin.Mode() == "test" {
		return
	}
	if gin.Mode() == "debug" {
		err = godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}

	conn, err = amqp.Dial(os.Getenv("RABBITMQ_HOST"))
	failOnError(err, "Failed to connect to RabbitMQ")
	channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	queue, err = channel.QueueDeclare(
		"program_style", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")
}

// ErrUnsupportedLanguage is returned when the language is not supported
var ErrUnsupportedLanguage = errors.New("unsupported language")

// Validate validates the style task
func (j *StyleTask) Validate() error {
	if _, ok := supportedLanguage[j.Language]; !ok {
		return ErrUnsupportedLanguage
	}
	return nil
}

// Run  Run a new submission
func (j *StyleTask) Run() (err error) {
	var data []byte

	if data, err = json.Marshal(j); err != nil {
		return
	}

	if gin.Mode() == "test" {
		return
	}

	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/json",
			Body:         data,
		},
	)

	return
}
