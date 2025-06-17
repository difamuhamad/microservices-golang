package clients

import (
	"order-service/clients/config"
	fieldClients "order-service/clients/field"
	paymentClients "order-service/clients/payment"
	clients "order-service/clients/user"
	configApp "order-service/config"
)

type ClientRegistry struct{}

type IClientRegistry interface {
	GetUser() clients.IUserClient
	GetPayment() paymentClients.IPaymentClient
	GetField() fieldClients.IFieldClient
}

func NewClientRegistry() IClientRegistry {
	return &ClientRegistry{}
}

// Get User
func (c *ClientRegistry) GetUser() clients.IUserClient {
	return clients.NewUserClient(
		config.NewClientConfig(
			config.WithBaseURL(configApp.Config.InternalService.User.Host),
			config.WithSignatureKey(configApp.Config.InternalService.User.SignatureKey),
		))
}

// Get Payment
func (c *ClientRegistry) GetPayment() paymentClients.IPaymentClient {
	return paymentClients.NewPaymentClient(
		config.NewClientConfig(
			config.WithBaseURL(configApp.Config.InternalService.Payment.Host),
			config.WithSignatureKey(configApp.Config.InternalService.Payment.SignatureKey),
		))
}

// Get Field
func (c *ClientRegistry) GetField() fieldClients.IFieldClient {
	return fieldClients.NewFieldClient(
		config.NewClientConfig(
			config.WithBaseURL(configApp.Config.InternalService.Field.Host),
			config.WithSignatureKey(configApp.Config.InternalService.Field.SignatureKey),
		))
}
