//  Copyright (c) 2024 JC Cormier
//  All rights reserved.
//  SPDX-License-Identifier: MIT
//  For full license text, see LICENSE file in the repo root or https://opensource.org/licenses/MIT

package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"scheduler/models"
)

type Writer struct {
	producer *kafka.Producer
}

func NewQueueWriter() *Writer {
	p, err := kafka.NewProducer(NewKafkaOptions().ConfigMap())
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	return &Writer{p}
}

func (k *Writer) Write(job models.Job) error {
	topic := job.Queue
	msg, err := json.Marshal(job)
	if err != nil {

	}

	err = k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg}, nil)

	return err
}
