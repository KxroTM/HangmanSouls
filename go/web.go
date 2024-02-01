package hangweb

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var home = template.Must(template.ParseFiles("./src/templates/eng/home.html"))
var hangman = template.Must(template.ParseFiles("./src/templates/eng/hangman.html"))
var hangoftheday = template.Must(template.ParseFiles("./src/templates/eng/hangoftheday.html"))
var beta = template.Must(template.ParseFiles("./src/templates/eng/beta.html"))
var Register = template.Must(template.ParseFiles("./src/templates/eng/register.html"))
var Admin = template.Must(template.ParseFiles("./src/templates/eng/admin.html"))
var AdminHOTD = template.Must(template.ParseFiles("./src/templates/eng/hang-of-the-day-admin.html"))
var Account = template.Must(template.ParseFiles("./src/templates/eng/account.html"))
var Manage = template.Must(template.ParseFiles("./src/templates/eng/manage-admin.html"))
var ManageSB = template.Must(template.ParseFiles("./src/templates/eng/scoreboard-admin.html"))
var ManageShop = template.Must(template.ParseFiles("./src/templates/eng/shop-admin.html"))
var Error = template.Must(template.ParseFiles("./src/templates/eng/error.html"))
var Login = template.Must(template.ParseFiles("./src/templates/eng/login.html"))
var WrongLogin = template.Must(template.ParseFiles("./src/templates/eng/wronglogin.html"))
var WrongRegister = template.Must(template.ParseFiles("./src/templates/eng/wrongregister.html"))
var Scoreboard = template.Must(template.ParseFiles("./src/templates/eng/scoreboard.html"))
var shop = template.Must(template.ParseFiles("./src/templates/eng/shop.html"))
var registerShop = template.Must(template.ParseFiles("./src/templates/eng/registershop.html"))
var Multiplayer = template.Must(template.ParseFiles("./src/templates/eng/multiplayer.html"))
var RegisterCustom = template.Must(template.ParseFiles("./src/templates/eng/registercustom.html"))
var WrongRegisterCustom = template.Must(template.ParseFiles("./src/templates/eng/wrongregistercustom.html"))
var CustomSolo = template.Must(template.ParseFiles("./src/templates/eng/customsolo.html"))
var CustomMulti = template.Must(template.ParseFiles("./src/templates/eng/custommulti.html"))

