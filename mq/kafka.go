package mq

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

const (
	topic     = "byteurl"
	partition = 0
)

var MqWriter = &kafka.Writer{
	Addr:         kafka.TCP("127.0.0.1:9092"),
	Topic:        topic,
	Balancer:     &kafka.LeastBytes{}, // 指定分区的balancer模式为最小字节分布
	RequiredAcks: kafka.RequireAll,    // ack模式
	Async:        false,               // 同步
}

var MqReader = kafka.NewReader(kafka.ReaderConfig{
	Brokers:   []string{"127.0.0.1:9092"},
	Topic:     topic,
	Partition: partition,
	MaxBytes:  10e6,
})

type MqEntity struct {
	Key   []byte
	Value []byte
}

func WriteMessage(ctx context.Context, msgs []MqEntity) error {
	kmsgs := make([]kafka.Message, 0, len(msgs))
	for _, msg := range msgs {
		kmsgs = append(kmsgs, kafka.Message{
			Key:   msg.Key,
			Value: msg.Value,
		})
	}
	fmt.Printf("[mq] kmsgs: %v\n", kmsgs)
	if err := MqWriter.WriteMessages(ctx, kmsgs...); err != nil {
		return fmt.Errorf("write kmsgs: %#v, failed, err: %v", kmsgs, err)
	}
	return nil
}
