package gcloudbot

import (
	"context"
	"log"
	"smuel1414/gcloud.vms/vms"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/api/compute/v1"
)

type GcloudbotConfig struct {
	InstanceConfig vms.InstanceConfig
	Bot            *tgbotapi.BotAPI
	Update         *tgbotapi.Update
	ComputeService *compute.Service
	Ctx            *context.Context
}

// Helper to log and send bot message
func SendAndLog(msgText string, bot *tgbotapi.BotAPI, msg *tgbotapi.MessageConfig) {
	// Define error message
	msg.Text = msgText
	// Log error
	log.Println(msg.Text)
	// Answer with error
	bot.Send(msg)
}
