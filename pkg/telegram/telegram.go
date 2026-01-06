package telegram

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// LoadDotEnv loads simple KEY=VALUE lines from a .env file into the process environment.
// It is intentionally lightweight to avoid adding third-party deps.
func LoadDotEnv(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		os.Setenv(key, val)
	}
	return nil
}

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
	return SendTelegramMessage(token, chatID, text)
}