func SBPage(w http.ResponseWriter, r *http.Request) {
	scores := readScoresFromFile()
	err := Scoreboard.ExecuteTemplate(w, "scoreboard.html", scores)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func StartMultiplayerPage(w http.ResponseWriter, r *http.Request) {
	multiplayerStart()
	StartGameMultiplayer("./dictionnary/words.txt")
	http.Redirect(w, r, "/multiplayer", http.StatusSeeOther)
}

func StartCustomPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var useMultiplayer bool
		var useItem bool
		life := r.FormValue("life")
		difficulty := r.FormValue("difficulty")
		items := r.FormValue("item")
		multiplayerValue := r.FormValue("multiplayer")
		intLife, err := strconv.Atoi(life)
		if err != nil {
			p := "error"
			err := RegisterCustom.ExecuteTemplate(w, "registercustom.html", p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if difficulty == "" || intLife <= 0 {
			p := "error"
			err := WrongRegisterCustom.ExecuteTemplate(w, "wrongregistercustom.html", p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if difficulty == "easy" {
			difficulty = "./dictionnary/words.txt"
		} else if difficulty == "normal" {
			difficulty = "./dictionnary/words2.txt"
		} else if difficulty == "hard" {
			difficulty = "./dictionnary/words3.txt"
		}
		if items == "on" {
			useItem = true
		}
		if multiplayerValue == "on" {
			useMultiplayer = true
		}
		customStart(intLife, difficulty, useItem, useMultiplayer)
		if useMultiplayer {
			if useItem {
				multiplayer.Items = true
			}
			http.Redirect(w, r, "/custommulti", http.StatusSeeOther)
		} else {
			if useItem {
				hang.Items = true
			}
			http.Redirect(w, r, "/customsolo", http.StatusSeeOther)
		}
	}

	err := RegisterCustom.ExecuteTemplate(w, "registercustom.html", custom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CustomMultiPage(w http.ResponseWriter, r *http.Request) {
	if multiplayer.Turn == 1 {
		rep1 := r.FormValue("answer1")
		findAndReplacePlayer1(rep1)
		testEndGameMulti()
		multiplayer.Turn = 2
	} else if multiplayer.Turn == 2 {
		rep2 := r.FormValue("answer2")
		findAndReplacePlayer2(rep2)
		testEndGameMulti()
		multiplayer.Turn = 1
	}

	err := CustomMulti.ExecuteTemplate(w, "custommulti.html", multiplayer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CustomSoloPage(w http.ResponseWriter, r *http.Request) {
	rep := r.FormValue("answer")
	FindAndReplace(rep)
	testEndGame()

	err := CustomSolo.ExecuteTemplate(w, "customsolo.html", hang)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func MultiplayerPage(w http.ResponseWriter, r *http.Request) {
	if multiplayer.Turn == 1 {
		rep1 := r.FormValue("answer1")
		findAndReplacePlayer1(rep1)
		testEndGameMulti()
		multiplayer.Turn = 2
	} else if multiplayer.Turn == 2 {
		rep2 := r.FormValue("answer2")
		findAndReplacePlayer2(rep2)
		testEndGameMulti()
		multiplayer.Turn = 1
	}

	err := Multiplayer.ExecuteTemplate(w, "multiplayer.html", multiplayer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	hang.Username = ""
	doublePoint = 1
	WinStreak = 0
	p := "Home page"
	err := home.ExecuteTemplate(w, "home.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		difficulty := r.FormValue("difficulty")
		username := r.FormValue("username")

		if difficulty == "" || len(username) == 0 {
			p := "Difficulty not chosen or nickname not specified !"
			err := WrongRegister.ExecuteTemplate(w, "wrongregister.html", p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		} else if difficulty == "easy" {
			HangmanStart(username, difficulty)
			if !isUsernameInSave(username, "./saves/inventory.txt") {
				inventory.Username = username
				inventory.Spell1 = 0
				inventory.Spell2 = 0
				inventory.Spell3 = 0
				hang.Spell1 = inventory.Spell1
				hang.Spell2 = inventory.Spell2
				hang.Spell3 = inventory.Spell3
				writeInventoryToFile(inventory)
			} else {
				inventory = shopBuild(findUsernameLineNumber(username, "./saves/inventory.txt"))
				hang.Spell1 = inventory.Spell1
				hang.Spell2 = inventory.Spell2
				hang.Spell3 = inventory.Spell3
			}
			StartGame("./dictionnary/words.txt")
			http.Redirect(w, r, "/hangman", http.StatusSeeOther)
		} else if difficulty == "normal" {
			HangmanStart(username, difficulty)
			if !isUsernameInSave(username, "./saves/inventory.txt") {
				inventory.Username = username
				inventory.Spell1 = 0
				inventory.Spell2 = 0
				inventory.Spell3 = 0
				hang.Spell1 = inventory.Spell1
				hang.Spell2 = inventory.Spell2
				hang.Spell3 = inventory.Spell3
				writeInventoryToFile(inventory)
			} else {
				inventory = shopBuild(findUsernameLineNumber(username, "./saves/inventory.txt"))
				hang.Spell1 = inventory.Spell1
				hang.Spell2 = inventory.Spell2
				hang.Spell3 = inventory.Spell3
			}
			StartGame("./dictionnary/words2.txt")
			http.Redirect(w, r, "/hangman", http.StatusSeeOther)
		} else if difficulty == "hard" {
			HangmanStart(username, difficulty)
			if !isUsernameInSave(username, "./saves/inventory.txt") {
				inventory.Username = username
				inventory.Spell1 = 0
				inventory.Spell2 = 0
				inventory.Spell3 = 0
				hang.Spell1 = inventory.Spell1
				hang.Spell2 = inventory.Spell2
				hang.Spell3 = inventory.Spell3
				writeInventoryToFile(inventory)
			} else {
				inventory = shopBuild(findUsernameLineNumber(username, "./saves/inventory.txt"))
				hang.Spell1 = inventory.Spell1
				hang.Spell2 = inventory.Spell2
				hang.Spell3 = inventory.Spell3
			}
			StartGame("./dictionnary/words3.txt")
			http.Redirect(w, r, "/hangman", http.StatusSeeOther)
		}
	}

	p := "Page d'inscription"
	err := Register.ExecuteTemplate(w, "register.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HangmanPage(w http.ResponseWriter, r *http.Request) {
	if hang.Username == "" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}
	rep := r.FormValue("answer")
	FindAndReplace(rep)
	testEndGame()
	if hang.GameState == 1 && hang.Username != "" {
		WinStreak += 1
		if isUsernameInSave(hang.Username, "./saves/save.txt") {
			ligne := findUsernameLineNumber(hang.Username, "./saves/save.txt")
			user = playerBuild(ligne)
			if WinStreak != 0 {
				user.Score += ((CheckDifficultyMult(hang.Difficulty) * (len(hang.Word) / 2)) * WinStreak / 2) * doublePoint
			} else {
				user.Score += (CheckDifficultyMult(hang.Difficulty) * (len(hang.Word) / 2)) * doublePoint
			}
			user.Money = user.Money + (CheckDifficultyMult(hang.Difficulty) / 5)
			writeStruct(user, ligne-1)

		} else {
			user = Player{
				Username: hang.Username,
			}
			if WinStreak != 0 {
				user.Score += ((CheckDifficultyMult(hang.Difficulty) * (len(hang.Word) / 2)) * WinStreak / 2) * doublePoint
			} else {
				user.Score += (CheckDifficultyMult(hang.Difficulty) * (len(hang.Word) / 2)) * doublePoint
			}
			user.Money = user.Money + (CheckDifficultyMult(hang.Difficulty) / 5)
			writeStructToFile(user)
		}
	} else if hang.GameState == 2 && hang.Username != "" {
		WinStreak = 0
	}

	err := hangman.ExecuteTemplate(w, "hangman.html", hang)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RetryPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		doublePoint = 1
		hang.GameState = 0
		if hang.Difficulty == "easy" {
			HangmanStart(hang.Username, hang.Difficulty)
			StartGame("./dictionnary/words.txt")
			http.Redirect(w, r, "/hangman", http.StatusSeeOther)
		} else if hang.Difficulty == "normal" {
			HangmanStart(hang.Username, hang.Difficulty)
			StartGame("./dictionnary/words2.txt")
			http.Redirect(w, r, "/hangman", http.StatusSeeOther)
		} else if hang.Difficulty == "hard" {
			HangmanStart(hang.Username, hang.Difficulty)
			StartGame("./dictionnary/words3.txt")
			http.Redirect(w, r, "/hangman", http.StatusSeeOther)
		} else {
			HangmanStart("", "beta")
			StartGame("./dictionnary/beta.txt")
			http.Redirect(w, r, "/beta", http.StatusSeeOther)
		}
	}
}

func UseSpell1(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if inventory.Spell1 > 0 {
			hang.TryLeft = hang.TryLeft + hang.TryLeft
			inventory.Spell1 -= 1
			hang.Spell1 -= 1
			updateShopData(inventory, findUsernameLineNumber(inventory.Username, "./saves/inventory.txt")-1)
		} else {
			fmt.Println("no item")
		}
	}
	http.Redirect(w, r, "/hangman", http.StatusSeeOther)
}

func UseSpell1Multi(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if multiplayer.Turn == 1 {
			multiplayer.TryLeft1 = multiplayer.TryLeft1 + multiplayer.TryLeft1
		} else {
			multiplayer.TryLeft2 = multiplayer.TryLeft2 + multiplayer.TryLeft2
		}
	}
	http.Redirect(w, r, "/custommulti", http.StatusSeeOther)
}

func UseSpell1Custom(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		hang.TryLeft = hang.TryLeft + hang.TryLeft
	}
	http.Redirect(w, r, "/customsolo", http.StatusSeeOther)
}

func UseSpell2(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if inventory.Spell2 > 0 {
			doublePoint = 2
			inventory.Spell2 -= 1
			hang.Spell2 -= 1
			updateShopData(inventory, findUsernameLineNumber(inventory.Username, "./saves/inventory.txt")-1)
		} else {
			fmt.Println("no item")
		}
	}
	http.Redirect(w, r, "/hangman", http.StatusSeeOther)
}

func UseSpell3(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if inventory.Spell3 > 0 {
			hang.WordHidden = revealLetter(hang.Word, hang.WordHidden)
			inventory.Spell3 -= 1
			hang.Spell3 -= 1
			updateShopData(inventory, findUsernameLineNumber(inventory.Username, "./saves/inventory.txt")-1)
		} else {
			fmt.Println("no item")
		}
	}
	http.Redirect(w, r, "/hangman", http.StatusSeeOther)
}

func UseSpell3Multi(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if multiplayer.Turn == 1 {
			multiplayer.WordHidden1 = revealLetter(multiplayer.Word1, multiplayer.WordHidden1)
		} else {
			multiplayer.WordHidden2 = revealLetter(multiplayer.Word2, multiplayer.WordHidden2)
		}
	}
	http.Redirect(w, r, "/custommulti", http.StatusSeeOther)
}

func UseSpell3Custom(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		hang.WordHidden = revealLetter(hang.Word, hang.WordHidden)
	}
	http.Redirect(w, r, "/customsolo", http.StatusSeeOther)
}

func BetaPage(w http.ResponseWriter, r *http.Request) {
	HangmanStart("", "beta")
	StartGame("./dictionnary/beta.txt")
	rep := r.FormValue("answer")
	FindAndReplace(rep)
	testEndGame()

	err := beta.ExecuteTemplate(w, "beta.html", hang)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HangofthedayPage(w http.ResponseWriter, r *http.Request) {
	StartGameHangOfTheDay("")
	if hang.GameState == 1 {
		WinStreak += 1
	} else if hang.GameState == 2 {
		WinStreak = 0
	}

	rep := r.FormValue("answer")
	FindAndReplace(rep)
	testEndGame()

	err := hangoftheday.ExecuteTemplate(w, "hangoftheday.html", hang)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	if Username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	p := "Admin page"
	err := Admin.ExecuteTemplate(w, "admin.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ManagePage(w http.ResponseWriter, r *http.Request) {
	p := "Admin manage page"
	err := Manage.ExecuteTemplate(w, "manage-admin.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ManageSBPage(w http.ResponseWriter, r *http.Request) {
	p := "Admin manage page"
	err := ManageSB.ExecuteTemplate(w, "scoreboard-admin.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ManageShopPage(w http.ResponseWriter, r *http.Request) {
	p := "Admin manage page"
	err := ManageShop.ExecuteTemplate(w, "shop-admin.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AdminHOTDPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		word := r.FormValue("word")
		ChangeWordOfTheDay(word)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
	p := "AdminHOTD page"
	err := AdminHOTD.ExecuteTemplate(w, "hang-of-the-day-admin.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AccountPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		password := r.FormValue("password")
		if password != "" {
			Mdp = password
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}
	}
	err := Account.ExecuteTemplate(w, "account.html", Mdp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ResetAllPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		resetAll("./saves/save.txt")
		resetAll("./saves/inventory.txt")
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func ResetScoreboardPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		resetScoreboard()
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func ResetShopPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		resetShop()
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func CheckSpell1(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if enoughtMoney(user.Money, item1) {
			inventory.Spell1 += 1
			user.Money -= item1
			ligne := findUsernameLineNumber(inventory.Username, "./saves/inventory.txt")
			updateShopData(inventory, ligne-1)
			line := findUsernameLineNumber(user.Username, "./saves/save.txt")
			writeStruct(user, line-1)
		} else {
			fmt.Println("Not enought money")
		}
	}
	http.Redirect(w, r, "/shop", http.StatusSeeOther)
}

func CheckSpell2(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if enoughtMoney(user.Money, item2) {
			inventory.Spell2 += 1
			user.Money -= item2
			ligne := findUsernameLineNumber(inventory.Username, "./saves/inventory.txt")
			updateShopData(inventory, ligne-1)
			line := findUsernameLineNumber(user.Username, "./saves/save.txt")
			writeStruct(user, line-1)
		} else {
			fmt.Println("Not enought money")
		}
	}
	http.Redirect(w, r, "/shop", http.StatusSeeOther)
}

func CheckSpell3(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if enoughtMoney(user.Money, item3) {
			inventory.Spell3 += 1
			user.Money -= item3
			ligne := findUsernameLineNumber(inventory.Username, "./saves/inventory.txt")
			updateShopData(inventory, ligne-1)
			line := findUsernameLineNumber(user.Username, "./saves/save.txt")
			writeStruct(user, line-1)
		} else {
			fmt.Println("Not enought money")
		}
	}
	http.Redirect(w, r, "/shop", http.StatusSeeOther)
}

func ShopPage(w http.ResponseWriter, r *http.Request) {
	err := shop.ExecuteTemplate(w, "shop.html", user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RegisterShopPage(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	if username != "" {
		if isUsernameInSave(username, "./saves/inventory.txt") {
			ligne := findUsernameLineNumber(username, "./saves/inventory.txt")
			inventory = shopBuild(ligne)
			if isUsernameInSave(username, "./saves/save.txt") {
				ligne := findUsernameLineNumber(username, "./saves/save.txt")
				user = playerBuild(ligne)
			} else {
				user.Username = username
				writeStructToFile(user)
			}
			http.Redirect(w, r, "/shop", http.StatusSeeOther)
		} else {
			inventory.Username = username
			writeInventoryToFile(inventory)
			if isUsernameInSave(username, "./saves/save.txt") {
				ligne := findUsernameLineNumber(username, "./saves/save.txt")
				user = playerBuild(ligne)
			} else {
				user.Username = username
				writeStructToFile(user)
			}
			http.Redirect(w, r, "/shop", http.StatusSeeOther)
		}
	}

	err := registerShop.ExecuteTemplate(w, "registershop.html", inventory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		Username = r.FormValue("username")
		password := r.FormValue("password")

		if login(Username, password) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		} else {
			p := "Invalid username or password"
			err := WrongLogin.ExecuteTemplate(w, "wronglogin.html", p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	p := "Login page"
	err := Login.ExecuteTemplate(w, "login.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	p := "Page not found"
	err := Error.ExecuteTemplate(w, "error.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
