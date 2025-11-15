package main

import (
	"fmt"
	"os"
)

func writeResults(correctWord string, attempts int, success bool) {
	status := "was found"
	if correctWord == "" {
		status = "was not found"
	}
	stringToWrite := fmt.Sprintf("Word \"%s\" %s in %d tries\n", correctWord, status, attempts)
	os.WriteFile("output.txt", []byte(stringToWrite), 0644)
}
