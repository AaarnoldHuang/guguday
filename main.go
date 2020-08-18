package main

import (
	"database/sql"
	"fmt"
	"guguday/Module"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yanzay/tbot/v2"
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
var tempWelcome string
var tempGroupUsername string

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
		if m.Chat.Type == "supergroup" {
			return
		} else if m.Chat.Type == "private" {
			c.SendMessage(m.Chat.ID, "我是骚鸡，请把我拉进你的群，我会咕咕Day!")
		}
	})

	bot.HandleMessage("/setmine", func(message *tbot.Message) {
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
					msg, _ := c.SendMessage(message.Chat.ID, "他是大猛1惹，假1罚石那种。", tbot.OptReplyToMessageID(message.ReplyToMessage.MessageID))
					time.Sleep(10 * time.Second)
					_ = c.DeleteMessage(message.Chat.ID, msg.MessageID)
				} else if result.Role == "0" {
					msg, _ := c.SendMessage(message.Chat.ID, "他是站街女惹，一晚接八个那种。", tbot.OptReplyToMessageID(message.ReplyToMessage.MessageID))
					time.Sleep(10 * time.Second)
					_ = c.DeleteMessage(message.Chat.ID, msg.MessageID)
				}
			} else {
				msg, _ := c.SendMessage(message.Chat.ID, fmt.Sprintf("我没有他的数据惹，快快来补充吧 [ %s ](tg://user?id= %d )",
					wantedUser.FirstName, wantedUser.ID), tbot.OptReplyToMessageID(message.ReplyToMessage.MessageID), tbot.OptParseModeMarkdown)
				time.Sleep(10 * time.Second)
				_ = c.DeleteMessage(message.Chat.ID, msg.MessageID)
			}
		} else {
			msg, _ := c.SendMessage(message.Chat.ID, "要查哪个就直接回复他的消息，并使用命令\"/checkinfo\", 不要手贱瞎点。",
				tbot.OptReplyToMessageID(message.MessageID))
			time.Sleep(10 * time.Second)
			_ = c.DeleteMessage(message.Chat.ID, msg.MessageID)
		}
	})

	bot.HandleMessage("/setwelcome .+", func(message *tbot.Message) {

		if message.Chat.Type == "private" {
			//check user is admin or not
			cmd := fmt.Sprintf("SELECT * FROM `admin_info` where `admin_uid` = '%d';", message.From.ID)
			result := Module.SelectAdminInfo(DB, cmd)
			if result.Admin_uid == strconv.Itoa(message.From.ID) {
				text := strings.TrimPrefix(message.Text, "/setwelcome ")
				_, _ = c.SendMessage(message.Chat.ID, "你要设置的群组是 @"+text+" ，对吗?")
				_, _ = c.SendMessage(message.Chat.ID, "如果不对，请重新发送，如果对的话请回复 /welcome +欢迎语。例子如下:")
				_, _ = c.SendMessage(message.Chat.ID, "/welcome 欢迎加入群聊。")
				tempGroupUsername = text
			}
		}
	})

	bot.HandleMessage("/setaskrole .+", func(message *tbot.Message) {

		if message.Chat.Type == "private" {
			//check user is admin or not
			cmd := fmt.Sprintf("SELECT * FROM `admin_info` where `admin_uid` = '%d';", message.From.ID)
			result := Module.SelectAdminInfo(DB, cmd)
			if result.Admin_uid == strconv.Itoa(message.From.ID) {
				text := strings.TrimPrefix(message.Text, "/setaskrole ")
				_, _ = c.SendMessage(message.Chat.ID, "你要设置的群组是 @"+text+" ，对吗?")
				_, _ = c.SendMessage(message.Chat.ID, "如果不对，请重新发送，如果对的话请选择:")
				tempGroupUsername = text
				app.changeAskRoleHandler(message)
			}
		}
	})

	bot.HandleMessage("/welcome.*", func(message *tbot.Message) {

		if message.Chat.Type == "private" {
			//check user is admin or not
			cmd := fmt.Sprintf("SELECT * FROM `admin_info` where `admin_uid` = '%d';", message.From.ID)
			result := Module.SelectAdminInfo(DB, cmd)
			if result.Admin_uid == strconv.Itoa(message.From.ID) {
				text := strings.TrimPrefix(message.Text, "/welcome ")
				_, _ = c.SendMessage(message.Chat.ID, "欢迎词将会变为如下信息:")
				_, _ = c.SendMessage(message.Chat.ID, text)
				_, _ = c.SendMessage(message.Chat.ID, "确定请回复 /Done")
				tempWelcome = text
			}
		}
	})

	bot.HandleMessage("/Done", func(message *tbot.Message) {
		if message.Chat.Type == "private" {
			//check user is admin or not
			cmd := fmt.Sprintf("SELECT * FROM `admin_info` where `admin_uid` = '%d';", message.From.ID)
			result := Module.SelectAdminInfo(DB, cmd)
			if result.Admin_uid == strconv.Itoa(message.From.ID) {
				cmd := fmt.Sprintf("INSERT INTO `welcome_message` (`group_username`, `group_welcome`) VALUES ('%s','%s') ON DUPLICATE KEY UPDATE group_welcome='%s';",
					tempGroupUsername, tempWelcome, tempWelcome)

				fmt.Println(cmd)
				//如果执行失败，返回信息
				insertResult, _ := Module.InserttoDB(DB, cmd)
				if !insertResult {
					_, _ = c.SendMessage(message.Chat.ID, "数据保存失败，请重试。")
				}
				_, _ = c.SendMessage(message.Chat.ID, "好了!")

			}
		}
	})

	bot.HandleMessage(".*", func(message *tbot.Message) {
		if message.NewChatMembers != nil {
			newuser := message.NewChatMembers[0]
			if !newuser.IsBot {
				cmd := fmt.Sprintf("SELECT * FROM `whore_info` WHERE `whore_uid` = '%d';",
					newuser.ID)
				result := Module.SelectUserInfo(DB, cmd)

				cmd2 := fmt.Sprintf("SELECT * FROM `welcome_message` WHERE `group_username` = '%s';", message.Chat.Username)
				welcome := Module.SelectWelcome(DB, cmd2)

				//有设置欢迎词
				if welcome.Group_welcome != "" {
					msg, _ := c.SendMessage(message.Chat.ID, welcome.Group_welcome, tbot.OptReplyToMessageID(message.MessageID), tbot.OptParseModeMarkdown)
					if welcome.Ask_role == 1 {
						if result.Uid != 0 {
							if result.Role == "1" {
								msg, _ := c.SendMessage(message.Chat.ID, "他是大猛1惹，假1罚石那种。")
								time.Sleep(10 * time.Second)
								_ = c.DeleteMessage(message.Chat.ID, msg.MessageID)
							} else if result.Role == "0" {
								msg, _ := c.SendMessage(message.Chat.ID, "他是站街女惹，一晚接八个那种。")
								time.Sleep(10 * time.Second)
								_ = c.DeleteMessage(message.Chat.ID, msg.MessageID)
							}
						} else {
							app.votingHandler(message, newuser)
						}
					}
					time.Sleep(30 * time.Second)
					_ = c.DeleteMessage(message.Chat.ID, msg.MessageID)
				} else {
					if result.Uid != 0 {
						if result.Role == "1" {
							msg, _ := c.SendMessage(message.Chat.ID, "他是大猛1惹，假1罚石那种。")
							time.Sleep(10 * time.Second)
							_ = c.DeleteMessage(message.Chat.ID, msg.MessageID)
						} else if result.Role == "0" {
							msg, _ := c.SendMessage(message.Chat.ID, "他是站街女惹，一晚接八个那种。")
							time.Sleep(10 * time.Second)
							_ = c.DeleteMessage(message.Chat.ID, msg.MessageID)
						}
					} else {
						app.votingHandler(message, newuser)
					}
				}
				//没有设置欢迎词，默认开启鸡叫模式
			}
		}
	})

	err := bot.Start()

	if err != nil {
		log.Fatal(err)
	}
}

