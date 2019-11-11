package consumer

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/nats.go"
	"github.com/nfrush/garlic/eventlistener"
	"log"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Start(broker_url string){
	nc, err := nats.Connect(broker_url)
	if err != nil {
		log.Fatal("NATS err", err)
	}
	go nc.Subscribe("events", func(msg *nats.Msg){
		var message eventlistener.Message
		json.Unmarshal(msg.Data, &message)
		fmt.Printf("RECEIVED AT CONSUMER: %v\n", message)
	})
}