package gcloudbot

import (
	"fmt"
	"log"
	"smuel1414/gcloud.vms/vms"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// func BotCreateInstance(instanceConfig vms.InstanceConfig, bot tgbotapi.BotAPI, update tgbotapi.Update, computeService *compute.Service, ctx context.Context) {
func BotCreateInstance(gcloudbotConfig GcloudbotConfig) {

	chatID := gcloudbotConfig.Update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	msg.ParseMode = "Markdown"

	// First answer
	msgText := fmt.Sprintf("You supplied the following instance name: *%s*", gcloudbotConfig.Update.Message.CommandArguments())
	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
	// Wait message
	msgText = "Now creating your instance please wait ..."
	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)

	// TODO validate vm name
	instanceName := gcloudbotConfig.Update.Message.CommandArguments()
	// matched, err := regexp.MatchString(`(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?)`, gcloudbotConfig.Update.Message.CommandArguments())
	// if err != nil {
	// 	msgText := fmt.Sprintf("Not creating instance, not matched regex\n%s", err)
	// 	SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
	// 	return
	// }
	if !(len(instanceName) > 0) {
		msgText := "Not creating instance, Missing instance name"
		SendAndLog(msgText, gcloudbotConfig.Bot, &msg)
		return
	}
	gcloudbotConfig.InstanceConfig.Name = gcloudbotConfig.Update.Message.CommandArguments()

	// Create VM
	err := vms.CreateInstance(gcloudbotConfig.ComputeService, *gcloudbotConfig.Ctx, &gcloudbotConfig.InstanceConfig)
	if err != nil {
		// Define error message
		msg.Text = fmt.Sprintf("Couldn't create Instance: *%s*\n%s", gcloudbotConfig.InstanceConfig.Name, err.Error())
		// Log error
		log.Println(msg.Text)
		// Answer with error
		gcloudbotConfig.Bot.Send(msg)
		return
	}

	// Wait message
	msg.Text = fmt.Sprintf("Instance: *%s* created\n Wating for instance to start ...", gcloudbotConfig.InstanceConfig.Name)
	log.Println(msg.Text)
	gcloudbotConfig.Bot.Send(msg)

	// Wait for VM to start
	err = vms.PollForSerialOutput(gcloudbotConfig.ComputeService, *gcloudbotConfig.Ctx, &gcloudbotConfig.InstanceConfig, "DONE INITIALIZING STARTUP SCRIPT", "error is now")
	if err != nil {
		// Define error message
		msg.Text = fmt.Sprintf("Faild waiting for Instance: *%s* serial port\n%s\n", gcloudbotConfig.InstanceConfig.Name, err.Error())
		// Log error
		log.Println(msg.Text)
		// Answer with error
		gcloudbotConfig.Bot.Send(msg)
		return
	}

	// Get instance state
	instanceDetails, err := gcloudbotConfig.ComputeService.Instances.Get(gcloudbotConfig.InstanceConfig.ProjectID,
		gcloudbotConfig.InstanceConfig.Zone,
		gcloudbotConfig.InstanceConfig.Name).Do()
	if err != nil {
		// Define error message
		msg.Text = fmt.Sprintf("Faild getting Instance: *%s* state\n%s\n", gcloudbotConfig.InstanceConfig.Name, err.Error())
		// Log error
		log.Println(err)
		// Answer with error
		gcloudbotConfig.Bot.Send(msg)
		return
	}

	// results
	msg.Text = fmt.Sprintf("Instance: *%s* created succesfuly\n", gcloudbotConfig.InstanceConfig.Name)
	log.Println(msg.Text)
	gcloudbotConfig.Bot.Send(msg)

	msg.Text = fmt.Sprintf("Created VM State is *%s*: ", instanceDetails.Status)
	log.Println(msg.Text)
	gcloudbotConfig.Bot.Send(msg)
}
