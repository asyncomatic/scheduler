//  Copyright (c) 2024 JC Cormier
//  All rights reserved.
//  SPDX-License-Identifier: MIT
//  For full license text, see LICENSE file in the repo root or https://opensource.org/licenses/MIT

package queue

import (
	"reflect"
	"scheduler/models"
	"scheduler/queue/kafka"
)

var QueueWriterRegistry = map[string]QueueWriter{
	"_default": kafka.NewQueueWriter(),
	"kafka":    kafka.NewQueueWriter(),
}

type QueueWriter interface {
	Write(models.Job) error
}

func NewQueueWriter(queueType string) QueueWriter {
	return reflect.ValueOf(QueueWriterRegistry[queueType]).Interface().(QueueWriter)
}
