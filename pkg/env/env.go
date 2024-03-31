package env

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func InitEnvironment() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Fatal(err)
	}
}

func must(key string) string {
	if val, exist := os.LookupEnv(key); exist {
		return val
	}

	log.Fatalf("no %s in env", key)
	return ""
}

func MustOpenAIAPIKey() string {
	return must("OPENAI_API_KEY")
}

func MustTGBotTOKEN() string {
	return must("TG_BOT_TOKEN")
}
