package handlers

import (
	"context"
	"encoding/json"
	"log"
	"smarthome/db"
    "smarthome/models"
    "time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Event описывает сообщение из RabbitMQ
type Event struct {
    EventType string    `json:"event_type"`
    Sensor    models.Sensor    `json:"sensor"`
    Timestamp time.Time `json:"timestamp"`
}

// RunRabbitMQConsumer запускает consumer для очереди sensor_events
func RunRabbitMQConsumer(ctx context.Context, amqpURL string, database *db.DB) error {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"sensor_events", // имя очереди
		true,            // durable
		false,           // auto-delete
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	// Горутина для обработки сообщений
	go func() {
		for {
			select {
			case d := <-msgs:
				var e Event
				if err := json.Unmarshal(d.Body, &e); err != nil {
					log.Printf("Error decoding event: %v\n", err)
					continue
				}
				log.Printf("[RabbitMQ] Received event: %+v\n", e)

				// Сохранение в БД
				switch e.EventType {
				case "sensor_created":
					// Пример: создаем новый сенсор на основе данных из события
					 sensor := e.Sensor
					_, err := database.CreateSensor(ctx, models.SensorCreate{
						Name:     sensor.Name,
						Type:     sensor.Type,
						Location: sensor.Location,
						Unit:     sensor.Unit,
					})
					if err != nil {
						log.Printf("Error saving sensor to DB: %v", err)
					} else {
						log.Printf("Sensor saved to DB from event %s", e.EventType)
					}
				default:
					log.Printf("Unknown event type: %s", e.EventType)
				}

			case <-ctx.Done():
				log.Println("[RabbitMQ] Consumer shutting down")
				return
			}
		}
	}()

	log.Println("[RabbitMQ] Consumer started, waiting for messages...")
	<-ctx.Done()
	return nil
}
