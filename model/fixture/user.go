package fixture

import (
	"fmt"
	"github.com/sentrionic/valkyrie/model"
	"time"
)

func GetMockUser() *model.User {
	email := Email()
	return &model.User{
		BaseModel: model.BaseModel{
			ID:        RandID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: Username(),
		Email:    email,
		Password: RandStr(8),
		Image:    fmt.Sprintf("https://gravatar.com/avatar/%s?d=identicon", getMD5Hash(email)),
	}
}
