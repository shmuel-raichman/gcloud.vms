package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"smuel1414/gcloud.vms/structs"
	"smuel1414/gcloud.vms/utils"
	"smuel1414/gcloud.vms/vms"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	script string = `

	export GITHUB_PASSWORD=` + os.Getenv("GITHUB_PASSWORD") + `
	export GITHUB_USERNAME=` + os.Getenv("GITHUB_USERNAME") + `
	export GITHUB_USERNAME=` + os.Getenv("GITHUB_INITAL_REPO") + `
	export VM_GCLOUD_USER= ` + os.Getenv("VM_GCLOUD_USER") + `
	export VM_SSH_USER= ` + os.Getenv("VM_SSH_USER") + `
	export DOCKER_COMPOSE_VERSION= ` + os.Getenv("DOCKER_COMPOSE_VERSION") + `

	sudo apt-get install git tree vim -y
	cd /opt
	mkdir init
	cd init
	git clone https://$GITHUB_USERNAME:$GITHUB_PASSWORD@github.com/$GITHUB_USERNAME/$GITHUB_INITAL_REPO

	cd $GITHUB_INITAL_REPO
	chmod +x installdocker.sh
	./installdocker.sh

	sudo usermod -aG docker $VM_GCLOUD_USER
	sudo usermod -aG docker $VM_SSH_USER

	echo "DONE INITIALIZING STARTUP SCRIPT"
	`
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Status all", `{"vm": "all", "action": "status-all"}`),
		tgbotapi.NewInlineKeyboardButtonData("delete", `{"vm": "all", "action": "delete-list"}`),
		tgbotapi.NewInlineKeyboardButtonData("Status List", `{"vm": "all", "action": "status-list"}`),
	),
)

var scopesForInst = []string{
	"https://www.googleapis.com/auth/devstorage.read_only",
	"https://www.googleapis.com/auth/logging.write",
	"https://www.googleapis.com/auth/monitoring.write",
	"https://www.googleapis.com/auth/servicecontrol",
	"https://www.googleapis.com/auth/service.management.readonly",
	"https://www.googleapis.com/auth/trace.append",
}

var instanceConfig = vms.InstanceConfig{
	ProjectID:     "projectID",
	Zone:          "zone",
	Name:          "instanceName",
	StartupScript: script,
	MachineType:   "g1-small",
	ImageProject:  "debian-cloud",
	ImageFamily:   "debian-10",
	Scopes:        scopesForInst,
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

	fmt.Print(".")
	for update := range updates {

		if update.CallbackQuery != nil {
			if !isAuthorized(bot, update.CallbackQuery.Message.Chat.ID, accountID) {
				continue
			}
			fmt.Println(update.CallbackQuery.Message)

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			msgTextByte := []byte(update.CallbackQuery.Data)
			VMAction, err := structs.UnmarshalWelcome(msgTextByte)
			if err != nil {
				continue
			}

			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))

			log.Println("Is CallbackQuery")

			doAction(bot, &update, VMAction, ctx, computeService, instanceConfig)
		}
		if update.Message != nil {
			if !isAuthorized(bot, update.Message.Chat.ID, accountID) {
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			switch update.Message.Text {
			case "open":
				msg.ReplyMarkup = numericKeyboard

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
					msg.ReplyMarkup = numericKeyboard
				case "clear":
					updates.Clear()
				case "create":
					msg.Text = "You supplied the following argument: " + update.Message.CommandArguments()
					instanceConfig.Name = update.Message.CommandArguments()
					err := vms.CreateInstance(computeService, ctx, &instanceConfig)
					if err != nil {
						log.Println(err)
					}

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					vms.CreateInstance(computeService, ctx, &instanceConfig)
					status, err := vms.GetVMStatus(computeService, projectID, zone, instanceConfig.Name)
					if err != nil {
						log.Println(err)
						msg.Text = err.Error()
						bot.Send(msg)
						continue
					}
					msg.Text = fmt.Sprintf("VM State is %s", status)
					bot.Send(msg)

					err = vms.PollForSerialOutput(computeService, ctx, &instanceConfig, "DONE INITIALIZING STARTUP SCRIPT", "error is now")
					if err != nil {
						log.Println(err)
						msg.Text = err.Error()
						bot.Send(msg)
						continue
					}

					msg.Text = fmt.Sprintf("Creating vm %s", instanceConfig.Name)
					bot.Send(msg)

					status, err = vms.GetVMStatus(computeService, projectID, zone, instanceConfig.Name)
					if err != nil {
						log.Println(err)
						msg.Text = err.Error()
						bot.Send(msg)
						continue
					}
					msg.Text = fmt.Sprintf("VM State is %s", status)
					bot.Send(msg)

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

func doAction(bot *tgbotapi.BotAPI, update *tgbotapi.Update, action structs.VMAction, ctx context.Context, computeService *compute.Service, InstanceConfig vms.InstanceConfig) {
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
		msg = tgbotapi.NewMessage(chatID, update.CallbackQuery.Data)
		instanceConfig.Name = action.VM
		status, err := vms.GetVMStatus(computeService, instanceConfig.ProjectID, instanceConfig.Zone, instanceConfig.Name)
		if err != nil {
			msg.Text = err.Error() + fmt.Sprintf("VM %s deleted succesfuly\n", instanceConfig.Name)
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, msg.Text))
			bot.Send(msg)
			log.Println(err)
			log.Println(msg.Text)
		}
		msg.Text = fmt.Sprintf("VM State is %s\n", status)
		bot.Send(msg)

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
		msg := tgbotapi.NewMessage(chatID, "")
		vmList, err := vms.GetVMs(computeService, instanceConfig.ProjectID, instanceConfig.Zone)
		if err != nil {
			log.Println(err)
			msg.Text = err.Error()
			bot.Send(msg)
		}

		var row []tgbotapi.InlineKeyboardButton

		for _, vm := range vmList {
			str := `{"vm": "` + vm.Name + `", "action": "status"}`
			vmBotten := tgbotapi.InlineKeyboardButton{
				Text:         "Status: " + vm.Name,
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
	case "delete":
		msg = tgbotapi.NewMessage(chatID, update.CallbackQuery.Data)
		instanceConfig.Name = action.VM
		vms.DeleteInstance(computeService, ctx, &instanceConfig)
		status, err := vms.GetVMStatus(computeService, instanceConfig.ProjectID, instanceConfig.Zone, instanceConfig.Name)
		if err != nil {
			msg.Text = err.Error() + fmt.Sprintf("VM %s deleted succesfuly\n", instanceConfig.Name)
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, msg.Text))
			bot.Send(msg)
			log.Println(err)
			log.Println(msg.Text)
		}
		msg.Text = fmt.Sprintf("VM State is %s\n", status)
		bot.Send(msg)
	case "create":
		// vmList, err := vms.GetVMs(computeService, instanceConfig.ProjectID, instanceConfig.Zone)
		// if err != nil {
		// 	log.Println(err)
		// 	msg.Text = err.Error()
		// 	bot.Send(msg)
		// }

		// var row []tgbotapi.InlineKeyboardButton

		// for _, vm := range vmList {
		// 	vmBotten := tgbotapi.NewInlineKeyboardButtonData("Delete: "+vm.Name, `{"vm": "`+vm.Name+`", "action": "delete"}`)
		// 	row = append(row, vmBotten)
		// }

		// var vmListKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		// 	tgbotapi.NewInlineKeyboardRow(
		// 		row...,
		// 	),
		// )
		// msg.ReplyMarkup = vmListKeyboard
		bot.Send(msg)
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

	default:

	}
}
