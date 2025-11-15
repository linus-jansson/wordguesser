package main

import (
	"encoding/json"
	"os"
	"unicode/utf8"
)

var allWords []string

func loadWords() {
	data, err := os.ReadFile("valid-words.json")
	if err != nil {
		panic(err)
	}

	var wordsFromFile []string
	if err := json.Unmarshal(data, &wordsFromFile); err != nil {
		panic(err)
	}

	for _, w := range wordsFromFile {
		if utf8.RuneCountInString(w) == 5 {
			allWords = append(allWords, w)
		}
	}
}
