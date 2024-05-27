package model

import (
	"encoding/json"
	"github.com/google/uuid"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
)

type User struct {
	Id     *uuid.UUID  `json:"id"`
	Email  string      `json:"email"`
	Groups []*AceGroup `json:"groups,omitempty"`
}

func NewUserFromJwt(jwt *dto.TokenClaimsReq, email string) *User {
	return &User{
		Id:     jwt.UserId,
		Email:  email,
		Groups: nil,
	}
}

func NewUserFromId(id *uuid.UUID) *User {
	return &User{
		Id:     id,
		Groups: nil,
	}
}

func NewUserFromDto(userDto *dto.UserUpdateReq) *User {
	groups := make([]*AceGroup, len(userDto.Groups))
	for i, group := range userDto.Groups {
		groups[i] = NewAceGroupFromDto(group)
	}

	return &User{
		Id:     userDto.Id,
		Email:  userDto.Email,
		Groups: groups,
	}
}

func (u *User) ToDto() *dto.UserUpdateReq {
	dtoGroups := make([]*dto.ACEGroupReq, len(u.Groups))
	for i, group := range u.Groups {
		dtoGroups[i] = group.ToDto()
	}
	return &dto.UserUpdateReq{
		Id:     u.Id,
		Email:  u.Email,
		Groups: dtoGroups,
	}
}

func (u *User) SetGroups(groups []*AceGroup) {
	for _, group := range groups {
		u.Groups = append(u.Groups, group)
	}
}

func (u *User) ToJson() ([]byte, error) {
	return json.Marshal(u)
}

func UserFromJson(data []byte) (*User, error) {
	u := &User{}
	if err := json.Unmarshal(data, u); err != nil {
		return nil, err
	}
	return u, nil
}
