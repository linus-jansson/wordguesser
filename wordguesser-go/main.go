package main

func main() {

	loadWords()

	resulting := play(allWords, allWords)

	writeResults(resulting.Word, resulting.Try, resulting.Word != "")

	SendFileToDiscordWebhook("output.txt", getDiscordWebhookURL())
}
