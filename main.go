package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/Syfaro/telegram-bot-api"
    "log"
    "net/http"
)

func main() {
    tgbot, err := tgbotapi.NewBotAPI(tg_token)
    if err != nil { log.Panic(err) }

    tgbot.Debug = true

    log.Printf("Telegram connected: %s", tgbot.Self.UserName)

    tgupdate := tgbotapi.NewUpdate(0)
    tgupdate.Timeout = 60

    bc_queue := make(chan []byte)
    go send_to_bc(bc_queue)

    err = tgbot.UpdatesChan(tgupdate)
    if err != nil { log.Panic(err) }

    for update := range tgbot.Updates {
        log.Printf("Message: [%s] %s", update.Message.From.UserName, update.Message.Text)
        if len(update.Message.Text) != 0 {
            var msg_text string
            if len(update.Message.From.LastName) != 0 {
                msg_text = fmt.Sprintf("[%s %s] %s", update.Message.From.FirstName, update.Message.From.LastName, update.Message.Text)
            } else {
                msg_text = fmt.Sprintf("[%s] %s", update.Message.From.FirstName, update.Message.Text)
            }
            bcmsg := map[string]interface{}{
                "text": msg_text,
                "markdown": false,
            }
            bcmsg_json, err := json.Marshal(bcmsg)
            if err != nil {
                log.Println(err)
                continue
            }
            bc_queue <- bcmsg_json
        }
    }
}

func send_to_bc(queue chan []byte) {
    var http_client http.Client
    for msg := range queue {
        msg_buffer := bytes.NewBuffer(msg)
        http_req, err := http.NewRequest("POST", bc_token_in, msg_buffer)
        if err != nil {
            log.Println(err)
            continue
        }
        http_req.Header.Set("Content-Type", "application/json")
        _, err = http_client.Do(http_req)
        if err != nil {
            log.Println(err)
            continue
        }
    }
}
