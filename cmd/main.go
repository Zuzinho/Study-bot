package main

import (
	"StudyTGBot/pkg/env"
	"StudyTGBot/pkg/handler"
	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	log.Debug("Initializing environment vars...")
	env.InitEnvironment()

	log.Debug("Getting TG_BOT_TOKEN...")
	tgBotToken := env.MustTGBotTOKEN()
	log.Debug("Getting IPENAI_API_KEY...")
	openAIAPIKey := env.MustOpenAIAPIKey()

	log.Debug("Creatign openai client...")
	client := openai.NewClient(openAIAPIKey)
	h := handler.NewOpenAIHandler(client)

	log.Debug("Trying connect to tg bot...")
	botApi, err := api.NewBotAPI(tgBotToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Started bot with name %s", botApi.Self.UserName)

	log.Debug("Getting updates chan...")
	u := api.NewUpdate(0)
	updates, err := botApi.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Ready for work")

	for update := range updates {
		log.Infof("Getting update '%d'", update.UpdateID)

		msg := new(api.MessageConfig)

		if message := update.Message; message != nil {
			chatID := message.Chat.ID

			log.Infof("Gettign message from chat '%d' '%s' with text: %s",
				chatID, message.From.UserName, message.Text)

			switch message.Text {
			case "/start":
				*msg = api.NewMessage(chatID, "Здравствуйте, отправьте мне вопрос, и я дам вам ответ)")
			default:
				log.Info("Trying handle message with openai")

				msg, err = h.HandleMessage(chatID, message.Text)
				if err != nil {
					log.Error(err)

					*msg = api.NewMessage(chatID,
						"Произошла ошибка при попытке сгенерировать ответ. Попробуйте заново")
				}
			}
		} else {
			log.Info("Unknown request type")

			*msg = api.NewMessage(update.Message.Chat.ID,
				"Неизвестный тип запроса. Введите пожалуйста текст своего вопроса")
		}

		if _, err = botApi.Send(*msg); err != nil {
			log.Println(err)
		}
	}
}
