package model

import "encoding/json"

type Log struct {
	Id        int64  `json:"id"`
	Operation string `json:"operation"`
	Executor  *User  `json:"executor"`
	CreatedAt string `json:"createdAt"`
	Meta      []byte `json:"meta,omitempty"`
}

func (l *Log) ToJson() *[]byte {
	logJson, err := json.Marshal(l)
	if err != nil {
		return nil
	}
	return &logJson
}
