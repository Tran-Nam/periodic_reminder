package telegram

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// SendTelegramMessage sends a text message using the Bot API.
// token is the bot token (e.g. 123456:ABC-DEF...), chatID is the destination chat id.
func SendTelegramMessage(token string, chatID string, text string) error {
	if token == "" {
		return errors.New("telegram token empty")
	}
	api := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", url.PathEscape(token))
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", text)
	fmt.Println(chatID)
	fmt.Println(token)

	resp, err := http.PostForm(api, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: status=%d body=%s", resp.StatusCode, string(body))
	}
	return nil
}

// SendMessageFromEnv reads TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID from the environment
// and sends the provided text.
func SendMessageFromEnv(text string) error {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	if token == "" {
		return errors.New("TELEGRAM_BOT_TOKEN not set")
	}
	if chatID == "" {
		return errors.New("TELEGRAM_CHAT_ID not set")
	}
	fmt.Println(token)
	fmt.Println(chatID)
	return SendTelegramMessage(token, chatID, text)
}
