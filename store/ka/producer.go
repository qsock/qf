package ka

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"time"
)

func NewProducer(brokers []string, name ...string) error {
	cfg := sarama.NewConfig()
	cfg.Net.KeepAlive = 60 * time.Second
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Version = sarama.V0_10_2_0
	cfg.Producer.Flush.Frequency = time.Second
	cfg.Producer.Flush.MaxMessages = 10
	return NewProducerWithCfg(brokers, cfg, name...)
}

func NewProducerWithInterceptor(brokers []string, interceptor sarama.ProducerInterceptor, name ...string) error {
	cfg := sarama.NewConfig()
	cfg.Net.KeepAlive = 60 * time.Second
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Version = sarama.V0_10_2_0
	cfg.Producer.Flush.Frequency = time.Second
	cfg.Producer.Flush.MaxMessages = 10
	if interceptor != nil {
		cfg.Producer.Interceptors = []sarama.ProducerInterceptor{interceptor}
	}
	return NewProducerWithCfg(brokers, cfg, name...)
}

func NewProducerWithCfg(brokers []string, cfg *sarama.Config, name ...string) error {
	producer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return err
	}
	p := new(Producer)
	p.producer = producer
	p.config = cfg
	go p.handle()
	producerName := gDefaultName
	if len(name) > 0 {
		producerName = name[0]
	}
	_, ok := producers[producerName]
	if ok {
		return errors.New("producer exists")
	}
	producers[producerName] = p
	return nil
}

func (p *Producer) handle() {
	for {
		select {
		case err, ok := <-p.producer.Errors():
			{
				if ok {
					if p.e != nil {
						p.e(err)
					}
				}
			}
		case msg, ok := <-p.producer.Successes():
			{
				if ok {
					if p.s != nil {
						p.s(msg)
					}
				}
			}
		}
	}
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

func (p *Producer) HandleError(e HandleErrorFunc) {
	p.e = e
}

func (p *Producer) HandleSucceed(s HandleSucceedFunc) {
	p.s = s
}

func (p *Producer) Send(topic string, data []byte) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}
	p.producer.Input() <- msg
}

func (p *Producer) TopicEvent(topic, data string, e interface{}) {
	msg := &eventProducer{
		Type:     data,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		From:     localAddr,
		Hostname: hostname,
		Data:     e,
	}
	b, _ := json.Marshal(msg)
	p.Send(topic, b)
}
