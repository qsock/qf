package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/qsock/qf/store/ka"
	"time"
)

var (
	brokers       = []string{"kafka-0.kafka.infra:9092"}
	topic         = "test-topic-1"
	consumerGroup = "test-topic-1-consumer"
)

func main() {
	defer ka.Close()
	// consumer
	{
		c, err := ka.NewConsumer(&ka.ConsumerConfig{Brokers: brokers, Topic: topic, Group: consumerGroup, Workers: 2})
		if err != nil {
			panic(err)
		}
		c.HandleSucceed(func(b []byte) error {
			fmt.Println("consumer000", string(b))
			return nil
		})
		c.HandleError(func(e error) {
			fmt.Println("err consumer000", e)
		})

		c.Run()
	}

	// consumer
	{
		c, err := ka.NewConsumer(&ka.ConsumerConfig{Brokers: brokers, Topic: topic, Group: consumerGroup + "-2", Workers: 2}, "hello")
		if err != nil {
			panic(err)
		}
		c.HandleSucceed(func(b []byte) error {
			fmt.Println("consumer111", string(b))
			return nil
		})
		c.HandleError(func(e error) {
			fmt.Println("err consumer111", e)
		})
		c.Run()
	}

	// producer
	{
		p, _ := ka.NewProducer(brokers)
		p.HandleSucceed(func(s *sarama.ProducerMessage) {
			b, _ := s.Value.Encode()
			fmt.Println("producer", string(b))
		})
		p.HandleError(func(e error) {
			fmt.Println("err producer", e)

		})
		var i uint64
		for {
			i++
			if i == 10 {
				return
			}
			time.Sleep(time.Millisecond * 100)
			ka.TopicEvent(topic, "1111", map[string]interface{}{"k": time.Now().String()})
		}
	}
}
