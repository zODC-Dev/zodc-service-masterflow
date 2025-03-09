package externals

import (
	"github.com/go-resty/resty/v2"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/results"
)

type UserAPI struct{}

func NewUserAPI() *UserAPI {
	return &UserAPI{}
}

func (u *UserAPI) FindUsersByUserIds(userIds []int32) (results.UserApiResult, error) {

	var body struct {
		UserIds []int32 `json:"userIds"`
	}
	body.UserIds = userIds

	result := results.UserApiResult{}

	client := resty.New()

	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&result).
		Post("https://zodc-api.thanhf.dev/auth/internal/users")

	if err != nil {
		return result, err
	}

	return result, nil

}
