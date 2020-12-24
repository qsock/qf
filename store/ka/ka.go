package ka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"sync"
)

const (
	gDefaultName = "default"
	TestPrefix   = "test_consumer_"
	gMsgDone     = "done"
)

var (
	producers = map[string]*Producer{}
	consumers = map[string]*Consumer{}
)

type (
	HandleErrorFunc    func(error)
	HandleSucceedFunc  func(*sarama.ProducerMessage)
	HandleConsumerFunc func([]byte) error
)

type Producer struct {
	config   *sarama.Config
	producer sarama.AsyncProducer
	e        HandleErrorFunc
	s        HandleSucceedFunc
}

type Consumer struct {
	ctx    context.Context
	cancel context.CancelFunc

	group  sarama.ConsumerGroup
	config *sarama.Config
	e      HandleErrorFunc
	c      HandleConsumerFunc
	wg     *sync.WaitGroup

	cfg *ConsumerConfig
}

type E struct {
	Time     string      `json:"atime,omitempty"`
	Type     string      `json:"atype,omitempty"`
	From     string      `json:"afrom,omitempty"`
	Hostname string      `json:"hostname,omitempty"`
	Msg      interface{} `json:"msg,omitempty"`
}

type EConsumer struct {
	Time     string          `json:"atime,omitempty"`
	Type     string          `json:"atype,omitempty"`
	From     string          `json:"afrom,omitempty"`
	Hostname string          `json:"hostname,omitempty"`
	Msg      json.RawMessage `json:"msg,omitempty"`
}

func GetProducer(name ...string) *Producer {
	producerName := gDefaultName
	if len(name) > 0 {
		producerName = name[0]
	}
	return producers[producerName]
}

func GetConsumer(name ...string) *Consumer {
	consumerName := gDefaultName
	if len(name) > 0 {
		consumerName = name[0]
	}
	return consumers[consumerName]
}

func TopicEvent(topic, t string, e interface{}) {
	GetProducer().TopicEvent(topic, t, e)
}

func Close() error {
	for _, producer := range producers {
		_ = producer.Close()
	}
	for _, consumer := range consumers {
		_ = consumer.Close()
	}
	return nil
}
