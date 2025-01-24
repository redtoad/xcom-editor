package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/redtoad/xcom-editor/savegame"
)

func (e *Editor) loadGames() ([]*savegame.Savegame, error) {
	var err error
	games := make([]*savegame.Savegame, 10)
	for i := 1; i < 10; i++ {
		path := filepath.Join(e.config.GamePath, fmt.Sprintf("GAME_%d", i))
		games[i-1], err = savegame.Load(path)
		if err != nil {
			games[i-1] = nil
		}
	}
	return games, nil
}

func (e *Editor) handleGames(rootPath string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			games, err := e.loadGames()
			if err != nil {
				http.Error(w, "could not load games", http.StatusInternalServerError)
				return
			}
			component := GamesList(games)
			component.Render(r.Context(), w)
		},
	)
}
