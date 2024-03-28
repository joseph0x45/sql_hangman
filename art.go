package main

var initialState = `
----
|   
|  
|  
|  
|
|
`

var GuessedOneWrong = `
----
|  |
|
|
|
|
|
`

var GuessedTwoWrong = `
----
|  |
|  ()
|
|
|
|
`

var GuessedThreeWrong = `
----
|   |
|   ()
|   ||
|
|
|
`
var GuessedFourWrong = `
----
|   |
|   ()
|  /||
|
|
|
`
var GuessedFiveWrong = `
----
|   |
|   ()
|  /||\
|
|
|
`
var GuessedSixWrong = `
----
|   |
|   ()
|  /||\
|   /
|
|
`
var GuessedSevenWrong = `
----
|   |
|   ()
|  /||\
|   /\
|
|
`

var arts = map[int]string{
	0: initialState,
	1: GuessedOneWrong,
	2: GuessedTwoWrong,
	3: GuessedThreeWrong,
	4: GuessedFourWrong,
	5: GuessedFiveWrong,
	6: GuessedSixWrong,
	7: GuessedSevenWrong,
}

func RenderArt(wrongGuesses int) string {
	return arts[wrongGuesses]
}
