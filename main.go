package main

import (
	"database/sql"
	"fmt"
	"github.com/blockloop/scan/v2"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
)

func isAlpha(key string) bool {
	letters := "abcdefghijklmnopqrstuvwxyz"
	return strings.Contains(letters, key)
}

func parseIntArray(intArrayStr string) []int {
	var result []int
	if intArrayStr == "{}" {
		return result
	}
	intArrayStr = strings.TrimPrefix(intArrayStr, "{")
	intArrayStr = strings.TrimSuffix(intArrayStr, "}")
	parts := strings.Split(intArrayStr, ",")
	for _, part := range parts {
		intValue, err := strconv.Atoi(part)
		if err != nil {
			panic(err)
		}
		result = append(result, intValue)
	}
	return result
}

func parseStrArray(strArray string) []string {
	var result []string
	if strArray == "{}" {
		return result
	}
	strArray = strings.TrimPrefix(strArray, "{")
	strArray = strings.TrimSuffix(strArray, "}")
	parts := strings.Split(strArray, ",")
	for _, part := range parts {
		result = append(result, part)
	}
	return result
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
	AlreadyGuessed  []string
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
		AlreadyGuessed:  []string{},
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
	ui := ""
	for i := 0; i < game.NumberOfLetters; i++ {
		ui += fmt.Sprintf("%s ", game.GuessedLetters[i])
	}
	strings.TrimSuffix(ui, " ")
	ui += "\t"
	ui += "Guess history: "
	wrongGuesses := ""
	for _, guess := range game.AlreadyGuessed {
		wrongGuesses += fmt.Sprintf("%s ", guess)
	}
	ui += wrongGuesses
	hangman := RenderArt(game.WrongGuesses)
	ui += fmt.Sprintf("\n\n%s", hangman)
	if game.Finished {
		if game.WordToGuess != "" {
			ui += fmt.Sprintf("\nYou failed to guess: %s\n\n", game.WordToGuess)
		} else {
			ui += fmt.Sprintf("\nYay you win :)\n\n")
		}
	}
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
			game.AlreadyGuessed = append(game.AlreadyGuessed, key)
			wrongGuessesCount, gameIsFinished, positionsStr, WordToGuess, alreadyGuessed := 0, false, "", "", ""
			err := game.db.QueryRow("select * from process_guess($1, $2)", key, game.ID).Scan(&wrongGuessesCount, &gameIsFinished, &positionsStr, &WordToGuess, &alreadyGuessed)
			if err != nil {
				panic(err)
			}
			game.WrongGuesses = wrongGuessesCount
			game.Finished = gameIsFinished
			game.WordToGuess = WordToGuess
			game.AlreadyGuessed = parseStrArray(alreadyGuessed)
			positions := parseIntArray(positionsStr)
			for _, position := range positions {
				game.GuessedLetters[position-1] = key
			}
			if game.Finished {
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
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	initalModel := initModel(dbConn)
	program := tea.NewProgram(initalModel)
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
