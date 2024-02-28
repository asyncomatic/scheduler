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
