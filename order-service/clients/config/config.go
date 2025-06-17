package config

import "github.com/parnurzeal/gorequest"

type CLientConfig struct {
	client       *gorequest.SuperAgent
	baseURL      string
	signatureKey string
}

type IClientConfig interface {
	Client() *gorequest.SuperAgent
	BaseURL() string
	SignatureKey() string
}

type Option func(*CLientConfig)

func NewClientConfig(options ...Option) IClientConfig {
	clientConfig := &CLientConfig{
		client: gorequest.New().Set("Content-Type", "application/json").Set("Accept", "application/json"),
	}
	for _, option := range options {
		option(clientConfig)
	}
	return clientConfig
}

func (c *CLientConfig) Client() *gorequest.SuperAgent {
	return c.client
}

func (c *CLientConfig) BaseURL() string {
	return c.baseURL
}

func (c *CLientConfig) SignatureKey() string {
	return c.signatureKey
}

func WithBaseURL(baseURL string) Option {
	return func(c *CLientConfig) {

		c.baseURL = baseURL
	}
}

func WithSignatureKey(signatureKey string) Option {
	return func(c *CLientConfig) {
		c.signatureKey = signatureKey
	}
}
