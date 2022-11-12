package glog2kafka

import (
	"errors"
	"github.com/Shopify/sarama"
	"io"
)

type KafkaEndPoint struct {
	addressList []string
	topic       string
	Cfg         *sarama.Config
	client      sarama.Client
	producer    sarama.AsyncProducer
}

var _ io.Writer = &KafkaEndPoint{}

func NewKafkaEndpoint(addresses []string, topic string, cfg *sarama.Config) (*KafkaEndPoint, error) {
	p := &KafkaEndPoint{
		addressList: addresses,
		topic:       topic,
		Cfg:         cfg,
	}
	err := p.Init()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *KafkaEndPoint) Init() error {
	var err error
	var client sarama.Client
	client, err = sarama.NewClient(p.addressList, p.Cfg)
	if err != nil {
		println(err)
		return err
	}
	p.client = client

	var producer sarama.AsyncProducer
	producer, err = sarama.NewAsyncProducerFromClient(p.client)
	if err != nil {
		println("create kafka async producer failed: ", err)
		return err
	}
	p.producer = producer

	// Track errors
	go func() {
		for errMsg := range p.producer.Errors() {
			println(errMsg.Error(), errMsg.Msg.Value.Length())
		}
	}()

	return nil
}

func (p *KafkaEndPoint) Ping() error {
	return p.client.RefreshMetadata()
}

func (p *KafkaEndPoint) produceMsg(topic string, partition int32, msgByte []byte) error {
	if p.producer == nil {
		return errors.New("please use Init before produceMsg()")
	}
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: partition,
		Value:     sarama.StringEncoder(msgByte),
	}

	p.producer.Input() <- msg
	return nil
}

func (p *KafkaEndPoint) Write(msgByte []byte) (n int, err error) {
	return 0, p.produceMsg(p.topic, 0, msgByte)
}

func (p *KafkaEndPoint) Close() {
	if p.producer != nil {
		_ = p.producer.Close()
	}
	if p.client != nil {
		_ = p.client.Close()
	}
}