func (a *application) votingHandler(m *tbot.Message, newUser *tbot.User) {
	buttons := makeButtons()
	checkMap[newUser.ID] = false
	msg, _ := a.client.SendMessage(m.Chat.ID, "你是1还是0？这个世界上没有0.5！",
		tbot.OptInlineKeyboardMarkup(buttons),
		tbot.OptReplyToMessageID(m.MessageID))
	votingID := fmt.Sprintf("%s:%d", m.Chat.ID, msg.MessageID)
	a.votings[votingID] = &voting{}
}

func (a *application) changeAskRoleHandler(m *tbot.Message) {
	buttons := askRoleButtons()
	checkMap[m.From.ID] = false
	msg, _ := a.client.SendMessage(m.Chat.ID, "请选择:",
		tbot.OptInlineKeyboardMarkup(buttons),
		tbot.OptReplyToMessageID(m.MessageID))
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
					_ = a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("数据保存失败，请与我私聊重试。"))
				}
				_, _ = a.client.EditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, "天惹！老公！快给我们康康鸡儿~~~")
				_ = a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("如果你想补充更多信息，请与我私聊。"))
				time.Sleep(10 * time.Second)
				_ = a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
			}
			if cq.Data == "bottom" {
				cmd := fmt.Sprintf("INSERT INTO `whore_info` (`whore_uid`, `whore_age`, `whore_role`, `whore_height`, `whore_bodytype`, `whore_size` ) VALUES ('%d', 'null', '%s', 'null', 'null', 'null') ON DUPLICATE KEY UPDATE whore_role=0;",
					cq.From.ID, "0")
				//如果执行失败，返回信息
				insertResult, _ := Module.InserttoDB(DB, cmd)
				if !insertResult {
					_ = a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("数据保存失败，请与我私聊重试。"))
				}
				_, _ = a.client.EditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, "姐妹！快给我们康康菊~~~")
				_ = a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("如果你想补充更多信息，请与我私聊。"))
				time.Sleep(10 * time.Second)
				_ = a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
			}
			if cq.Data == "moreinfo" {
				_ = a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("如果你想补充更多信息，请与我私聊。"))
				_ = a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
			}
			if cq.Data == "1" {
				cmd := fmt.Sprintf("UPDATE `welcome_message` SET `ask_role`='1' WHERE `group_username`='%s';",
					tempGroupUsername)
				//如果执行失败，返回信息
				insertResult, _ := Module.InserttoDB(DB, cmd)
				if !insertResult {
					_ = a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("数据保存失败，请重试。"))
				}
				_, _ = a.client.EditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, "好了!")

			}
			if cq.Data == "0" {
				cmd := fmt.Sprintf("UPDATE `welcome_message` SET `ask_role`='0' WHERE `group_username`='%s';",
					tempGroupUsername)
				//如果执行失败，返回信息
				insertResult, _ := Module.InserttoDB(DB, cmd)
				if !insertResult {
					_ = a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("数据保存失败，请重试。"))
				}
				_, _ = a.client.EditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, "好了!")
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

func askRoleButtons() *tbot.InlineKeyboardMarkup {
	button1 := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("问"),
		CallbackData: "1",
	}
	button2 := tbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("不问"),
		CallbackData: "0",
	}

	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{button1, button2},
		},
	}
}
