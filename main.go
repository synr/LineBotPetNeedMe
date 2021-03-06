// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

//PetDB :
var PetDB *Pets

func main() {
	var err error
	PetDB = NewPets()
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				var pet *Pet
				log.Println(message.Text)
				inText := strings.ToLower(message.Text)
				if strings.Contains(inText, "狗") || strings.Contains(inText, "dog") {
					pet = PetDB.GetNextDog()
				} else if strings.Contains(inText, "貓") || strings.Contains(inText, "cat") {
					pet = PetDB.GetNextCat()
				}

				if pet == nil {
					pet = PetDB.GetNextPet()
				}

				out := fmt.Sprintf("您好，目前的動物：名為%s, 所在地為:%s, 敘述: %s 電話為:%s 圖片網址在: %s", pet.Name, pet.Resettlement, pet.Note, pet.Phone, pet.ImageName)
				//obj_message := linebot.NewImageMessage(pet.ImageName, pet.ImageName)
				//因為圖片不是 https 協定，所以 LINE API 不能直接秀圖

				//if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(out),obj_message).Do(); err != nil {
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(out)).Do(); err != nil {
					log.Print(err)
				}

				log.Println("Img:", pet.ImageName)

				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(pet.ImageName, pet.ImageName)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
