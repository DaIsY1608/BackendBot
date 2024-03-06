package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type SignUpStruct struct {
	Name          string
	TelegramLogin string
	Password      string
}

var SignUpSlice = []SignUpStruct{} // ? empty
func main() {
	r := gin.Default()

	r.Use(Cors)
	r.POST("/signup", SignUp)
	go Recovery()
	r.Run(":3434")
}

func Recovery() {
	ReadUser()
	botresult, err := tgbotapi.NewBotAPI("6711222013:AAEoRHYaHuu86J4lgA22taX3Yr4dcVcL4Y0")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println("Connected!")
	updates := tgbotapi.NewUpdate(0)
	allUpdates, updateErr := botresult.GetUpdatesChan(updates)
	if updateErr != nil {
		fmt.Printf("updateErr: %v\n", updateErr)
	}

	for each := range allUpdates {
		if each.Message.IsCommand() {
			if each.Message.Command() == "reset" {
				for _, v := range SignUpSlice {
					if v.TelegramLogin == each.Message.Chat.UserName {
						msg := tgbotapi.NewMessage(each.Message.Chat.ID, "Enter New Password")
						botresult.Send(msg)
					}
				}
			}
		} else {

			yast := false
			for i, user := range SignUpSlice {
				if user.TelegramLogin == each.Message.Chat.UserName {
					yast = true
					SignUpSlice[i].Password = each.Message.Text
					msg := tgbotapi.NewMessage(each.Message.Chat.ID, "Password changed")
					botresult.Send(msg)
					WriteUser()

				}
			}
			if !yast {
				msg := tgbotapi.NewMessage(each.Message.Chat.ID, "Username not found")
				botresult.Send(msg)
			}
		}
	}
}

func SignUp(c *gin.Context) {
	var SignUpTemp SignUpStruct
	c.ShouldBindJSON(&SignUpTemp)

	if SignUpTemp.Name == "" || SignUpTemp.Password == "" || SignUpTemp.TelegramLogin == "" {
		c.JSON(404, "Empty field")
	} else {
		ReadUser()
		SignUpSlice = append(SignUpSlice, SignUpTemp)
		WriteUser()
	}
}

func WriteUser() {
	marsheledData, _ := json.Marshal(SignUpSlice)
	ioutil.WriteFile("app.json", marsheledData, 0644)
}

func ReadUser() {
	readByte, _ := ioutil.ReadFile("app.json")
	json.Unmarshal(readByte, &SignUpSlice)
}

func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Al low-Origin", "http://192.168.43.246:5500")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	}

	c.Next()
}
