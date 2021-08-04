package gcloudbot

import (
	"context"
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
