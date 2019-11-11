package main

import (
	"github.com/nats-io/nats.go"
	"github.com/nfrush/garlic/consumer"
	"github.com/nfrush/garlic/eventlistener"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go eventlistener.Start("/var/run/salt/master/master_event_pub.ipc", nats.DefaultURL)
	wg.Add(1)
	go consumer.Start(nats.DefaultURL)
	wg.Wait()
}
