package hangweb

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
)

var item1 = 50
var item2 = 30
var item3 = 15

type Shop struct {
	Username string
	Spell1   int
	Spell2   int
	Spell3   int
}

var inventory Shop

// Fonction pour écrire la structure dans un fichier texte
func writeInventoryToFile(inv Shop) {
	file, err := os.OpenFile("./saves/inventory.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	val := reflect.ValueOf(inv)
	typ := reflect.TypeOf(inv)

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

// Fonction pour lire les valeurs d'une structure à partir d'une ligne spécifiée
func shopBuild(ligne int) Shop {
	file, err := os.Open("./saves/inventory.txt")
	if err != nil {
		return Shop{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nbLigne := 0

	var inventory Shop

	for scanner.Scan() {
		nbLigne++
		if nbLigne < ligne {
			continue
		}
		line := scanner.Text()
		if nbLigne == ligne {
			inventory.Username = getValueFromLine("Username", line)
		} else if nbLigne == ligne+1 {
			inventory.Spell1 = getIntValueFromLine("Spell1", line)
		} else if nbLigne == ligne+2 {
			inventory.Spell2 = getIntValueFromLine("Spell2", line)
		} else if nbLigne == ligne+3 {
			inventory.Spell3 = getIntValueFromLine("Spell3", line)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return Shop{}
	}

	return inventory
}

// Fonction pour update les attributs de la structure dans un fichier texte
func updateShopData(s Shop, ligne int) {
	file, err := os.OpenFile("./saves/inventory.txt", os.O_RDWR|os.O_CREATE, 0644)
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

// Fonction qui vérifie si on peut acheter un item
func enoughtMoney(money int, item int) bool {
	return money >= item
}
