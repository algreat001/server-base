package model

import (
	"github.com/google/uuid"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
)

type AceGroup struct {
	Id   *uuid.UUID `json:"id"`
	Name string     `json:"name"`
}

func GetAceGroup(id *uuid.UUID, name string) *AceGroup {
	return &AceGroup{
		Id:   id,
		Name: name,
	}
}

func NewAceGroupFromDto(req *dto.ACEGroupReq) *AceGroup {
	return GetAceGroup(req.Id, req.Name)
}

func (u *AceGroup) ToDto() *dto.ACEGroupReq {
	return &dto.ACEGroupReq{
		Id:   u.Id,
		Name: u.Name,
	}
}
