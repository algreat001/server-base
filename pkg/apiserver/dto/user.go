package dto

import (
	"github.com/google/uuid"
)

type UserReq struct {
	Id    *uuid.UUID `json:"id"`
	Email *string    `json:"email,omitempty"`
}

type UserInAccessGroupReq struct {
	UserId  *uuid.UUID `json:"userId"`
	GroupId *uuid.UUID `json:"groupId"`
}

type UserUpdateReq struct {
	Id     *uuid.UUID     `json:"id,omitempty"`
	Email  string         `json:"email"`
	Groups []*ACEGroupReq `json:"groups,omitempty"`
}

type UsersPageReq struct {
	Page       int    `json:"numPage"`
	PageSize   int    `json:"sizePage"`
	Order      string `json:"order"`
	Descending *bool  `json:"descending"`
	Filter     string `json:"filter,omitempty"`
}

type UserEmailReq struct {
	Email string `json:"email"`
}

type UserPage struct {
	Data           []*UserUpdateReq `json:"data"`
	CountAllRecord int              `json:"countAllRecord"`
}

func UserReqFromJson(data []byte) (*UserReq, error) {
	return FromJson[UserReq](data)
}
func UserInAccessGroupReqFromJson(data []byte) (*UserInAccessGroupReq, error) {
	return FromJson[UserInAccessGroupReq](data)
}
func UserUpdateReqFromJson(data []byte) (*UserUpdateReq, error) {
	return FromJson[UserUpdateReq](data)
}
func UsersPageReqFromJson(data []byte) (*UsersPageReq, error) {
	return FromJson[UsersPageReq](data)
}
func UserEmailReqFromJson(data []byte) (*UserEmailReq, error) {
	return FromJson[UserEmailReq](data)
}
func UserPageFromJson(data []byte) (*UserPage, error) {
	return FromJson[UserPage](data)
}

type ClientToUserReq struct {
	ClientId *uuid.UUID `json:"client_id"`
	UserId   *uuid.UUID `json:"user_id"`
}
