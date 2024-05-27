package model

import "github.com/gin-gonic/gin"

type UserInfo struct {
	Ip        string `json:"ip"`
	UserAgent string `json:"ua"`
}

func NewUserInfo(c *gin.Context) *UserInfo {
	return &UserInfo{
		Ip:        c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}
}
