package tgd

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT         int
	FRONTEND_URL string
}

var AppConfig = Config{}

func LoadConfigs() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln(".env file is required to run the backend", err)
	}
	port := os.Getenv("PORT")
	if portInt, err := strconv.ParseInt(port, 10, 32); err != nil {
		log.Fatalln("Invalid PORT value:", err)
	} else {
		AppConfig.PORT = int(portInt)
	}
	AppConfig.FRONTEND_URL = os.Getenv("FRONTEND_URL")
}
