package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var TOKEN string

type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

func sayHello(chatID string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", TOKEN, chatID, "hello")
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}
	return nil
}

func Handler(res http.ResponseWriter, req *http.Request) {
	body := &webhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}
	if !strings.ContainsAny(strings.ToLower(body.Message.Text), "telegom") {
		return
	}
	err := sayHello(fmt.Sprint(body.Message.Chat.ID))
	if err != nil {
		fmt.Println("error in sending reply:", err)
		return
	}
	fmt.Println("reply sent")
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	TOKEN = os.Getenv("TOKEN")

	fmt.Println(TOKEN)
	http.ListenAndServe(":3000", http.HandlerFunc(Handler))
}
