package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type GuessResponse struct {
	Letters []CorrectLetter `json:"letters"`
	Error   string          `json:"error"`
}

type CorrectLetter int

const (
	Wrong   CorrectLetter = -1 // grey / black
	Kinda   CorrectLetter = 0  // yellow
	Correct CorrectLetter = 1  // green
)

type Pattern [5]CorrectLetter

func daysSinceOrdelEpochWithMagicNumber() int {
	// wordgame epoch is 2024-06-19
	YYYYMMDD := "2006-01-02"
	then, err := time.Parse(YYYYMMDD, "2024-06-19")
	check(err)
	now := time.Now()
	diff := now.Sub(then)
	magic_number := 904 + int(diff.Hours()/24) - 2
	return magic_number
}

func submitGuess(guess string) ([]CorrectLetter, error) {
	apiUrl := os.Getenv("API_URL")
	tryParam := fmt.Sprintf("n=%d", 0)
	guessParam := fmt.Sprintf("guess=%s", guess)
	id := fmt.Sprintf("id=%d", daysSinceOrdelEpochWithMagicNumber())
	url := fmt.Sprintf("%s?%s&%s&%s", apiUrl, tryParam, guessParam, id)
	fmt.Printf("Submitting '%s' to API\n", guess)
	// When calling the API I either get back a list of CorrectLetters '{"letters":[1,0,0,-1,-1]}', or I get a error '{"error":"INVALID_WORD"}'
	// The API seems to also always return 200 OK even on error.
	res, err := http.Get(url)
	check(err)
	defer res.Body.Close()

	var gr GuessResponse
	if err := json.NewDecoder(res.Body).Decode(&gr); err != nil {
		return nil, err
	}

	// Check for server-side “error”
	if gr.Error != "" {
		return nil, fmt.Errorf("server error: %s", gr.Error)
	}

	return gr.Letters, nil
}

func patternFor(guess, target string) Pattern {
	gr := []rune(guess)
	tr := []rune(target)

	if len(gr) != len(tr) {
		log.Printf("length mismatch: guess %q (%d runes), target %q (%d runes)",
			guess, len(gr), target, len(tr))
		return Pattern{}
	}

	n := len(gr)
	var pat Pattern
	counts := make(map[rune]int)

	// Count runes in target
	for i := 0; i < n; i++ {
		counts[tr[i]]++
	}

	// First pass: greens
	for i := 0; i < n; i++ {
		if gr[i] == tr[i] {
			pat[i] = Correct
			counts[gr[i]]--
		}
	}

	// Second pass: yellows / greys
	for i := 0; i < n; i++ {
		if pat[i] == Correct {
			continue
		}
		ch := gr[i]
		if counts[ch] > 0 {
			pat[i] = Kinda
			counts[ch]--
		} else {
			pat[i] = Wrong
		}
	}

	return pat
}

func patternFromSlice(s []CorrectLetter) Pattern {
	var p Pattern
	for i := 0; i < len(p) && i < len(s); i++ {
		p[i] = s[i]
	}
	return p
}

func matchesPattern(candidate, guess string, want Pattern) bool {
	got := patternFor(guess, candidate)
	for i := 0; i < len(want); i++ {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func filterCandidates(candidates []string, guess string, pat Pattern) []string {
	var out []string
	for _, c := range candidates {
		if matchesPattern(c, guess, pat) {
			out = append(out, c)
		}
	}
	return out
}

// scores a word by sum of frequencies of its unique letters
func scoreWordByLetterFrequency(word string, freq map[rune]int) int {
	seen := make(map[rune]bool)
	score := 0
	for _, r := range word {
		if !seen[r] {
			score += freq[r]
			seen[r] = true
		}
	}
	return score
}

func pickBestGuess(candidates, allWords []string) string {
	freq := make(map[rune]int)
	for _, w := range candidates {
		seen := make(map[rune]bool)
		for _, r := range w {
			if !seen[r] {
				freq[r]++
				seen[r] = true
			}
		}
	}

	search := candidates
	if len(search) == 0 {
		search = allWords
	}

	best := ""
	bestScore := -1

	for _, w := range search {
		s := scoreWordByLetterFrequency(w, freq)
		if s > bestScore {
			bestScore = s
			best = w
		}
	}

	return best
}
