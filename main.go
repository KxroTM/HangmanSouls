package main

import (
	hangweb "hangweb/go"
	"net/http"
	"os/exec"
	"runtime"
)

func main() {
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./src/style"))))
	http.HandleFunc("/", hangweb.NotFoundHandler)
	http.HandleFunc("/home", hangweb.HomePage)
	http.HandleFunc("/register", hangweb.RegisterPage)
	http.HandleFunc("/hangman", hangweb.HangmanPage)
	http.HandleFunc("/retry", hangweb.RetryPage)
	http.HandleFunc("/use-spell1", hangweb.UseSpell1)
	http.HandleFunc("/use-spell2", hangweb.UseSpell2)
	http.HandleFunc("/use-spell3", hangweb.UseSpell3)
	http.HandleFunc("/beta", hangweb.BetaPage)
	http.HandleFunc("/hangoftheday", hangweb.HangofthedayPage)
	http.HandleFunc("/login", hangweb.LoginPage)
	http.HandleFunc("/admin", hangweb.AdminPage)
	http.HandleFunc("/account", hangweb.AccountPage)
	http.HandleFunc("/adminhodt", hangweb.AdminHOTDPage)
	http.HandleFunc("/manage", hangweb.ManagePage)
	http.HandleFunc("/managescoreboard", hangweb.ManageSBPage)
	http.HandleFunc("/manageshop", hangweb.ManageShopPage)
	http.HandleFunc("/reset-all", hangweb.ResetAllPage)
	http.HandleFunc("/reset-shop", hangweb.ResetShopPage)
	http.HandleFunc("/reset-scoreboard", hangweb.ResetScoreboardPage)
	http.HandleFunc("/scoreboard", hangweb.SBPage)
	http.HandleFunc("/shop", hangweb.ShopPage)
	http.HandleFunc("/registershop", hangweb.RegisterShopPage)
	http.HandleFunc("/buy-spell1", hangweb.CheckSpell1)
	http.HandleFunc("/buy-spell2", hangweb.CheckSpell2)
	http.HandleFunc("/buy-spell3", hangweb.CheckSpell3)
	http.HandleFunc("/multiplayer", hangweb.MultiplayerPage)
	http.HandleFunc("/startmulti", hangweb.StartMultiplayerPage)
	http.HandleFunc("/registercustom", hangweb.StartCustomPage)
	http.HandleFunc("/customsolo", hangweb.CustomSoloPage)
	http.HandleFunc("/custommulti", hangweb.CustomMultiPage)
	http.HandleFunc("/use-spell1multi", hangweb.UseSpell1Multi)
	http.HandleFunc("/use-spell3multi", hangweb.UseSpell3Multi)
	http.HandleFunc("/use-spell1custom", hangweb.UseSpell1Custom)
	http.HandleFunc("/use-spell3custom", hangweb.UseSpell3Custom)
	linkToOpen := "http://localhost:8080/home"
	openLink(linkToOpen)
	http.ListenAndServe(":8080", nil)
}

func openLink(link string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", link)
	case "linux":
		cmd = exec.Command("xdg-open", link)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", link)
	default:
		panic("Système d'exploitation non pris en charge")
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
