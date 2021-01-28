package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/qsock/qf/store/ka"
	"time"
)

var (
	cfg = &ka.Config{
		Brokers: []string{"127.0.0.1:9092"},
		Consumers: []*ka.ConsumerConfig{
			{
				Topic: "test_consumer",
				Group: "test_group",
			},
		},
	}
)

func main() {
	const (
		type1 = "type1"
		type2 = "type2"
	)
	ka.SetVersion(sarama.V0_10_2_0)
	kafka, err := ka.New(cfg)
	defer kafka.Close()

	if err != nil {
		panic(err)
	}

	{
		_ = kafka.Consume(type1, func(e *ka.Event) error {
			fmt.Println(e.Type, string(e.Data))
			return nil
		})
		_ = kafka.Consume(type2, func(e *ka.Event) error {
			fmt.Println(e.Type, string(e.Data))
			return nil
		})
	}

	kafka.Start()
	kafka.Producer().HandleError(func(e error) {
		fmt.Println(e)
	})
	kafka.Consumer("test_consumer").HandleError(func(e error) {
		fmt.Println("consume", e)
	})
	var i uint64
	for {
		i++
		if i == 10 {
			return
		}
		time.Sleep(time.Millisecond * 100)
		kafka.Push("test_consumer", type1, map[string]interface{}{"kolo1": time.Now().String()})
		kafka.Push("test_consumer", type2, map[string]interface{}{"tick2": time.Now().String()})
	}
	fmt.Println("done")
}
