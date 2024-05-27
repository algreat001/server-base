package model

type Lang struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Meta []byte `json:"meta,omitempty"`
}

func GetLang(code string, name string, meta []byte) *Lang {
	return &Lang{
		Code: code,
		Name: name,
		Meta: meta,
	}
}
