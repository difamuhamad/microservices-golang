package kafka

import (
	"order-service/config"
	"order-service/controllers/kafka"
	kafkaPayment "order-service/controllers/kafka/payment"

	"golang.org/x/exp/slices"
)

type Kafka struct {
	consumer *ConsumerGroup
	kafka    kafka.IKafkaRegistry
}

type IKafka interface {
	Register()
}

func NewKafkaConsumer(consumer *ConsumerGroup, kafka kafka.IKafkaRegistry) IKafka {
	return &Kafka{consumer: consumer, kafka: kafka}
}

func (k *Kafka) Register() {
	k.paymentHandler()
}

func (k *Kafka) paymentHandler() {
	// check if array contain the Topic
	if slices.Contains(config.Config.Kafka.Topics, kafkaPayment.PaymentTopic) {
		k.consumer.RegisterHandler(kafkaPayment.PaymentTopic, k.kafka.GetPayment().HandlePayment)
	}
}
