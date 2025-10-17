package models

import (
	"fmt"
	"time"
)

type Difficulty struct {
	Name string
	Rows int
	Cols int
}

var Difficulties = map[string]Difficulty{
	"easy":   {"Facile", 6, 7},
	"normal": {"Normal", 6, 9},
	"hard":   {"Difficile", 7, 8},
}

const (
	EMPTY   = 0
	PLAYER1 = 1
	PLAYER2 = 2
)

type Game struct {
	Board           [][]int
	Rows            int
	Cols            int
	CurrentPlayer   int
	GameState       string
	Winner          int
	Player1Name     string
	Player2Name     string
	Difficulty      string
	GravityInverted bool
	LastGravityFlip time.Time
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
		Board:           board,
		Rows:            difficulty.Rows,
		Cols:            difficulty.Cols,
		CurrentPlayer:   PLAYER1,
		GameState:       "playing",
		Winner:          EMPTY,
		Player1Name:     player1Name,
		Player2Name:     player2Name,
		Difficulty:      difficultyKey,
		GravityInverted: false,
		LastGravityFlip: time.Now(),
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
	// Vérifier l'inversion de gravité avant le coup
	g.CheckGravityInversion()

	if !g.IsValidMove(col) {
		return fmt.Errorf("coup invalide dans la colonne %d", col)
	}

	// Placer la pièce selon la gravité
	if g.GravityInverted {
		// Gravité inversée : placer du haut vers le bas
		for row := 0; row < g.Rows; row++ {
			if g.Board[row][col] == EMPTY {
				g.Board[row][col] = g.CurrentPlayer
				break
			}
		}
	} else {
		// Gravité normale : placer du bas vers le haut
		for row := g.Rows - 1; row >= 0; row-- {
			if g.Board[row][col] == EMPTY {
				g.Board[row][col] = g.CurrentPlayer
				break
			}
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

func (g *Game) CheckGravityInversion() {
	// Inverser la gravité toutes les 30 secondes
	if time.Since(g.LastGravityFlip) >= 30*time.Second {
		g.InvertGravity()
		g.LastGravityFlip = time.Now()
	}
}

func (g *Game) InvertGravity() {
	g.GravityInverted = !g.GravityInverted
	g.ApplyGravity()
}

func (g *Game) ApplyGravity() {
	// Appliquer la physique selon la gravité actuelle
	for col := 0; col < g.Cols; col++ {
		pieces := []int{}
		// Collecter toutes les pièces de la colonne
		for row := 0; row < g.Rows; row++ {
			if g.Board[row][col] != EMPTY {
				pieces = append(pieces, g.Board[row][col])
			}
			g.Board[row][col] = EMPTY
		}

		// Replacer les pièces selon la gravité
		if g.GravityInverted {
			// Gravité inversée : vers le haut
			for i, piece := range pieces {
				g.Board[i][col] = piece
			}
		} else {
			// Gravité normale : vers le bas
			for i, piece := range pieces {
				g.Board[g.Rows-1-i][col] = piece
			}
		}
	}
}
