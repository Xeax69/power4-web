package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"power4-web/models"
	"strconv"
)

var funcMap = template.FuncMap{
	"iterate": func(count int) []int {
		var items []int
		for i := 0; i < count; i++ {
			items = append(items, i)
		}
		return items
	},
	"add": func(a, b int) int {
		return a + b
	},
	"isValidMove": func(game *models.Game, col int) bool {
		return game.IsValidMove(col)
	},
}
var templates *template.Template

func init() {
	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	game := models.NewGame()
	data := struct {
		Title         string
		Game          *models.Game
		Board         [6][7]int
		CurrentPlayer int
		GameState     string
		Winner        int
	}{
		Title:         "jeu",
		Game:          game,
		Board:         game.Board,
		CurrentPlayer: game.CurrentPlayer,
		GameState:     game.GameState,
		Winner:        game.Winner,
	}
	err := templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func moveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode interdite", http.StatusMethodNotAllowed)
		return
	}
	colStr := r.FormValue("column")
	col, err := strconv.Atoi(colStr)
	if err != nil {
		http.Error(w, "Colonne invalide", http.StatusBadRequest)
		return
	}
	fmt.Printf("Coup reçu dans la colonne %d\n", col)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("Reset demandé")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", gameHandler)
	http.HandleFunc("/move", moveHandler)
	http.HandleFunc("/reset", resetHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Serveur Puissance4 démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
