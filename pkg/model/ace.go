package model

import "github.com/google/uuid"

type Ace struct {
	Id    *uuid.UUID `json:"id"`
	Path  string     `json:"path"`
	Group *AceGroup  `json:"group"`
}

func GetACE(id *uuid.UUID, path string, group *AceGroup) *Ace {
	return &Ace{
		Id:    id,
		Path:  path,
		Group: group,
	}
}
