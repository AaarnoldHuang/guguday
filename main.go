package main

import (
	"database/sql"
	"fmt"
	"github.com/yanzay/tbot/v2"
	"guguday/Module"
	"log"
	"os"
	"time"
)

type application struct {
	client  *tbot.Client
	votings map[string]*voting
}
type voting struct {
	top    int
	bottom int
}

var tempid int
var DB *sql.DB
var checkMap map[int]bool

func main() {
	DB = Module.ConnectDB()
	checkMap = make(map[int]bool)
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
		app.votingHandler(message, message.From)
	})

	bot.HandleMessage("/checkinfo", func(message *tbot.Message) {
		if message.ReplyToMessage != nil {
			wantedUser := message.ReplyToMessage.From
			cmd := fmt.Sprintf("SELECT * FROM `whore_info` WHERE `whore_uid` = '%d';",
				wantedUser.ID)
			result := Module.SelectUserInfo(DB, cmd)
			if result.Uid != 0 {
				if result.Role == "1" {
					c.SendMessage(message.Chat.ID, "他是大猛1惹，假1罚石那种。")
				} else if result.Role == "0" {
					c.SendMessage(message.Chat.ID, "他是站街女惹，一晚接八个那种。")
				}
			} else {
				c.SendMessage(message.Chat.ID, "我没有他的数据惹，快快来补充吧 @"+wantedUser.Username)
			}
		}
	})

	bot.HandleMessage(".*", func(message *tbot.Message) {
		if message.NewChatMembers != nil {
			newuser := message.NewChatMembers[0]
			c.SendMessage(message.Chat.ID, "欢迎新骚鸡进群。\n来，大家热烈欢迎 @"+newuser.Username)
			app.votingHandler(message, newuser)
		}
	})
	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func (a *application) votingHandler(m *tbot.Message, newuser *tbot.User) {
	buttons := makeButtons()
	checkMap[newuser.ID] = false
	msg, _ := a.client.SendMessage(m.Chat.ID, "你是1还是0？这个世界上没有0.5！", tbot.OptInlineKeyboardMarkup(buttons))
	votingID := fmt.Sprintf("%s:%d", m.Chat.ID, msg.MessageID)
	a.votings[votingID] = &voting{}
}

func (a *application) callbackHandler(cq *tbot.CallbackQuery) {

	for key, _ := range checkMap {
		if key == cq.From.ID {
			delete(checkMap, cq.From.ID)
			if cq.Data == "top" {
				cmd := fmt.Sprintf("INSERT INTO `whore_info` (`whore_uid`, `whore_age`, `whore_role`, `whore_height`, `whore_bodytype`, `whore_size` ) VALUES ('%d', 'null', '%s', 'null', 'null', 'null') ON DUPLICATE KEY UPDATE whore_role=1;",
					cq.From.ID, "1")
				//如果执行失败，返回信息
				insertResult, _ := Module.InserttoDB(DB, cmd)
				if !insertResult {
					a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("数据保存失败，请与我私聊重试。"))
				}
				a.client.EditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, "天惹！老公！快给我们康康鸡儿~~~")
				a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("如果你想补充更多信息，请与我私聊。"))
				time.Sleep(20 * time.Second)
				a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
			}
			if cq.Data == "bottom" {
				cmd := fmt.Sprintf("INSERT INTO `whore_info` (`whore_uid`, `whore_age`, `whore_role`, `whore_height`, `whore_bodytype`, `whore_size` ) VALUES ('%d', 'null', '%s', 'null', 'null', 'null') ON DUPLICATE KEY UPDATE whore_role=0;",
					cq.From.ID, "0")
				//如果执行失败，返回信息
				insertResult, _ := Module.InserttoDB(DB, cmd)
				if !insertResult {
					a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("数据保存失败，请与我私聊重试。"))
				}
				a.client.EditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, "姐妹！快给我们康康菊~~~")
				a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("如果你想补充更多信息，请与我私聊。"))
				time.Sleep(20 * time.Second)
				a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
			}
			if cq.Data == "moreinfo" {
				a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("如果你想补充更多信息，请与我私聊。"))
				a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
			}
		}
	}
}

func makeButtons() *tbot.InlineKeyboardMarkup {
	button1 := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("我是1"),
		CallbackData: "top",
	}
	button2 := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("我是0"),
		CallbackData: "bottom",
	}

	button3 := tbot.InlineKeyboardButton{
		Text:                         fmt.Sprintf("我想完善我的信息"),
		URL:                          "",
		LoginURL:                     nil,
		CallbackData:                 "moreinfo",
		SwitchInlineQuery:            nil,
		SwitchInlineQueryCurrentChat: nil,
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{button1, button2, button3},
		},
	}
}
