package types

import (
	"strconv"
	"strings"
)

type (
	User struct {
		Id             int64                  `json:"id"`
		Username       string                 `json:"username"`
		FirstName      string                 `json:"first_name"`
		LastName       string                 `json:"last_name"`
		Fullname       string                 `json:"fullname"`
		Avatar         string                 `json:"avatar"`
		Role           []string               `json:"role"`
		State          string                 `json:"state"`
		AuthProvider   string                 `json:"auth_provider"`
		AuthProviderId string                 `json:"auth_provider_id"`
		AuthToken      string                 `json:"auth_token,omitempty"`
		Options        map[string]interface{} `json:"options"`
		Deleted        bool                   `json:"deleted"`
		Phone          string                 `json:"phone"`
		Password       string                 `json:"password,omitempty"`
		Email          string                 `json:"email,omitempty"`
	}
)

func (u *User) IdString() string {
	return strconv.FormatInt(u.Id, 10)
}

func (u *User) GetRoleAsString() string {
	return strings.Join(u.Role, "_")
}