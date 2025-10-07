package models

type Game struct {
	Board         [6][7]int
	CurrentPlayer int
	GameState     string
	Winner        int
}

const (
	ROWS    = 6
	COLS    = 7
	EMPTY   = 0
	PLAYER1 = 1
	PLAYER2 = 2
)

func NewGame() *Game {
	game := &Game{
		Board:         [6][7]int{},
		CurrentPlayer: PLAYER1,
		GameState:     "play",
		Winner:        EMPTY,
	}
	return game
}
