package notify

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var urlTemplate = "https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s"
var apiTimeout = 5 * time.Second
var tokenParam = "token"
var channelIdParam = "channel-id"

type TelegramSender struct {
	token     string
	channelId string
}

func NewTelegramSender(params map[string]string) (*TelegramSender, error) {
	token, ok := params[tokenParam]
	if !ok || len(token) == 0 {
		return nil, fmt.Errorf("'%s' is a required parameter for telegram sender", tokenParam)
	}
	channelId, ok := params[channelIdParam]
	if !ok || len(channelId) == 0 {
		return nil, fmt.Errorf("'%s' is a required parameter for telegram sender", channelIdParam)
	}

	return &TelegramSender{token, channelId}, nil
}

func (ts *TelegramSender) Send(text string) {
	encodedText := url.QueryEscape(text)
	apiUrl := fmt.Sprintf(urlTemplate, ts.token, ts.channelId, encodedText)

	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		fmt.Printf("Couldn't construct request for: %s, err: %v", apiUrl, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil || (resp != nil && resp.StatusCode != 200) {
		status := 500
		if resp != nil {
			status = resp.StatusCode
		}
		fmt.Printf("Error calling %s, err: %v, status code: %d\n", apiUrl, err, status)
	}
}
