package main

import (
	"database/sql"
	"fmt"
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
	WrongGuesses    int `db:"wrong_guesses"`
	GuessedLetters  []string
	Finished        bool   `db:"game_state"`
	WordToGuess     string `db:"word_to_guess"`
	db              *sql.DB
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
	letters := []string{}
	for i := 0; i < game.NumberOfLetters; i++ {
		letters = append(letters, "_")
	}
	gameModel := GameModel{
		ID:              game.ID,
		NumberOfLetters: game.NumberOfLetters,
		WrongGuesses:    0,
		GuessedLetters:  letters,
		Finished:        false,
		WordToGuess:     "",
		db:              dbConn,
	}
	return gameModel
}

func (game GameModel) Init() tea.Cmd {
	return nil
}

func (game GameModel) View() string {
	if game.Finished {
		return fmt.Sprintf("You died!! The word to guess was %s\n", game.WordToGuess)
	}
	ui := ""
	for i := 0; i < game.NumberOfLetters; i++ {
		ui += fmt.Sprintf("%s ", game.GuessedLetters[i])
	}
	strings.TrimSuffix(ui, " ")
	hangman := RenderArt(game.WrongGuesses)
	ui += fmt.Sprintf("\n\n%s", hangman)
	return ui
}

func (game GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "esc":
			return game, tea.Quit
		}
		if isAlpha(key) {
			wrongGuessesCount, gameIsFinished, positions, WordToGuess := 0, false, []uint8{}, ""
			err := game.db.QueryRow("select * from process_guess($1, $2)", key, game.ID).Scan(&wrongGuessesCount, &gameIsFinished, &positions, &WordToGuess)
			if err != nil {
				panic(err)
			}
			game.WrongGuesses = wrongGuessesCount
			game.Finished = gameIsFinished
			game.WordToGuess = WordToGuess
			for _, position := range positions {
				game.GuessedLetters[position-1] = key
			}
			if gameIsFinished {
				return game, tea.Quit
			}
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
