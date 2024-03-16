package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func main() {
	var err error
	Bot, err = tgbotapi.NewBotAPI(Config.Token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	go push()

	bot()
}

func bot() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := Bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.FromChat().ID != Config.Admin {
			continue
		}

		if update.Message != nil && update.Message.IsCommand() {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "add":
				args := update.Message.CommandArguments()
				if args == "" {
					msg.Text = "格式错误：/add 番剧名称"
				} else {
					link, err := search(args)
					if err != nil {
						fmt.Println("search error:", err)
						msg.Text = "搜索时出现错误"
					} else {
						links := getRss(link, true)
						for _, link := range links {
							LinkMap[link.Name] = ItemInfoLink{
								Url:        link.Url,
								Downloaded: false,
								Pushed:     false,
							}
						}
						saveLink()
						msg.Text = "已添加"
					}
				}
			case "del":
				args := update.Message.CommandArguments()
				if args == "" {
					msg.Text = "格式错误：/del 番剧名称"
				} else {
					link, err := search(args)
					if err != nil {
						fmt.Println("search error:", err)
						msg.Text = "搜索时出现错误"
					} else {
						isDel := false
						for i, j := range Data.Tasks {
							if j == link {
								Data.Tasks = append(Data.Tasks[:i], Data.Tasks[i+1:]...)
								isDel = true
								break
							}
						}
						if isDel {
							saveData()
							msg.Text = "已删除"
						} else {
							msg.Text = "没有找到目标"
						}
					}
				}
			case "get":
				var links = make([]string, 0, len(LinkMap))
				for _, v := range LinkMap {
					if !v.Downloaded {
						links = append(links, v.Url)
					}
				}

				if len(links) == 0 {
					msg.Text = "不存在未下载的项目"
				} else {
					msg.ParseMode = "markdown"
					msg.Text = "```\n" + strings.Join(links, "\n") + "\n```\n警告：按钮会把当前所有的项目全部设为已下载，请确保这里的数据是最新的"
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("全部标记为已下载", "clear"),
						),
					)
				}
			case "help":
				msg.Text = "/add <番剧名称> 添加番剧\n/del <番剧名称> 移除番剧\n/get 获取当前所有标记为未下载的种子"
			default:
				msg.Text = "请使用 /help 获取可以使用的命令"
			}
			_, err := Bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		} else if update.CallbackQuery != nil {
			callback := update.CallbackQuery

			switch callback.Data {
			case "set":
				k := strings.Split(callback.Message.Text, "\n\n")[0]
				fmt.Println(k)
				if v, ok := LinkMap[string(k)]; ok {
					v.Downloaded = true
					LinkMap[string(k)] = v

					saveLink()

					msg := strings.ReplaceAll(callback.Message.Text, "[", "\\[")
					editMsg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "#Finish\n"+msg)
					editMsg.ParseMode = "markdown"
					Bot.Send(editMsg)
				}
			case "clear":
				for k, v := range LinkMap {
					if !v.Downloaded {
						v.Downloaded = true
						LinkMap[k] = v
					}
				}
				saveLink()

				editMsg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "#Finish\n以上的所有推送已标记为下载完成")
				Bot.Send(editMsg)
			}

			newCallback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := Bot.Request(newCallback); err != nil {
				log.Println("Callback Err:", err)
			}
		}
	}
}

func push() {
	for _, v := range Data.Tasks {
		links := getRss(v, false)
		for _, link := range links {
			LinkMap[link.Name] = ItemInfoLink{
				Url:        link.Url,
				Downloaded: false,
				Pushed:     false,
			}
		}
	}

	for k, v := range LinkMap {
		if !v.Pushed {
			msg := tgbotapi.NewMessage(Config.Admin, strings.ReplaceAll(k, "[", "\\[")+"\n\n`"+v.Url+"`")
			msg.ParseMode = "markdown"
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("标记为已下载", "set"),
				),
			)
			_, err := Bot.Send(msg)
			if err != nil {
				log.Println("Push Error:", err)
			}

			v.Pushed = true
			LinkMap[k] = v
		}
	}
	saveLink()

	time.AfterFunc(30*time.Minute, push)
}
