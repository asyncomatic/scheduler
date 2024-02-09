package internal

import (
	"encoding/json"
	"fmt"
	kafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

type KafkaQueue struct {
	producer *kafka.Producer
}

func NewKafkaQueue() *KafkaQueue {

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"client.id":         os.Getenv("KAFKA_CLIENT_ID"),
		//"bootstrap.servers": "localhost:9094",
		//"client.id":         "devcloud",
		"acks": "all"}) // "all"

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	return &KafkaQueue{p}
}

func (k *KafkaQueue) Write(test Test) error {
	topic := test.Queue
	msg, err := json.Marshal(test)
	if err != nil {

	}

	err = k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg}, nil)

	return err
}
