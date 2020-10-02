package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type RMQMessage struct {
	Timestamp int64             `json:"ts"`
	Value     float64           `json:"value"`
	Topic     string            `json:"topic"`
	Scope     string            `json:"scope"`
	Tags      map[string]string `json:"tags"`
}

func getCredential(host string, id string, pw string) string {
	return fmt.Sprintf("amqp://%s:%s@%s", id, pw, host)
}

func main() {
	value := flag.Float64("value", 0., "Value to rain_gauge")
	flag.Parse()

	conn, _ := amqp.Dial(getCredential("localhost:5672", "worker", "worker"))
	ch, _ := conn.Channel()

	var msg RMQMessage
	msg.Timestamp = time.Now().UnixNano()
	msg.Topic = "env.rain_gauge"
	msg.Value = *value
	msg.Scope = "node"

	jsonString, _ := json.Marshal(msg)

	_ = ch.Publish(
		"messages", // exchange
		"",         // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonString,
		})
}
