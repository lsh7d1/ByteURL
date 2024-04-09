package main

import (
	"context"
	"fmt"

	"byteurl/mq"
	"byteurl/service"
)

func main() {
	shorturl, err := service.GenShortURL(context.TODO(), "http://www.bilibili.com")
	if err != nil {
		fmt.Printf("service.GenShortURL failed, err: %v", err)
		return
	}
	fmt.Println("shorturl:", shorturl)
	mq.MqWriter.Close()
}
