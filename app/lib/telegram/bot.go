package telegram

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const TELEGRAM_URL = "https://api.telegram.org"

type teleBot struct {
	token   string
	baseUrl string
}

func NewTelegramBot(token string) *teleBot {
	return &teleBot{
		token:   token,
		baseUrl: fmt.Sprintf("%s/bot%s", TELEGRAM_URL, token),
	}
}

func (t *teleBot) httpCall(path string) (*http.Response, error) {

	var client = &http.Client{}
	url := fmt.Sprintf("%s/%s", t.baseUrl, path)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (t *teleBot) SendMessage(chatId, text, parseMode string) error {
	var param = url.Values{}

	param.Set("chat_id", chatId)
	param.Set("text", text)
	param.Set("parse_mode", parseMode)

	fullPath := fmt.Sprintf("%s?%s", "sendMessage", param.Encode())
	response, err := t.httpCall(fullPath)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check the response status
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Telegram API request failed: %s", response.Status)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	fmt.Println("Success send message to telegram")

	return nil
}
