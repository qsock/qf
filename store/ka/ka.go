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
	HandleErrorFunc   func(error)
	HandleSucceedFunc func(*sarama.ProducerMessage)
	// 如果返回error，则这一条消费不会被mark消费成功
	HandleConsumerFunc func(e *Event) error
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
	wg     *sync.WaitGroup

	handlers    map[string]HandleConsumerFunc
	handlerLock sync.RWMutex

	cfg *ConsumerConfig
}

type Event struct {
	Time     string          `json:"Time,omitempty"`
	Hostname string          `json:"Hostname,omitempty"`
	From     string          `json:"From,omitempty"`
	Type     string          `json:"Type,omitempty"`
	Data     json.RawMessage `json:"Data,omitempty"`
}

type eventProducer struct {
	Time     string      `json:"Time,omitempty"`
	From     string      `json:"From,omitempty"`
	Hostname string      `json:"Hostname,omitempty"`
	Type     string      `json:"Type,omitempty"`
	Data     interface{} `json:"Data,omitempty"`
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
