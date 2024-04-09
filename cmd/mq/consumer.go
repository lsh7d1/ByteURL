package main

import (
	"context"
	"encoding/json"
	"fmt"

	"byteurl/dal/dao"
	"byteurl/dal/model"
	"byteurl/mq"
)

func main() {
	reader := mq.MqReader
	for {
		msg, err := reader.FetchMessage(context.TODO())
		if err != nil {
			break
		}
		// 写MySQL
		s := new(model.Short)
		_ = json.Unmarshal([]byte(msg.Value), s)
		fmt.Printf("Short entity: %#v\n", s)
		fmt.Printf("message at offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))
		var _ = msg
		// fail continue
		if err := dao.InsertShortURL(context.TODO(), s); err != nil {
			continue
		}
		reader.CommitMessages(context.Background(), msg)
	}
	if err := reader.Close(); err != nil {
		fmt.Printf("reader.Close failed, err: %v", err)
	}
}
