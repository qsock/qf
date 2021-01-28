package ka

type Kafka struct {
	producer  *Producer
	consumers map[string]*Consumer
	cfg       *Config
}

type Config struct {
	// broker的集群地址
	Brokers   []string          `toml:"brokers"`
	Consumers []*ConsumerConfig `toml:"consumers"`
}

func NewPlat(consumer *ConsumerConfig) (*Kafka, error) {
	cfg := new(Config)
	cfg.Brokers = consumer.Brokers
	if consumer != nil {
		cfg.Consumers = []*ConsumerConfig{consumer}
	}
	return New(cfg)
}

func New(cfg *Config) (*Kafka, error) {
	var (
		err error
		ka  = new(Kafka)
	)
	ka.producer, err = NewProducer(cfg.Brokers)
	if err != nil {
		return nil, err
	}
	if len(cfg.Consumers) == 0 {
		return ka, nil
	}
	ka.cfg = cfg
	ka.consumers = make(map[string]*Consumer)
	for _, consumerCfg := range cfg.Consumers {
		consumerCfg.Brokers = cfg.Brokers
		consumer, err := NewConsumer(consumerCfg)
		if err != nil {
			return nil, err
		}
		ka.consumers[consumerCfg.Topic] = consumer
	}
	return ka, nil
}

func (k *Kafka) Consumer(topic string) *Consumer {
	return k.consumers[topic]
}

func (k *Kafka) Consumers() map[string]*Consumer {
	return k.consumers
}

func (k *Kafka) Producer() *Producer {
	return k.producer
}

func (k *Kafka) Push(topic, eventType string, event interface{}) {
	k.producer.PushEvent(topic, eventType, event)
}

func (k *Kafka) ConsumeTopic(topic string, handler HandleConsumerMsgFunc) error {
	consumer := k.Consumer(topic)
	if consumer == nil {
		return ErrConsumerNotFound
	}
	consumer.HandleMsg(handler)
	return nil
}

func (k *Kafka) ConsumeTopicEvent(topic string, eventType string, handler HandleConsumerFunc) error {
	consumer := k.Consumer(topic)
	if consumer == nil {
		return ErrConsumerNotFound
	}
	consumer.HandleEvent(eventType, handler)
	return nil
}

// 如果在消费组中只有一个consumer
// 可以默认不传topic
func (k *Kafka) Consume(eventType string, handler HandleConsumerFunc) error {
	var consumer *Consumer
	for _, v := range k.consumers {
		consumer = v
	}
	if consumer == nil {
		return ErrConsumerNotFound
	}
	consumer.HandleEvent(eventType, handler)
	return nil
}

func (k *Kafka) Start() {
	for _, consumer := range k.consumers {
		consumer.Run()
	}
}

func (k *Kafka) Close() error {
	err := k.producer.Close()
	for _, consumer := range k.consumers {
		_ = consumer.Close()
	}
	return err
}
