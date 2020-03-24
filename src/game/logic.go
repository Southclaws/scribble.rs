package game

type Player struct {
	Name  string
	Score int
}

type GameState struct {
	ID      string
	Canvas  Canvas
	Players []Player
}

type PacketType uint8

const (
	PacketTypeDraw PacketType = 0x0
	PacketTypeChat PacketType = 0x1
)

type Packet struct {
	T PacketType
	P string
	D interface{}
}

type Game struct {
	state    GameState
	incoming chan Packet            // packets received from players
	outgoing map[string]chan Packet // packets to send to players
}

func (g *Game) loop() {
	for {
		select {
		case p := <-g.incoming:
			g.handleIncoming(p)
		}
	}
}

func (g *Game) handleIncoming(p Packet) {
	switch p.T {
	case PacketTypeDraw:
		px := p.D.(Draws)
		g.state.Canvas.draw(px)
		g.writeForRest(p, px)

	case PacketTypeChat:
	}
}

func (g *Game) writeForRest(p Packet, px Draws) {
	for _, player := range g.state.Players {
		if player.Name == p.P {
			continue
		}
		g.outgoing[player.Name] <- p
	}
}
