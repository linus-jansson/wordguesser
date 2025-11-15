package main

import "fmt"

type CorrectGuess struct {
	Word string
	Try  int
}

func play(allWords, answers []string) CorrectGuess {
	candidates := append([]string(nil), answers...)

	for turn := 1; turn <= 6; turn++ {
		guess := pickBestGuess(candidates, allWords)
		fmt.Printf("Turn %d, guess: %s\n", turn, guess)

		apiPatternSlice, err := submitGuess(guess)
		if err != nil {
			fmt.Println("probably wrong word", err)
		}

		pat := patternFromSlice(apiPatternSlice)

		// Check for solved (all Correct)
		solved := true
		for _, v := range pat {
			if v != Correct {
				solved = false
				break
			}
		}
		if solved {
			fmt.Println("Solved in", turn, "guesses!")
			return CorrectGuess{
				Word: guess,
				Try:  turn,
			}
		}

		// Narrow candidates for next turn
		candidates = filterCandidates(candidates, guess, pat)
		fmt.Printf("Remaining candidates: %d\n", len(candidates))

		if len(candidates) == 0 {
			fmt.Println("No candidates left – something is inconsistent.")
			return CorrectGuess{}
		}
	}

	fmt.Println("Out of guesses :(")
	return CorrectGuess{}
}
