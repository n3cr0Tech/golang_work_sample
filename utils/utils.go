package utils

import (
	"fmt"
	"log"

	// "net/http"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	// "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var EnvEntries map[string]string

func GetEnvEntries() map[string]string {
	EnvEntries = loadDotEnv()
	return EnvEntries
}

func GetRabbitMQURL() string {
	return "amqp://" + EnvEntries["RABBITMQ_USER"] + ":" + EnvEntries["RABBITMQ_PWD"] + "@" + EnvEntries["RABBITMQ_HOST"]
}

// fileName example: "config.json"
func GetConfigFromJSON(fileDir string, fileName string) map[string]string {
	if len(fileDir) == 0 {
		fileDir = getCurrentDir()
	}
	filePath := filepath.Join(fileDir, fileName)
	log.Println("- Attempting to open file at: ")
	log.Println(filePath)
	file := openFile(filePath)
	byteValue, _ := io.ReadAll(file)

	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)
	defer file.Close()
	return result
}

func loadDotEnv() map[string]string {
	err := godotenv.Load("../.env") // .env file is at root dir of project
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	keys := make(map[string]string)
	keys["JWT_SECRET"] = os.Getenv("JWT_SECRET")
	keys["AUTH_HEADER"] = os.Getenv("AUTH_HEADER")
	keys["RABBITMQ_USER"] = os.Getenv("RABBITMQ_USER")
	keys["RABBITMQ_PWD"] = os.Getenv("RABBITMQ_PWD")
	keys["RABBITMQ_HOST"] = os.Getenv("RABBITMQ_HOST")
	keys["MONGO_URL"] = os.Getenv("MONGO_URL")
	keys["DB_NAME"] = os.Getenv("DB_NAME")
	return keys
}

func getCurrentDir() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
	return path
}

func openFile(filePath string) *os.File {
	// Open our jsonFile
	jsonFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	return jsonFile
}

/*
--- SAMPLE REQUEST ---
curl http://localhost:8080/greeting \
    --include --header \
    "Content-Type: application/json" \
    --request "POST" --data \
    '{"name": "Luke"}'
*/
