package hangweb

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ScoreB struct {
	Username []string
	Score    []int
}

type Player struct {
	Username string
	Score    int
	Money    int
}

var user Player
var doublePoint = 1
var scores ScoreB

// Fonction qui permet de trier le scoreboard et de l'ajouter a une struct pour l'afficher
func readScoresFromFile() ScoreB {
	var scores ScoreB

	file, err := os.Open("./saves/save.txt")
	if err != nil {
		return scores
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Username:") {
			username := strings.TrimSpace(strings.TrimPrefix(line, "Username:"))
			scores.Username = append(scores.Username, username)
		} else if strings.HasPrefix(line, "Score:") {
			scoreStr := strings.TrimSpace(strings.TrimPrefix(line, "Score:"))
			score, err := strconv.Atoi(scoreStr)
			if err != nil {
				return scores
			}
			scores.Score = append(scores.Score, score)
		}
	}

	if err := scanner.Err(); err != nil {
		return scores
	}

	// Structure temporaire
	var tempData []struct {
		Username string
		Score    int
	}

	for i, username := range scores.Username {
		tempData = append(tempData, struct {
			Username string
			Score    int
		}{Username: username, Score: scores.Score[i]})
	}

	// Tri des scores
	sort.SliceStable(tempData, func(i, j int) bool {
		return tempData[i].Score > tempData[j].Score
	})

	// Update struct
	for i, data := range tempData {
		scores.Username[i] = data.Username
		scores.Score[i] = data.Score
	}

	return scores
}
