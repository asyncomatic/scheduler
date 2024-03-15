//  Copyright (c) 2024 JC Cormier
//  All rights reserved.
//  SPDX-License-Identifier: MIT
//  For full license text, see LICENSE file in the repo root or https://opensource.org/licenses/MIT

package kafka

import (
	"github.com/caitlinelfring/go-env-default"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaOptions struct {
	servers string
	client  string
	ackmode string
}

func NewKafkaOptions() *KafkaOptions {
	return &KafkaOptions{
		servers: env.GetDefault("KAFKA_BOOTSTRAP_SERVERS", "localhost:9094"),
		client:  env.GetDefault("KAFKA_CLIENT_ID", "devcloud"),
		ackmode: env.GetDefault("KAFKA_ACK_MODE", "all")}
}

func (k *KafkaOptions) ConfigMap() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": k.servers,
		"client.id":         k.client,
		"acks":              k.ackmode}
}
