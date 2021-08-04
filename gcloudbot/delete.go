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
		msgText := fmt.Sprintf("Deleting instance - *%s* faild\n%s", instanceName, err.Error())
		SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
		return
	}

	msgText = fmt.Sprintf("Instance - *%s* deleted succesfuly\n", instanceName)
	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
}

func DeleteAll() {
	// msg = tgbotapi.NewMessage(chatID, update.CallbackQuery.Data)

	// vmList, err := vms.GetVMs(computeService, instanceConfig.ProjectID, instanceConfig.Zone)
	// if err != nil {
	// 	log.Println(err)
	// 	msg.Text = err.Error()
	// 	bot.Send(msg)
	// }

	// for _, vm := range vmList {
	// 	instanceConfig.Name = vm.Name
	// 	if instanceConfig.Name == os.Getenv("BOT_VM") {
	// 		continue
	// 	}
	// 	vms.DeleteInstance(computeService, ctx, &instanceConfig)

	// 	_, err := vms.GetVMStatus(computeService, instanceConfig.ProjectID, instanceConfig.Zone, instanceConfig.Name)
	// 	if err != nil {
	// 		msg.Text = err.Error() + fmt.Sprintf("VM %s deleted succesfuly\n", instanceConfig.Name)
	// 		bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, msg.Text))
	// 		bot.Send(msg)
	// 		log.Println(err)
	// 		log.Println(msg.Text)
	// 	}
	// 	msg.Text = "Finished delete all vms"
	// 	bot.Send(msg)
	// }
}
