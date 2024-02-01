package hangweb

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const User = "admin"

var Mdp = "root"
var Username = ""
var WordOfTheDay = "mot"

func login(utilisateur string, mot_de_passe string) bool {
	if utilisateur == User && mot_de_passe == Mdp {
		return true
	} else {
		return false
	}
}

// Fonction pour changer le mot du jour
func ChangeWordOfTheDay(word string) {
	WordOfTheDay = word
}

// Fonction pour reset le scoreboard
func resetScoreboard() {
	file, err := os.OpenFile("./saves/save.txt", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	var previousLineEmpty bool

	for scanner.Scan() {
		line := scanner.Text()

		// Si la ligne commence par "Score:", on remplace par "Score: 0"
		if strings.HasPrefix(line, "Score:") {
			line = "Score: 0"
		}

		// Laisser une ligne vide a la fin
		if line != "" || !previousLineEmpty {
			lines = append(lines, line)
		}

		previousLineEmpty = line == ""
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	// Tronquer le fichier pour éliminer les données restantes
	if err := file.Truncate(0); err != nil {
		fmt.Println(err)
	}

	// Replacer le curseur au début du fichier
	file.Seek(0, 0)

	// Écrire les lignes update dans le fichier
	writer := bufio.NewWriter(file)
	for i, line := range lines {
		fmt.Fprintln(writer, line)
		if i == len(lines)-1 {
			fmt.Fprint(writer, "")
		}
	}
	writer.Flush()
}

// Fonction pour reset tous les joueurs
func resetAll(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	fmt.Println("Reset done.")
}

// Fonction pour reset le shop
func resetShop() {
	file, err := os.OpenFile("./saves/inventory.txt", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	var previousLineEmpty bool

	for scanner.Scan() {
		line := scanner.Text()

		// Si la ligne commence par "Score:", on remplace par "Score: 0"
		if strings.HasPrefix(line, "Spell1:") {
			line = "Spell1: 0"
		}
		// Si la ligne commence par "Score:", on remplace par "Score: 0"
		if strings.HasPrefix(line, "Spell2:") {
			line = "Spell2: 0"
		} // Si la ligne commence par "Score:", on remplace par "Score: 0"
		if strings.HasPrefix(line, "Spell3:") {
			line = "Spell3: 0"
		}

		// Laisser une ligne vide a la fin
		if line != "" || !previousLineEmpty {
			lines = append(lines, line)
		}

		previousLineEmpty = line == ""
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	// Tronquer le fichier pour éliminer les données restantes
	if err := file.Truncate(0); err != nil {
		fmt.Println(err)
	}

	// Replacer le curseur au début du fichier
	file.Seek(0, 0)

	// Écrire les lignes update dans le fichier
	writer := bufio.NewWriter(file)
	for i, line := range lines {
		fmt.Fprintln(writer, line)
		if i == len(lines)-1 {
			fmt.Fprint(writer, "")
		}
	}
	writer.Flush()
}
