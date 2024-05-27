package dto

import (
	"github.com/google/uuid"
)

type LogExecutor struct {
	Id    *uuid.UUID `json:"id"`
	Email string     `json:"email,omitempty"`
}

type PageReq struct {
	Page       int    `json:"numPage"`
	PageSize   int    `json:"sizePage"`
	Order      string `json:"order"`
	Descending bool   `json:"descending"`
	Filter     string `json:"filter"`
}

type LogReq struct {
	Id        int64    `json:"id"`
	Operation string   `json:"operation"`
	Executor  *UserReq `json:"executor"`
	CreatedAt string   `json:"createdAt"`
	Meta      []byte   `json:"meta,omitempty"`
}

type LogPage struct {
	Data           []*LogReq `json:"data"`
	CountAllRecord int       `json:"countAllRecord"`
}

func LogExecutorFromJson(data []byte) (*LogExecutor, error) {
	return FromJson[LogExecutor](data)
}
func PageReqFromJson(data []byte) (*PageReq, error) {
	return FromJson[PageReq](data)
}
func LogReqFromJson(data []byte) (*LogReq, error) {
	return FromJson[LogReq](data)
}
func LogPageFromJson(data []byte) (*LogPage, error) {
	return FromJson[LogPage](data)
}
