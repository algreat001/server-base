package dto

type SendReq struct {
	Message string `json:"message"`
}

type TelegramResponse struct {
	Ok          bool   `json:"ok"`
	Error_code  int    `json:"error_code"`
	Description string `json:"description"`
}
