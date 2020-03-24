# scribble.rs API

A quick API for a drawing game. Will be deleted probably.

Demonstrates an asynchronous game-loop style websocket thingymabob.

`func (g *GameModule) join` handles players joining, which spawns a goroutine that handles packet writes and reads.

`func (g *GameModule) create` handles creation of game lobbies, which creates an in-memory canvas and one player (the one that created the game.)

`logic.go` contains the basic game loop.
