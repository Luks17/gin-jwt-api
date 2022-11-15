package api

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetUri() string {
	LoadEnv()
	// loads MONGODB_URI Environment Variable
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Missing MONGODB_URI environment variable")
	}

	return uri
}

func LoadEnv() {
	// loads .env file to Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No env file found")
	}
}

func GetSecret() string {
	LoadEnv()
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatal("Missing SECRET_KEY environment variable")
	}

	return secret
}
