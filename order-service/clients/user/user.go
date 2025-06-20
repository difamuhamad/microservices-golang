package clients

import (
	"context"
	"fmt"
	"net/http"
	"order-service/clients/config"
	"order-service/common/utils"
	configApp "order-service/config"
	"order-service/constants"
	"time"

	"github.com/google/uuid"
)

type UserClient struct {
	client config.IClientConfig
}

type IUserClient interface {
	GetUserByToken(context.Context) (*UserData, error)
	GetUserByUUID(context.Context, uuid.UUID) (*UserData, error)
}

func NewUserClient(client config.IClientConfig) IUserClient {
	return &UserClient{client: client}
}

func (u *UserClient) GetUserByToken(ctx context.Context) (*UserData, error) {
	unixTime := time.Now().Unix()
	generateAPIKey := fmt.Sprintf("%s:%s:%d",
		configApp.Config.AppName,
		u.client.SignatureKey(),
		unixTime,
	)
	apiKey := utils.GenerateSHA256(generateAPIKey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	var response UserResponse
	request := u.client.Client().Clone().
		Set(constants.Authorization, bearerToken).
		Set(constants.XApiKey, apiKey).
		Set(constants.XServiceName, configApp.Config.AppName).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Get(fmt.Sprintf("%s/api/v1/auth/user", u.client.BaseURL()))

	res, _, errs := request.EndStruct(&response)
	if len(errs) > 0 {
		return nil, fmt.Errorf("request failed: %v", errs[0])
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, message: %s",
			res.StatusCode, response.Message)
	}

	return &response.Data, nil
}

func (u *UserClient) GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*UserData, error) {
	unixTime := time.Now().Unix()
	generateAPIKey := fmt.Sprintf("%s:%s:%d",
		configApp.Config.AppName,
		u.client.SignatureKey(),
		unixTime,
	)
	apiKey := utils.GenerateSHA256(generateAPIKey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	var response UserResponse
	request := u.client.Client().Clone().
		Set(constants.Authorization, bearerToken).
		Set(constants.XApiKey, apiKey).
		Set(constants.XServiceName, configApp.Config.AppName).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Get(fmt.Sprintf("%s/api/v1/auth/%s", u.client.BaseURL(), uuid))

	res, _, errs := request.EndStruct(&response)
	if len(errs) > 0 {
		return nil, fmt.Errorf("request failed: %v", errs[0])
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, message: %s",
			res.StatusCode, response.Message)
	}

	return &response.Data, nil
}
