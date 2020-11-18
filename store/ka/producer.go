package ka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/qsock/qf/net/ipaddr"
	"github.com/qsock/qf/qlog"
	"time"
)

var producer sarama.AsyncProducer

func Close() {
	if producer != nil {
		_ = producer.Close()
	}
}

func Init(addrs []string) (err error) {
	return InitWithDuration(addrs, time.Second)
}

func InitWithDuration(addrs []string, d time.Duration) (err error) {
	config := sarama.NewConfig()
	config.Net.KeepAlive = 60 * time.Second
	config.Producer.Return.Successes = false
	config.Producer.Flush.Frequency = d
	config.Producer.Flush.MaxMessages = 10

	producer, err = sarama.NewAsyncProducer(addrs, config)
	if err != nil {
		return err
	}

	if err == nil {
		go func() {
			for err := range producer.Errors() {
				qlog.Get().Logger().Error("kafka||err:%#v", err)
			}
		}()
	}

	return err
}

func Send(topic string, data []byte) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}
	if producer == nil {
		return
	}
	producer.Input() <- msg
}

type E struct {
	Time string      `json:"atime,omitempty"`
	Type string      `json:"atype,omitempty"`
	From string      `json:"afrom,omitempty"`
	Msg  interface{} `json:"msg,omitempty"`
}

type EConsumer struct {
	Time string          `json:"atime,omitempty"`
	Type string          `json:"atype,omitempty"`
	From string          `json:"afrom,omitempty"`
	Msg  json.RawMessage `json:"msg,omitempty"`
}

func TopicEvent(topic, t string, e interface{}) {
	msg := &E{
		Type: t,
		Time: time.Now().Format("2006-01-02 15:04:05"),
		From: ipaddr.GetHostname(),
		Msg:  e,
	}
	b, _ := json.Marshal(msg)
	Send(topic, b)
}
