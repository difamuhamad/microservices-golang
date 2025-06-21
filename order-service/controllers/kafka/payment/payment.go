package kafka

import (
	"context"
	"encoding/json"
	"order-service/common/utils"
	"order-service/domain/dto"
	"order-service/services"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

const PaymentTopic = "payment-service-callback"

type PaymentKafka struct {
	service services.IServiceRegistry
}

type IPaymentKafka interface {
	HandlePayment(context.Context, *sarama.ConsumerMessage) error
}

func NewPaymentKafka(service services.IServiceRegistry) IPaymentKafka {
	return &PaymentKafka{service: service}
}

func (p *PaymentKafka) HandlePayment(ctx context.Context, message *sarama.ConsumerMessage) error {
	defer utils.Recover()
	var body dto.PaymentContent

	err := json.Unmarshal(message.Value, &body)
	if err != nil {
		logrus.Errorf("failes to unmarshal message :%v", err)
		return err
	}

	data := body.Body.Data
	err = p.service.GetOrder().HandlePayment(ctx, &data)
	if err != nil {
		logrus.Errorf("failes to unmarshal message :%v", err)
		return err
	}
	logrus.Infof("success handle payment")
	return nil
}
