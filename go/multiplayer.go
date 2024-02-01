package hangweb

import (
	"math/rand"
	"regexp"
	"strings"
)

type HangmanMulti struct {
	TryLeft1       int
	WordHidden1    string
	Word1          string
	Guessedletter1 []string
	GameState      int
	Username1      string
	TryLeft2       int
	WordHidden2    string
	Word2          string
	Guessedletter2 []string
	Username2      string
	Turn           int
	Items          bool
}

var multiplayer HangmanMulti

func multiplayerStart() {
	multiplayer = HangmanMulti{
		TryLeft1:       10,
		Guessedletter1: []string{},
		GameState:      0,
		Username1:      "Player 1",
		TryLeft2:       10,
		Guessedletter2: []string{},
		Username2:      "Player 2",
		Turn:           1,
		Items:          false,
	}
}

func StartGameMultiplayer(filename string) {
	multiplayer.Guessedletter1 = []string{}
	multiplayer.Guessedletter2 = []string{}
	PickWordMulti(filename)
}

func PickWordMulti(filename string) {
	tw := ReadWord(filename)
	multiplayer.Word1 = tw[rand.Intn(len(tw))]
	multiplayer.Word2 = tw[rand.Intn(len(tw))]
	multiplayer.WordHidden1 = wordToUnderScorePlayer1()
	multiplayer.WordHidden2 = wordToUnderScorePlayer2()
}

func wordToUnderScorePlayer1() string {
	sampleRegexp := regexp.MustCompile("[a-z,A-Z]")

	input := multiplayer.Word1

	result := sampleRegexp.ReplaceAllString(input, "_")
	return (string(result))
}

func wordToUnderScorePlayer2() string {
	sampleRegexp := regexp.MustCompile("[a-z,A-Z]")

	input := multiplayer.Word2

	result := sampleRegexp.ReplaceAllString(input, "_")
	return (string(result))
}

func findAndReplacePlayer1(letterToReplace string) {
	isALetter, err := regexp.MatchString("^[a-zA-Z]", letterToReplace)
	if !isALetter || err != nil {
		return
	}

	if len(letterToReplace) != 0 {
		for _, guess := range multiplayer.Guessedletter1 {
			if letterToReplace == guess {
				return
			}
		}
		multiplayer.Guessedletter1 = append(multiplayer.Guessedletter1, letterToReplace)
	}
	if len(letterToReplace) > 1 {
		if letterToReplace == multiplayer.Word1 {
			multiplayer.WordHidden1 = multiplayer.Word1
		} else {
			multiplayer.TryLeft1 -= 2
		}
		if multiplayer.TryLeft1 < 0 {
			multiplayer.TryLeft1 = 0
		}
		return
	}

	isFound := strings.Index(multiplayer.Word1, letterToReplace)
	if isFound == -1 {
		if multiplayer.TryLeft1 >= 1 {
			multiplayer.TryLeft1--

		}

	} else {
		str3 := []rune(multiplayer.WordHidden1)
		for i, lettre := range multiplayer.Word1 {
			if string(lettre) == letterToReplace {
				str3[i] = lettre
				multiplayer.WordHidden1 = string(str3)
			}
		}
	}
}

func findAndReplacePlayer2(letterToReplace string) {
	isALetter, err := regexp.MatchString("^[a-zA-Z]", letterToReplace)
	if !isALetter || err != nil {
		return
	}

	if len(letterToReplace) != 0 {
		for _, guess := range multiplayer.Guessedletter2 {
			if letterToReplace == guess {
				return
			}
		}
		multiplayer.Guessedletter2 = append(multiplayer.Guessedletter2, letterToReplace)
	}
	if len(letterToReplace) > 1 {
		if letterToReplace == multiplayer.Word2 {
			multiplayer.WordHidden2 = multiplayer.Word2
		} else {
			multiplayer.TryLeft2 -= 2
		}
		if multiplayer.TryLeft2 < 0 {
			multiplayer.TryLeft2 = 0
		}
		return
	}

	isFound := strings.Index(multiplayer.Word2, letterToReplace)
	if isFound == -1 {
		if multiplayer.TryLeft2 >= 1 {
			multiplayer.TryLeft2--

		}

	} else {
		str3 := []rune(multiplayer.WordHidden2)
		for i, lettre := range multiplayer.Word2 {
			if string(lettre) == letterToReplace {
				str3[i] = lettre
				multiplayer.WordHidden2 = string(str3)
			}
		}
	}
}

func testEndGameMulti() {
	if multiplayer.WordHidden1 == multiplayer.Word1 {
		multiplayer.GameState = 1
	} else if multiplayer.TryLeft1 <= 0 {
		multiplayer.GameState = 2
	}
	if multiplayer.WordHidden2 == multiplayer.Word2 {
		multiplayer.GameState = 1
	} else if multiplayer.TryLeft2 <= 0 {
		multiplayer.GameState = 2
	}
}
