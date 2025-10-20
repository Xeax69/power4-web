package models

import "fmt"

type DifficultyConfig struct {
	Name string
	Rows int
	Cols int
}

var Difficulties = map[string]DifficultyConfig{
	"easy":   {"Facile", 6, 7},
	"normal": {"Normal", 6, 7},
	"hard":   {"Difficile", 8, 9},
}

const (
	EMPTY   = 0
	PLAYER1 = 1
	PLAYER2 = 2
)

type Game struct {
	Board         [][]int
	Rows          int
	Cols          int
	CurrentPlayer int
	GameState     string
	Winner        int
	Player1Name   string
	Player2Name   string
	Difficulty    string
}

func NewGame() *Game {
	return NewGameWithDifficulty("easy", "Joueur 1", "Joueur 2")
}

func NewGameWithDifficulty(difficultyKey, player1Name, player2Name string) *Game {
	difficulty, exists := Difficulties[difficultyKey]
	if !exists {
		difficulty = Difficulties["easy"]
		difficultyKey = "easy"
	}

	board := make([][]int, difficulty.Rows)
	for i := range board {
		board[i] = make([]int, difficulty.Cols)
	}

	game := &Game{
		Board:         board,
		Rows:          difficulty.Rows,
		Cols:          difficulty.Cols,
		CurrentPlayer: PLAYER1,
		GameState:     "playing",
		Winner:        EMPTY,
		Player1Name:   player1Name,
		Player2Name:   player2Name,
		Difficulty:    difficultyKey,
	}
	return game
}

func (g *Game) IsValidMove(col int) bool {
	if col < 0 || col >= g.Cols {
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
	for col := 0; col < g.Cols; col++ {
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

	// Placement normal : du bas vers le haut
	for row := g.Rows - 1; row >= 0; row-- {
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
	for row := 0; row < g.Rows; row++ {
		for col := 0; col < g.Cols; col++ {
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

		if newRow < 0 || newRow >= g.Rows || newCol < 0 || newCol >= g.Cols {
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
		if newRow < 0 || newRow >= g.Rows || newCol < 0 || newCol >= g.Cols {
			break
		}
		if g.Board[newRow][newCol] != player {
			break
		}
		count++
	}
	return count >= 4
}
