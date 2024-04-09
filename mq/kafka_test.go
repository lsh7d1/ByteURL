package mq

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

func TestWriteByWriter(t *testing.T) {
	// 创建一个writer 向topic发送消息
	w := &kafka.Writer{
		Addr:         kafka.TCP("127.0.0.1:9092"),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{}, // 指定分区的balancer模式为最小字节分布
		RequiredAcks: kafka.RequireAll,    // ack模式
		Async:        false,               // 异步
	}

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("One"),
			Value: []byte("111"),
		},
		kafka.Message{
			Key:   []byte("Two"),
			Value: []byte("222"),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	// if err := w.Close(); err != nil {
	// 	log.Fatal("failed to close writer:", err)
	// }
}

// TestReadByReader 通过Reader接收消息
func TestReadByReader(t *testing.T) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"127.0.0.1:9092"},
		Topic:     topic,
		Partition: partition,
		MaxBytes:  10e6,
	})

	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		r.CommitMessages(context.Background(), m)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

// TestWriteByConn 基于Conn发送消息
func TestWriteByConn(t *testing.T) {
	// 连接至Kafka的leader节点，需指定要连接的topic和partition
	conn, err := kafka.DialLeader(context.Background(), "tcp", "127.0.0.1:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	// 设置发送消息的超时时间
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// 发送消息
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	// 关闭连接
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

// TestReadByConn 基于Conn接收消息
func TestReadByConn(t *testing.T) {
	// 连接至Kafka的leader节点，需指定要连接的topic和partition
	conn, err := kafka.DialLeader(context.Background(), "tcp", "127.0.0.1:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	// 设置读取超时时间
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	// 读取一批消息，得到的batch是一系列消息的迭代器
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	// 遍历读取消息
	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

	// 关闭batch
	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	// 关闭连接
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}
