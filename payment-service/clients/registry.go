package clients

import (
	"payment-service/clients/config"
	clients "payment-service/clients/user"
	paymentConfig "payment-service/config"
)

type ClientRegistry struct{}

type IClientRegistry interface {
	GetUser() clients.IUserClient
}

func NewClientRegistry() IClientRegistry {
	return &ClientRegistry{}
}

func (c *ClientRegistry) GetUser() clients.IUserClient {
	return clients.NewUserClient(
		config.NewClientConfig(
			config.WithBaseURL(paymentConfig.Config.InternalService.User.Host),
			config.WithSignatureKey(paymentConfig.Config.InternalService.User.SignatureKey),
		))
}
