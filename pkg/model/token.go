package model

import "encoding/json"

type TokenPage struct {
	Data           []*Token `json:"data"`
	CountAllRecord int      `json:"countAllRecord"`
}

type Token struct {
	Token     string `json:"token"`
	CreatedAt string `json:"created_at,omitempty"`
}

func NewToken(token string) *Token {
	return &Token{
		Token: token,
	}
}

func (t *Token) ToJson() *[]byte {
	tokenJson, err := json.Marshal(t)
	if err != nil {
		return nil
	}
	return &tokenJson
}
