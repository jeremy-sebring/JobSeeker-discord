package main

import (
	//"time"
	"log"

	"github.com/joho/godotenv"
	bot "sebring.dev/JobSeeker-discord/Bot/v2"
	JobHunter "sebring.dev/JobSeeker-discord/JobHunter/v2"
)

func InitEnv() {
	err := godotenv.Load(".env")
	checkNilErr(err)
}

func checkNilErr(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func main() {
	InitEnv()

	bot.CreateJobthreads(JobHunter.GetSerp())

	//bot.Run() // call the run function of bot/bot.go
}
