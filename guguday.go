package main

import (
	"fmt"
	"github.com/yanzay/tbot/v2"
	"log"
	"os"
	)

type application struct {
	client  *tbot.Client
	votings map[string]*voting
}
type voting struct {
	top   int
	bottom int
}

var tempid int

func main() {
	bot := tbot.New(os.Getenv("TELEGRAM_TOKEN"))
	c := bot.Client()
	app := &application{
		votings: make(map[string]*voting),
	}
	app.client = bot.Client()
	bot.HandleCallback(app.callbackHandler)
	bot.HandleMessage("/start", func(m *tbot.Message) {
		c.SendMessage(m.Chat.ID, "我是骚鸡，请把我拉进你的群，我会咕咕Day!")
	})

	bot.HandleMessage("/testtesttest", func(message *tbot.Message) {
		app.votingHandler(message)
	})

	bot.HandleMessage(".*", func(message *tbot.Message) {
		if (message.NewChatMembers != nil) {
			fmt.Println("new join in")
			newuser := message.NewChatMembers[0]
			//fmt.Println(newuser.FirstName)
			//fmt.Println(newuser.ID)
			c.SendMessage(message.Chat.ID, "欢迎新骚鸡进群。\n来，大家热烈欢迎" + newuser.FirstName)
			app.votingHandler(message)
		}
	})
	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func (a *application) votingHandler(m *tbot.Message) {
	buttons := makeButtons()
	msg, _ := a.client.SendMessage(m.Chat.ID, "你是1还是0？这个世界上没有0.5！", tbot.OptInlineKeyboardMarkup(buttons))
	votingID := fmt.Sprintf("%s:%d", m.Chat.ID, msg.MessageID)
	//m.From.ID
	tempid = m.From.ID
	a.votings[votingID] = &voting{}
}

func (a *application) callbackHandler(cq *tbot.CallbackQuery) {
	if (tempid != cq.From.ID) {
		a.client.SendMessage(cq.Message.Chat.ID, "瞎jb点啊！")
	} else {
		if cq.Data == "top" {
			a.client.EditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, "天惹！老公！快给我们康康鸡儿~~~")
		}
		if cq.Data == "bottom" {
			a.client.EditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, "姐妹！快给我们康康菊~~~")
		}
		a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("OK"))
	}
}


func makeButtons() *tbot.InlineKeyboardMarkup{
	button1 := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("我是1"),
		CallbackData: "top",
	}
	button2 := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("我是0"),
		CallbackData: "bottom",
	}
	fmt.Println(button1, button2)
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{button1, button2},
		},
	}
}