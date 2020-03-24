package game

import (
	"encoding/json"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/pkg/errors"

	"github.com/Southclaws/scribble.rs/src/web"
)

type ReqBodyJoin struct {
	GameID string
	Player string
}

func (g *GameModule) join(w http.ResponseWriter, r *http.Request) {
	var j ReqBodyJoin
	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		web.StatusBadRequest(w, err)
		return
	}

	game, exists := g.Games[j.GameID]
	if !exists {
		web.StatusNotFound(w, errors.New("game not found"))
	}

	for _, p := range game.state.Players {
		if p.Name == j.Player {
			web.StatusUnauthorized(w, errors.New("name already in use"))
			return
		}
	}

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}

	go func() {
		defer conn.Close()

		for {
			select {
			// receive from functions like writeToRest and write packets
			case p := <-game.outgoing[j.Player]:
				msg, err := json.Marshal(p)
				if err != nil {
					web.StatusInternalServerError(w, err)
					return
				}
				err = wsutil.WriteServerMessage(conn, ws.OpText, msg)
				if err != nil {
					web.StatusInternalServerError(w, err)
					return
				}

			// if there's nothing to write, fallthrough and attempt to read
			default:
			}

			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				web.StatusInternalServerError(w, err)
				return
			}
			if op != ws.OpText {
				web.StatusBadRequest(w, errors.New("bad packet op, expected TEXT"))
				return
			}
			var p Packet
			if err := json.Unmarshal(msg, &p); err != nil {
				web.StatusBadRequest(w, errors.Wrap(err, "failed to decode as JSON"))
				return
			}
			game.incoming <- p
		}
	}()
}
