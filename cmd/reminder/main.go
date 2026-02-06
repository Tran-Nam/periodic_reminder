package main

import (
	"log"
	"encoding/json"
	"fmt"
	"os"
	"net/http"
	"strings"
	"bytes"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GistUpdate struct {
	Files map[string]GistFile `json:"files"`
}

type GistFile struct {
	Content string `json:"content"`
}

func getGistContent(gistID, token string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/gists/%s", gistID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	files := result["files"].(map[string]interface{})
	dbFile := files["database.csv"].(map[string]interface{})
	return dbFile["content"].(string), nil
}

func updateGist(gistID, token, content string) error {
	url := fmt.Sprintf("https://api.github.com/gists/%s", gistID)
	update := GistUpdate{
		Files: map[string]GistFile{
			"database.csv": {Content: content},
		},
	}
	jsonData, _ := json.Marshal(update)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer " + token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	gistID := os.Getenv("GIST_ID")
	ghToken := os.Getenv("GH_TOKEN")

	fmt.Println(botToken)
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Connected to Bot: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("Message: %s", update.Message.Text)
		msg := update.Message.Text 
		chatID := update.Message.Chat.ID 

		if strings.HasPrefix(msg, "/add") {
			input := strings.TrimPrefix(msg, "/add ")
			oldData, _ := getGistContent(gistID, ghToken)

			newData := oldData + fmt.Sprintf("\n%s,%s", os.Getenv("CURRENT_DATE"), input)
			updateGist(gistID, ghToken, newData)

			msg := tgbotapi.NewMessage(chatID, "added")
			bot.Send(msg)
		} else if msg == "/list" {
			content, _ := getGistContent(gistID, ghToken)
			msg := tgbotapi.NewMessage(chatID, content)
			bot.Send(msg)
		}
	}
}
