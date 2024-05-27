package dto

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type ServiceClaimsReq struct {
	Id           string `json:"id"`
	Database     int    `json:"database"`
	FileResource int    `json:"files"`
}

type TokenClaimsReq struct {
	Subject  string             `json:"sub"`
	UserId   *uuid.UUID         `json:"user_id"`
	Expired  int64              `json:"exp"`
	Services []ServiceClaimsReq `json:"services"`
	jwt.RegisteredClaims
}
