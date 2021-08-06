package gcloudbot

import (
	"fmt"

	"smuel1414/gcloud.vms/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func StatusList(gcloudbotConfig GcloudbotConfig, action string) {

	chatID := gcloudbotConfig.Update.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	msg.ParseMode = "Markdown"

	var projectID string = gcloudbotConfig.InstanceConfig.ProjectID
	var zone string = gcloudbotConfig.InstanceConfig.Zone

	msgText := fmt.Sprintf("Getting instaces list for project: *%s*\n", projectID)
	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)

	list, err := gcloudbotConfig.ComputeService.Instances.List(projectID, zone).Do()
	if err != nil {
		msgText := fmt.Sprintf("Couldn't list instances for project: *%s*\n%s", projectID, err.Error())
		SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
		return
	}

	if len(list.Items) == 0 {
		msgText := fmt.Sprintf("No VMs in project: *%s*\n", projectID)
		SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
		return
	}

	var row []tgbotapi.InlineKeyboardButton

	// Create inline keyboard rows
	for _, vm := range list.Items {
		str := `{"vm": "` + vm.Name + `", "action": "` + action + `"}`
		vmBotten := tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf("%s: %s", action, vm.Name),
			CallbackData: &str,
		}
		row = append(row, vmBotten)
	}

	// Create inline keyboard
	var vmListKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			row...,
		),
	)

	msg.ReplyMarkup = vmListKeyboard
	msgText = fmt.Sprintf("Json object \n```%s```\nList of instances for project: *%s*\n", gcloudbotConfig.Update.CallbackQuery.Data, projectID)
	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
	// bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
}

func ProjectInstanceStatus(gcloudbotConfig GcloudbotConfig) {

	chatID := gcloudbotConfig.Update.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	msg.ParseMode = "Markdown"

	var projectID string = gcloudbotConfig.InstanceConfig.ProjectID
	var zone string = gcloudbotConfig.InstanceConfig.Zone

	msgText := fmt.Sprintf("Getting instaces list for project: *%s*\n", projectID)
	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)

	list, err := gcloudbotConfig.ComputeService.Instances.List(projectID, zone).Do()
	if err != nil {
		msgText := fmt.Sprintf("Couldn't list instances for project: *%s*\n%s", projectID, err.Error())
		SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
		return
	}

	if len(list.Items) == 0 {
		msgText := fmt.Sprintf("No VMs in project: *%s*\n", projectID)
		SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
		return
	}

	numberOfInstances := len(list.Items)
	responseMsg := fmt.Sprintf("There are - %d - instances in project %s\n ---", numberOfInstances, projectID)
	// Create inline keyboard rows
	for _, instance := range list.Items {
		responseMsg += fmt.Sprintf("\n- Instance: %s state is - %s", instance.Name, instance.Status)
	}

	// msgText = fmt.Sprintf("Json object \n``` %s```\nList of instances for project: *%s*\n", gcloudbotConfig.Update.CallbackQuery.Data, projectID)
	SendAndLog(responseMsg, gcloudbotConfig.Bot, &msg)
	// bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
}

func StatusAndDetails(gcloudbotConfig GcloudbotConfig) {
	chatID := gcloudbotConfig.Update.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, gcloudbotConfig.Update.CallbackQuery.Data)
	msg.ParseMode = "Markdown"

	var projectID string = gcloudbotConfig.InstanceConfig.ProjectID
	var zone string = gcloudbotConfig.InstanceConfig.Zone
	var instanceName string = gcloudbotConfig.InstanceConfig.Name

	msgText := fmt.Sprintf("Getting instance: *%s* - please wait ...%s\n", instanceName, utils.Drizzle)
	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)

	instance, err := gcloudbotConfig.ComputeService.Instances.Get(projectID, zone, instanceName).Do()
	if err != nil {
		msgText := fmt.Sprintf("Falid to get instance: *%s*\n", instanceName)
		SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
		return
	}

	var externalIPs string = ""
	for _, netInterface := range instance.NetworkInterfaces {
		for _, accessConfig := range netInterface.AccessConfigs {
			externalIPs += fmt.Sprintf("\n- External IP: ```%s```", accessConfig.NatIP)
		}
	}

	var diskSizes string = ""
	for _, disk := range instance.Disks {
		diskSizes += fmt.Sprintf("Disk size GB: *%d*\n", disk.DiskSizeGb)
	}

	msgText = fmt.Sprintf(
		"Instance: *%s*\n---------\n- State: *%s* %s\n- %s", instanceName, instance.Status, externalIPs, diskSizes,
	)
	// gcloudbotConfig.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(gcloudbotConfig.Update.CallbackQuery.ID, "msgText"))
	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
}
