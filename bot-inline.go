package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"smuel1414/gcloud.vms/gcloudbot"
	"smuel1414/gcloud.vms/structs"
	"smuel1414/gcloud.vms/utils"
	"smuel1414/gcloud.vms/vms"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var instanceConfig = vms.InstanceConfig{
	ProjectID:     "projectID",
	Zone:          "zone",
	Name:          "instanceName",
	StartupScript: utils.Script,
	MachineType:   "g1-small",
	ImageProject:  "debian-cloud",
	ImageFamily:   "debian-10",
	Scopes:        utils.ScopesForInst,
}

func main() {

	projectID := os.Getenv("GCLOUD_PROJECT_ID")
	jsonPath := os.Getenv("GCLOUD_SERVICE_ACCOUNT_JSON_PATH")
	zone := "us-central1-a"

	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		log.Fatal(err)
	}

	instanceConfig.ProjectID = projectID
	instanceConfig.Zone = zone

	botToken, _ := utils.GetenvStr("BOT_TOKEN")
	// workAccountID, _ := utils.GetenvInt64("MYWORKACCOUNT")
	accountID, _ := utils.GetenvInt64("MYACCOUNT")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
	}

	gcloudbotConfig := gcloudbot.GcloudbotConfig{
		ComputeService: computeService,
		InstanceConfig: instanceConfig,
		Ctx:            &ctx,
		Bot:            bot,
		Update:         nil,
	}

	fmt.Print(".")
	for update := range updates {
		gcloudbotConfig.Update = &update

		if update.CallbackQuery != nil {
			if !isAuthorized(bot, update.CallbackQuery.Message.Chat.ID, accountID) {
				continue
			}
			// fmt.Println(update.CallbackQuery.Message)

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			msgTextByte := []byte(update.CallbackQuery.Data)
			VMAction, err := structs.UnmarshalWelcome(msgTextByte)
			if err != nil {
				continue
			}

			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))

			log.Println("Is CallbackQuery")

			doAction(gcloudbotConfig, bot, &update, VMAction, ctx, computeService, instanceConfig)
		}
		if update.Message != nil {
			if !isAuthorized(bot, update.Message.Chat.ID, accountID) {
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			switch update.Message.Text {
			case "open":
				msg.ReplyMarkup = gcloudbot.NumericKeyboard

			}

			fmt.Println(msg)
			bot.Send(msg)

			fmt.Println(update.Message.IsCommand())
			if update.Message.IsCommand() {
				fmt.Println(update.Message.Command())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				switch update.Message.Command() {
				case "help":
					msg.Text = "type /vid or /status."
				case "status":
					msg.Text = "I'm ok."
				case "options":
					msg.Text = update.Message.Command()
					msg.ReplyMarkup = gcloudbot.NumericKeyboard
				case "clear":
					updates.Clear()
				case "create":
					gcloudbot.BotCreateInstance(gcloudbotConfig)
				default:
					msg.Text = "I don't know that command"
				}
				bot.Send(msg)
			}

		}
	}
}

func isAuthorized(bot *tgbotapi.BotAPI, id int64, authID int64) bool {
	if id != authID {
		msg := tgbotapi.NewMessage(id, "")
		msg.Text = "Unauthorized"
		bot.Send(msg)
		return false
	}
	return true
}

func doAction(gcloudbotConfig gcloudbot.GcloudbotConfig, bot *tgbotapi.BotAPI, update *tgbotapi.Update, action structs.VMAction, ctx context.Context, computeService *compute.Service, InstanceConfig vms.InstanceConfig) {
	chatID := update.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	switch action.Action {
	case "status-all":
		vmList, err := vms.GetVMs(computeService, instanceConfig.ProjectID, instanceConfig.Zone)
		if err != nil {
			log.Println(err)
			msg.Text = err.Error()
			bot.Send(msg)
		}

		if len(vmList) == 0 {
			msg.Text = fmt.Sprintf("No VMs for project: %s", instanceConfig.ProjectID)
			bot.Send(msg)
		}

		// status := string[]
		vmsState := ""

		for _, vm := range vmList {
			vmsState += fmt.Sprintf("%s - state is: %s\n", vm.Name, vm.Status)
		}

		msg.Text = vmsState
		bot.Send(msg)
	case "status":

		gcloudbotConfig.InstanceConfig.Name = action.VM
		gcloudbot.StatusAndDetails(gcloudbotConfig)

	case "delete-list":
		msg := tgbotapi.NewMessage(chatID, "")
		vmList, err := vms.GetVMs(computeService, instanceConfig.ProjectID, instanceConfig.Zone)
		if err != nil {
			log.Println(err)
			msg.Text = err.Error()
			bot.Send(msg)
		}

		var row []tgbotapi.InlineKeyboardButton

		for _, vm := range vmList {
			str := `{"vm": "` + vm.Name + `", "action": "delete"}`
			vmBotten := tgbotapi.InlineKeyboardButton{
				Text:         "Delete: " + vm.Name,
				CallbackData: &str,
			}
			row = append(row, vmBotten)
		}

		var vmListKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				row...,
			),
		)

		msg = tgbotapi.NewMessage(chatID, update.CallbackQuery.Data)
		msg.ReplyMarkup = vmListKeyboard
		bot.Send(msg)
		// bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
	case "status-list":
		// List All instances
		gcloudbot.StatusList(gcloudbotConfig)
	case "delete":
		gcloudbotConfig.InstanceConfig.Name = action.VM
		gcloudbot.Delete(gcloudbotConfig)
	case "delete-all":
		msg = tgbotapi.NewMessage(chatID, update.CallbackQuery.Data)

		vmList, err := vms.GetVMs(computeService, instanceConfig.ProjectID, instanceConfig.Zone)
		if err != nil {
			log.Println(err)
			msg.Text = err.Error()
			bot.Send(msg)
		}

		for _, vm := range vmList {
			instanceConfig.Name = vm.Name
			if instanceConfig.Name == os.Getenv("BOT_VM") {
				continue
			}
			vms.DeleteInstance(computeService, ctx, &instanceConfig)

			_, err := vms.GetVMStatus(computeService, instanceConfig.ProjectID, instanceConfig.Zone, instanceConfig.Name)
			if err != nil {
				msg.Text = err.Error() + fmt.Sprintf("VM %s deleted succesfuly\n", instanceConfig.Name)
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, msg.Text))
				bot.Send(msg)
				log.Println(err)
				log.Println(msg.Text)
			}
			msg.Text = "Finished delete all vms"
			bot.Send(msg)
		}
	default:

	}
}
