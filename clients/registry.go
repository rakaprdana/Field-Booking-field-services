package clients

import (
	clientConfig "field-services/clients/config"
	client "field-services/clients/user"
	"field-services/config"
)

type ClientRegistry struct{}

// GetUser implements [IClientRegistry].
func (c *ClientRegistry) GetUser() client.IUserClient {
	return client.NewUserClient(
		clientConfig.NewClientConfig(
			clientConfig.WithBaseURL(config.Config.InternalService.User.Host),
			clientConfig.WithSignatureKey(config.Config.SignatureKey),
		),
	)
}

type IClientRegistry interface {
	GetUser() client.IUserClient
}

func NewClientRegistry() IClientRegistry {
	return &ClientRegistry{}
}
