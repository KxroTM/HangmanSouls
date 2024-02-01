package hangweb

type CustomSettings struct {
	Life        int
	Difficulty  string
	Items       bool
	Multiplayer bool
}

var custom CustomSettings

func customStart(life int, difficulty string, items bool, multi bool) {
	custom = CustomSettings{
		Life:        life,
		Difficulty:  difficulty,
		Items:       items,
		Multiplayer: multi,
	}

	if custom.Multiplayer {
		multiplayerStart()
		StartGameMultiplayer(difficulty)
		multiplayer.TryLeft1 = custom.Life
		multiplayer.TryLeft2 = custom.Life
	} else {
		HangmanStart("custom", difficulty)
		StartGame(difficulty)
		hang.TryLeft = life
	}
}
