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

func NewKafkaQueue(opts *KafkaOptions) *KafkaQueue {

	p, err := kafka.NewProducer(opts.ConfigMap())
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	return &KafkaQueue{p}
}

func (k *KafkaQueue) Write(job Job) error {
	topic := job.Queue
	msg, err := json.Marshal(job)
	if err != nil {

	}

	err = k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg}, nil)

	return err
}
