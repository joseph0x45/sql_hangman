package main

import (
	"database/sql"
	"github.com/blockloop/scan/v2"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/lib/pq"
	"strings"
)

func isAlpha(key string) bool {
	letters := "abcdefghijklmnopqrstuvwxyz"
	return strings.Contains(letters, key)
}

type GameData struct {
	ID              string `db:"game_id"`
	NumberOfLetters int    `db:"number_of_letters"`
}

type GameModel struct {
	ID              string
	NumberOfLetters int
	GuessesCount    int
	GuessedLetters  []string
	Finished        bool
}

func initModel(dbConn *sql.DB) GameModel {
	game := new(GameData)
	rows, err := dbConn.Query("select game_id, number_of_letters from start_game()")
	if err != nil {
		panic(err)
	}
	err = scan.Row(game, rows)
	if err != nil {
		panic(err)
	}
	gameModel := GameModel{
		ID:              game.ID,
		NumberOfLetters: game.NumberOfLetters,
		GuessesCount:    0,
		GuessedLetters:  []string{},
		Finished:        false,
	}
	return gameModel
}

func (game GameModel) Init() tea.Cmd {
	return nil
}

func (game GameModel) View() string {
	wordToGuess := ""
	for i := 0; i < game.NumberOfLetters; i++ {
		wordToGuess += "_ "
	}
	strings.TrimSuffix(wordToGuess, " ")
	return wordToGuess
}

func (game GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return game, tea.Quit
		}
		if !isAlpha(msg.String()) {
			return nil, nil
		}
	}
	return game, nil
}

func main() {
	dbConn, err := sql.Open("postgres", "postgres://postgres:pwd@localhost:5432/game?sslmode=disable")
	if err != nil {
		panic(err)
	}
	err = dbConn.Ping()
	if err != nil {
		panic(err)
	}
	initalModel := initModel(dbConn)
	program := tea.NewProgram(initalModel)
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
