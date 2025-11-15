package main

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"io"
	"log"
	"unicode/utf8"
)

//go:embed valid-words.json.gz
var embeddedWordsGz []byte

var allWords []string

func loadWords() {
	r, err := gzip.NewReader(bytes.NewReader(embeddedWordsGz))
	if err != nil {
		log.Fatalf("failed to create gzip reader: %v", err)
	}
	defer r.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		log.Fatalf("failed to decompress embedded words: %v", err)
	}

	if err := json.Unmarshal(buf.Bytes(), &allWords); err != nil {
		log.Fatalf("failed to unmarshal embedded words: %v", err)
	}

	filtered := allWords[:0]
	for _, w := range allWords {
		if utf8.RuneCountInString(w) == 5 {
			filtered = append(filtered, w)
		}
	}
	allWords = filtered
}
