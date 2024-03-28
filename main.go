package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func main() {
	dbConn, err := sql.Open("postgres", "")
	if err != nil {
		panic(err)
	}
	err = dbConn.Ping()
	if err != nil {
		panic(err)
	}
	var gameId, numOfLetters int
  rows, err := dbConn.Query("")

}
