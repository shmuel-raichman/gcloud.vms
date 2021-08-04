package gcloudbot

import (
	"fmt"
	"os"

	"smuel1414/gcloud.vms/vms"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Delete(gcloudbotConfig GcloudbotConfig) {

	chatID := gcloudbotConfig.Update.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	msg.ParseMode = "Markdown"

	var instanceName string = gcloudbotConfig.InstanceConfig.Name

	if instanceName == os.Getenv("BOT_VM") {
		msgText := fmt.Sprintf("Instance - *%s* should not be deleted ..\n", instanceName)
		SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
		return
	}

	msgText := fmt.Sprintf("Deleting instance - *%s* please wait ..\n", instanceName)
	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)

	err := vms.DeleteInstance(gcloudbotConfig.ComputeService, *gcloudbotConfig.Ctx, &gcloudbotConfig.InstanceConfig)
	if err != nil {
		msgText := fmt.Sprintf("Deleting instance - *%s* faild\n%s\n\n", instanceName, err.Error())
		SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
		return
	}

	msgText = fmt.Sprintf("Instance - *%s* deleted succesfuly\n", instanceName)
	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
}

func DeleteAll(gcloudbotConfig GcloudbotConfig) {

	chatID := gcloudbotConfig.Update.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	msg.ParseMode = "Markdown"

	var projectID string = gcloudbotConfig.InstanceConfig.ProjectID
	var zone string = gcloudbotConfig.InstanceConfig.Zone

	msgText := fmt.Sprintf("Getting instace list for project: *%s*\n", projectID)
	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)

	instances, err := gcloudbotConfig.ComputeService.Instances.List(projectID, zone).Do()
	if err != nil {
		msgText := fmt.Sprintf("Couldn't list instances for project: *%s*\n%s", projectID, err.Error())
		SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
		return
	}

	if len(instances.Items) == 0 {
		msgText := fmt.Sprintf("No instances to delete for project: *%s*\n", projectID)
		SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
		return
	}

	msgText = fmt.Sprintf("Deleting all project - *%s* -  instances please wait ..\n", projectID)
	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)

	for _, instace := range instances.Items {
		gcloudbotConfig.InstanceConfig.Name = instace.Name
		Delete(gcloudbotConfig)
	}

	msgText = fmt.Sprintf("Finished deletion for project - *%s*\n", projectID)
	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
}
