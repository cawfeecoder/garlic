package eventlistener

import (
	"github.com/cskr/pubsub"
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack"
	"io"
	"log"
	"net"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Message struct {
	Tag string
	Payload map[string]interface{}
}

func listen(r io.Reader, ps *pubsub.PubSub) {
	decoder := msgpack.NewDecoder(r)
	for {
		var message map[string]interface{}
		decoder.Decode(&message)
		body := message["body"].(string)
		var payload map[string]interface{}
		results := strings.SplitN(body, "\n\n", 2)
		tag := results[0]
		_ = msgpack.Unmarshal([]byte(results[1]), &payload)
		m := Message{
			Tag:     tag,
			Payload: payload,
		}
		ps.Pub(m, "master_sock")
	}
}

func Start(sock string, broker_url string) {
	nc, err := nats.Connect(broker_url)
	if err != nil {
		log.Fatal("NATS error", err)
	}
	b, err := net.Dial("unix", sock)
	if err != nil {
		log.Fatal("Dial error", err)
	}
	defer b.Close()

	ps := pubsub.New(10000)

	go listen(b, ps)

	ch := ps.Sub("master_sock")

	for {
		if msg, ok := <- ch; ok {
			b, _ := json.Marshal(&msg)
			nc.Publish("events", b)
		}
	}
}
