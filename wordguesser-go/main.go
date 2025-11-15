package main

import "os"

func ensureEnvVariables() {
	a := os.Getenv("API_URL")
	if a == "" {
		panic("API_URL environment variable not set")
	}
	b := os.Getenv("DISCORD_WEBHOOK_URL")
	if b == "" {
		panic("DISCORD_WEBHOOK_URL environment variable not set")
	}
}

func main() {
	ensureEnvVariables()

	loadWords()

	resulting := play(allWords, allWords)

	writeResults(resulting.Word, resulting.Try, resulting.Word != "")

	SendFileToDiscordWebhook("output.txt", getDiscordWebhookURL())
}
