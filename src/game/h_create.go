package game

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/Southclaws/scribble.rs/src/web"
)

func (g *GameModule) create(w http.ResponseWriter, r *http.Request) {
	var p Player
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		web.StatusNotAcceptable(w, err)
		return
	}

	u, err := uuid.NewUUID()
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}

	game := Game{
		state: GameState{
			ID: u.String(),
			Canvas: Canvas{
				Pixels: make([]Pixel, 512*512),
			},
			Players: []Player{p},
		},
	}
	g.Games[u.String()] = game

	go game.loop()

	if err := json.NewEncoder(w).Encode(game); err != nil {
		web.StatusInternalServerError(w, err)
		return
	}
}
