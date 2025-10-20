package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"power4-web/models"
	"strconv"
)

type Session struct {
	ID   string
	Game *models.Game
}

var sessions = make(map[string]*Session)

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
	var err error
	templates, err = template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Erreur lors du chargement des templates:", err)
	}
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	setSession(w, session)

	data := struct {
		Title         string
		Game          *models.Game
		Board         [][]int
		CurrentPlayer int
		GameState     string
		Winner        int
		Player1Name   string
		Player2Name   string
		Difficulty    string
		Rows          int
		Cols          int
	}{
		Title:         "jeu",
		Game:          session.Game,
		Board:         session.Game.Board,
		CurrentPlayer: session.Game.CurrentPlayer,
		GameState:     session.Game.GameState,
		Winner:        session.Game.Winner,
		Player1Name:   session.Game.Player1Name,
		Player2Name:   session.Game.Player2Name,
		Difficulty:    session.Game.Difficulty,
		Rows:          session.Game.Rows,
		Cols:          session.Game.Cols,
	}
	err := templates.ExecuteTemplate(w, "game.html", data)
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
	session := getSession(r)
	colStr := r.FormValue("column")
	col, err := strconv.Atoi(colStr)
	if err != nil {
		http.Error(w, "Colonne invalide", http.StatusBadRequest)
		return
	}
	err = session.Game.MakeMove(col)
	if err != nil {
		fmt.Printf("Erreur coup %v\n", err)
	} else {
		fmt.Printf("Coup joué dans la colonne %d\n", col)
	}
	if session.Game.GameState == "won" || session.Game.GameState == "draw" {
		http.Redirect(w, r, "/victory", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
	}
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	session := getSession(r)
	oldGame := session.Game
	session.Game = models.NewGameWithDifficulty(
		oldGame.Difficulty,
		oldGame.Player1Name,
		oldGame.Player2Name,
	)
	fmt.Println("Nouvelle partie créée")
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

func getSession(r *http.Request) *Session {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		sessionID := generateSessionID()
		session := &Session{
			ID:   sessionID,
			Game: models.NewGame(),
		}
		sessions[sessionID] = session
		return session
	}
	sessionID := cookie.Value
	if session, exists := sessions[sessionID]; exists {
		return session
	}
	sessionID = generateSessionID()
	session := &Session{
		ID:   sessionID,
		Game: models.NewGame(),
	}
	sessions[sessionID] = session
	return session
}

func setSession(w http.ResponseWriter, session *Session) {
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: session.ID,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Accueil",
	}
	err := templates.ExecuteTemplate(w, "start.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func startGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	player1 := r.FormValue("player1")
	player2 := r.FormValue("player2")
	difficulty := r.FormValue("difficulty")

	if player1 == "" || player2 == "" {
		http.Error(w, "Noms des joueurs requis", http.StatusBadRequest)
		return
	}
	sessionID := generateSessionID()
	session := &Session{
		ID:   sessionID,
		Game: models.NewGameWithDifficulty(difficulty, player1, player2),
	}
	sessions[sessionID] = session
	setSession(w, session)

	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func victoryHandler(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)

	data := struct {
		Title       string
		GameState   string
		Winner      int
		Player1Name string
		Player2Name string
		Difficulty  string
		Rows        int
		Cols        int
	}{
		Title:       "Résultat",
		GameState:   session.Game.GameState,
		Winner:      session.Game.Winner,
		Player1Name: session.Game.Player1Name,
		Player2Name: session.Game.Player2Name,
		Difficulty:  session.Game.Difficulty,
		Rows:        session.Game.Rows,
		Cols:        session.Game.Cols,
	}

	err := templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func rematchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	session := getSession(r)
	oldGame := session.Game
	session.Game = models.NewGameWithDifficulty(
		oldGame.Difficulty,
		oldGame.Player1Name,
		oldGame.Player2Name,
	)

	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	if _, err := os.Stat("templates"); os.IsNotExist(err) {
		log.Fatal("Le répertoire 'templates' n'existe pas")
	}
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		log.Fatal("Le répertoire 'static' n'existe pas")
	}

	http.HandleFunc("/", startHandler)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/start-game", startGameHandler)
	http.HandleFunc("/move", moveHandler)
	http.HandleFunc("/reset", resetHandler)
	http.HandleFunc("/victory", victoryHandler)
	http.HandleFunc("/rematch", rematchHandler)
	http.HandleFunc("/new-game", newGameHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Serveur Puissance4 démarré sur http://localhost:8080")
	log.Println("Configuration:")
	log.Println("- Templates chargés:", len(templates.Templates()), "fichiers")
	log.Println("- Répertoire statique: ./static")
	log.Println("- Écoute sur toutes les interfaces (0.0.0.0:8080)")

	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal("Erreur lors du démarrage du serveur:", err)
	}
}
