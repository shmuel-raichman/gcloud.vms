package gcloudbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

var NumericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Status all", `{"vm": "all", "action": "status-all"}`),
		tgbotapi.NewInlineKeyboardButtonData("delete", `{"vm": "all", "action": "delete-list"}`),
		tgbotapi.NewInlineKeyboardButtonData("Status List", `{"vm": "all", "action": "status-list"}`),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Delete all", `{"vm": "all", "action": "delete-all"}`),
	),
)


