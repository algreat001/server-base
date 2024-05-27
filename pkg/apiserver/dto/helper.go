package dto

import "encoding/json"

func FromJson[TRequestType any](data []byte) (*TRequestType, error) {
	var result TRequestType
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
