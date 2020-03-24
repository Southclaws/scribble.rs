package game

import "github.com/gorilla/mux"

type GameModule struct {
	Games map[string]Game
}

func New() GameModule {
	return GameModule{
		Games: make(map[string]Game),
	}
}

func (g *GameModule) Routes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/create", g.create).Methods("POST")
	r.HandleFunc("/join", g.join).Methods("POST")

	return r
}
