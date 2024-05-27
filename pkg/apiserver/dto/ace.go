package dto

import (
	"github.com/google/uuid"
)

type ACEGroupReq struct {
	Id   *uuid.UUID `json:"id,omitempty"`
	Name string     `json:"name,omitempty"`
}

type ACEGroupIdReq struct {
	Id *uuid.UUID `json:"id"`
}

type ACEIdReq struct {
	Id *uuid.UUID `json:"id"`
}

type ACEToGroupReq struct {
	Id    *uuid.UUID  `json:"id,omitempty"`
	Group ACEGroupReq `json:"group"`
	Path  string      `json:"path"`
}

func ACEGroupReqFromJson(data []byte) (*ACEGroupReq, error) {
	return FromJson[ACEGroupReq](data)
}
func ACEGroupIdReqFromJson(data []byte) (*ACEGroupIdReq, error) {
	return FromJson[ACEGroupIdReq](data)
}
func ACEIdReqFromJson(data []byte) (*ACEIdReq, error) {
	return FromJson[ACEIdReq](data)
}
func ACEToGroupReqFromJson(data []byte) (*ACEToGroupReq, error) {
	return FromJson[ACEToGroupReq](data)
}
