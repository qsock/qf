package ka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"sync"
)

const (
	gTestPrefix = "test_consumer_"
	gMsgDone    = "done"
)

var (
	ver = sarama.V0_10_2_0
)

var (
	ErrConsumerNotFound = errors.New("consumer not found")
)

type (
	HandleErrorFunc   func(error)
	HandleSucceedFunc func(*sarama.ProducerMessage)
	// 如果返回error，则这一条消费不会被mark消费成功
	HandleConsumerFunc    func(e *Event) error
	HandleConsumerMsgFunc func(message *sarama.ConsumerMessage)
)

type Producer struct {
	config   *sarama.Config
	producer sarama.AsyncProducer
	e        HandleErrorFunc
	s        HandleSucceedFunc
}

type Consumer struct {
	run bool

	ctx    context.Context
	cancel context.CancelFunc

	group  sarama.ConsumerGroup
	config *sarama.Config
	e      HandleErrorFunc
	wg     *sync.WaitGroup

	cfg   *ConsumerConfig
	topic string

	consumerMsgHandler HandleConsumerMsgFunc
	consumerHandlers   map[string]HandleConsumerFunc
	handlerLock        sync.RWMutex
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

func SetVersion(v sarama.KafkaVersion) {
	ver = v
}

func GetVersion() sarama.KafkaVersion {
	return ver
}
