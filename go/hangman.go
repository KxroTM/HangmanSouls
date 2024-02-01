package hangweb

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Hangman struct {
	TryLeft       int
	WordHidden    string
	Word          string
	Guessedletter []string
	GameState     int
	Username      string
	Difficulty    string
	Items         bool
	Spell1        int
	Spell2        int
	Spell3        int
}

var EasyMult = 10
var NormalMult = 20
var HardMult = 30
var WinStreak int
var hang Hangman

func StartGame(filename string) {
	hang.Guessedletter = []string{}
	PickWord(filename)
}

func StartGameHangOfTheDay(username string) {
	HangmanStart(username, "extra")
	hang.Guessedletter = []string{}
	hang.Word = WordOfTheDay
	hang.WordHidden = wordToUnderScore()
}

func PickWord(filename string) {
	tw := ReadWord(filename)
	hang.Word = tw[rand.Intn(len(tw))]
	hang.WordHidden = wordToUnderScore()
}

func ReadWord(filename string) []string {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	word1 := ""
	var todowordos []string
	for _, char := range string(file) {
		if char == '\n' {
			todowordos = append(todowordos, word1)
			word1 = ""
		} else {
			word1 += string(char)
		}
	}
	return todowordos
}

func wordToUnderScore() string {
	sampleRegexp := regexp.MustCompile("[a-z,A-Z]")

	input := hang.Word

	result := sampleRegexp.ReplaceAllString(input, "_")
	return (string(result))
}

func FindAndReplace(letterToReplace string) {
	isALetter, err := regexp.MatchString("^[a-zA-Z]", letterToReplace)
	if !isALetter || err != nil {
		return
	}

	if len(letterToReplace) != 0 {
		for _, guess := range hang.Guessedletter {
			if letterToReplace == guess {
				return
			}
		}
		hang.Guessedletter = append(hang.Guessedletter, letterToReplace)
	}
	if len(letterToReplace) > 1 {
		if letterToReplace == hang.Word {
			hang.WordHidden = hang.Word
		} else {
			hang.TryLeft -= 2
		}
		if hang.TryLeft < 0 {
			hang.TryLeft = 0
		}
		return
	}

	isFound := strings.Index(hang.Word, letterToReplace)
	if isFound == -1 {
		if hang.TryLeft >= 1 {
			hang.TryLeft--

		}

	} else {
		str3 := []rune(hang.WordHidden)
		for i, lettre := range hang.Word {
			if string(lettre) == letterToReplace {
				str3[i] = lettre
				hang.WordHidden = string(str3)
			}
		}
	}
}

func testEndGame() {
	if hang.WordHidden == hang.Word {
		hang.GameState = 1
	} else if hang.TryLeft <= 0 {
		hang.GameState = 2
	}
}

func HangmanStart(username string, difficulty string) {
	hang = Hangman{
		TryLeft:       10,
		Guessedletter: []string{},
		GameState:     0,
		Username:      username,
		Difficulty:    difficulty,
		Items:         false,
	}
}

func CheckDifficultyMult(difficulty string) int {
	if difficulty == "easy" {
		return EasyMult
	} else if difficulty == "normal" {
		return NormalMult
	}
	return HardMult
}

// Fonction pour écrire la structure dans un fichier texte
func writeStructToFile(s Player) {
	file, err := os.OpenFile("./saves/save.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		fieldValue := fmt.Sprintf("%v", field.Interface())
		line := fmt.Sprintf("%s: %s\n", fieldName, fieldValue)
		_, err := file.WriteString(line)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Fonction pour vérifier si un username existe dans la save
func isUsernameInSave(username string, filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Username: "+username) {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return false
}

// Fonction pour retourner le numéro de la ligne si un username est trouvé dans la save
func findUsernameLineNumber(username string, filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ligne := 1

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Username: "+username) {
			return ligne
		}
		ligne++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return -1
}

// Fonction pour lire les valeurs d'une structure à partir d'une ligne spécifiée
func playerBuild(ligne int) Player {
	file, err := os.Open("./saves/save.txt")
	if err != nil {
		return Player{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nbLigne := 0

	var user Player

	for scanner.Scan() {
		nbLigne++
		if nbLigne < ligne {
			continue
		}
		line := scanner.Text()
		if nbLigne == ligne {
			user.Username = getValueFromLine("Username", line)
		} else if nbLigne == ligne+1 {
			user.Score = getIntValueFromLine("Score", line)
		} else if nbLigne == ligne+2 {
			user.Money = getIntValueFromLine("Money", line)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return Player{}
	}

	return user
}

// Fonction pour obtenir la valeur d'un attribut de la struct à partir d'une ligne
func getValueFromLine(attribut, line string) string {
	parts := strings.SplitN(line, ": ", 2)
	if len(parts) == 2 && parts[0] == attribut {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

// Fonction pour obtenir la valeur entière d'un attribut de la struct partir d'une ligne
func getIntValueFromLine(attribut, line string) int {
	valueStr := getValueFromLine(attribut, line)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0
	}
	return value
}

// Fonction pour update les attributs de la structure dans un fichier texte
func writeStruct(s Player, ligne int) {
	file, err := os.OpenFile("./saves/save.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	scanner := bufio.NewScanner(file)
	lignes := []string{}

	// Lire toutes les lignes existantes
	for scanner.Scan() {
		lignes = append(lignes, scanner.Text())
	}

	// Si startLine est en dehors de la plage des lignes existantes, ne rien faire
	if ligne < 0 || ligne >= len(lignes) {
		fmt.Println("startLine hors de la plage des lignes existantes.")
		return
	}

	// Remplacer les données à partir de ligne
	for i := ligne; i < ligne+val.NumField(); i++ {
		field := val.Field(i - ligne)
		fieldName := typ.Field(i - ligne).Name
		fieldValue := fmt.Sprintf("%v", field.Interface())
		line := fmt.Sprintf("%s: %s", fieldName, fieldValue)

		if i < len(lignes) {
			lignes[i] = line
		} else {
			lignes = append(lignes, line)
		}
	}

	// Réécrire toutes les lignes dans le fichier
	file.Seek(0, 0)
	file.Truncate(0)
	writer := bufio.NewWriter(file)
	for _, line := range lignes {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()
}

func revealLetter(word, maskedWord string) string {

	var hiddenIndices []int
	for i, char := range maskedWord {
		if char == '_' {
			hiddenIndices = append(hiddenIndices, i)
		}
	}

	if len(hiddenIndices) == 0 {
		return maskedWord
	}

	// Choix aléatoire d'un indice pour révéler une lettre
	randomIndex := hiddenIndices[rand.Intn(len(hiddenIndices))]

	revealedWord := strings.Builder{}
	for i, char := range maskedWord {
		if i == randomIndex {
			revealedWord.WriteRune(rune(word[i]))
		} else {
			revealedWord.WriteRune(char)
		}
	}

	return revealedWord.String()
}
