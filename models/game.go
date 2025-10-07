package models

import "fmt"

const (
	ROWS    = 6
	COLS    = 7
	EMPTY   = 0
	PLAYER1 = 1
	PLAYER2 = 2
)

type Game struct {
	Board         [6][7]int
	CurrentPlayer int
	GameState     string
	Winner        int
}

func NewGame() *Game {
	game := &Game{
		Board:         [6][7]int{},
		CurrentPlayer: PLAYER1,
		GameState:     "playing",
		Winner:        EMPTY,
	}
	return game
}

func (g *Game) IsValidMove(col int) bool {
	if col < 0 || col >= COLS {
		return false
	}
	if g.GameState != "playing" {
		return false
	}
	if g.Board[0][col] != EMPTY {
		return false
	}
	return true
}

func (g *Game) IsFull() bool {
	for col := 0; col < COLS; col++ {
		if g.Board[0][col] == EMPTY {
			return false
		}
	}
	g.GameState = "draw"
	return true
}

func (g *Game) MakeMove(col int) error {
	if !g.IsValidMove(col) {
		return fmt.Errorf("coup invalide dans la colonne %d", col)
	}
	for row := ROWS - 1; row >= 0; row-- {
		if g.Board[row][col] == EMPTY {
			g.Board[row][col] = g.CurrentPlayer
			break
		}
	}
	if g.CheckWin() {
		return nil
	}
	if g.IsFull() {
		return nil
	}
	g.CurrentPlayer = 3 - g.CurrentPlayer
	return nil
}

func (g *Game) CheckWin() bool {
	for row := 0; row < ROWS; row++ {
		for col := 0; col < COLS; col++ {
			if g.Board[row][col] != EMPTY {
				player := g.Board[row][col]
				if g.checkDirection(row, col, 0, 1, player) ||
					g.checkDirection(row, col, 1, 0, player) ||
					g.checkDirection(row, col, 1, 1, player) ||
					g.checkDirection(row, col, 1, -1, player) {
					g.GameState = "won"
					g.Winner = player
					return true
				}
			}
		}
	}
	return false
}

func (g *Game) checkDirection(row, col, deltaRow, deltaCol, player int) bool {
	count := 1
	for i := 1; i < 4; i++ {
		newRow := row + i*deltaRow
		newCol := col + i*deltaCol

		if newRow < 0 || newRow >= ROWS || newCol < 0 || newCol >= COLS {
			break
		}
		if g.Board[newRow][newCol] != player {
			break
		}
		count++
	}
	for i := 1; i < 4; i++ {
		newRow := row - i*deltaRow
		newCol := col - i*deltaCol
		if newRow < 0 || newRow >= ROWS || newCol < 0 || newCol >= COLS {
			break
		}
		if g.Board[newRow][newCol] != player {
			break
		}
		count++
	}
	return count >= 4
}
