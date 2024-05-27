package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/config"
)

func Send2Telegram(messageText string) error {

	cfg := config.GetInstance()
	request := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?disable_web_page_preview=true&parse_mode=Markdown&chat_id=%s&text=%s", cfg.Telegram.BotId, cfg.Telegram.ChatId, messageText)

	resp, err := http.Get(request)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respJSON, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	respStruct := &dto.TelegramResponse{}

	err = json.Unmarshal(respJSON, &respStruct)

	if err != nil {
		return err
	}
	if !respStruct.Ok {
		return errors.New(fmt.Sprintf("Telegram response: Error code: %d, description: %s", respStruct.Error_code, respStruct.Description))
	}

	return nil

}
