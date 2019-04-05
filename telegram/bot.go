package telegram

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Bot struct {
	builtUrl string
}

func (tb *Bot) SendMessage(message string) {

	form := url.Values{}
	form.Add("chat_id", os.Getenv("TELEGRAM_CHANNEL"))
	form.Add("parse_mode", "markdown")
	form.Add("text", message)

	req, err := http.NewRequest("POST", tb.builtUrl, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("TELEGRAM: An error ocurred sending message: %s\n", err)
	}

	log.Printf("TELEGRAM: Result -> %s\n", body)
}

func New() *Bot {
	return &Bot{
		builtUrl: fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("TELEGRAM_BOT_KEY")),
	}
}
