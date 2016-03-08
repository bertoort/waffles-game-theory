package main

import (
	"encoding/json"
	"log"
)

// these constants represent different game status messages
const waitPaired = "Waiting to get paired"
const gameBegins = "Game begins!"
const over = "Game Over"
const resetWaitPaired = "Our friend has been disconnected... Waiting to get paired again"

//playerInfo is the struct inside gameState
type playerInfo struct {
	Status bool `json:"status"`
	ID     int  `json:"id"`
}

//results is the struct inside gameState
type results struct {
	PlayerOne string `json:"playerOne"`
	PlayerTwo string `json:"playerTwo"`
}

// gameState is the struct which represents the gameState between two players
type gameState struct {
	//renaming json values here to confirm the standard (lowercase var names)
	StatusMessage string     `json:"statusMessage"`
	Messages      []string   `json:"messages"`
	Started       bool       `json:"started"`
	Over          bool       `json:"over"`
	PlayerOne     playerInfo `json:"playerOne"`
	PlayerTwo     playerInfo `json:"playerTwo"`
	Results       results    `json:"results"`
	Reset         bool       `json:"reset"`
	//These are not exported to JSON
	numberOfPlayers int
	playerOneChoice string
	playerTwoChoice string
}

// newGameState is the constructor for the gameState struct and creates the initial gameState Struct (empty board)
func newGameState() gameState {
	pi := playerInfo{
		Status: false,
	}
	gs := gameState{
		StatusMessage: waitPaired,
		Messages:      make([]string, 0),
		Started:       false,
		PlayerOne:     pi,
		PlayerTwo:     pi,
		//These are not exported to JSON
		numberOfPlayers: 0,
	}
	return gs
}

// addPlayer informs the gameState about the new player and alters the statusMessage
func (gs *gameState) addPlayer() {
	gs.numberOfPlayers++
	switch gs.numberOfPlayers {
	case 1:
		gs.StatusMessage = waitPaired
	case 2:
		gs.StatusMessage = gameBegins
		gs.Started = true
	}
}

// recordMove takes in the message and records it
// set id to player
// define type and run funcs accordingly
// record to messages if message
// record move if otherwise
func (gs *gameState) recordMove(msg *incomingMessage) {
	switch msg.Type {
	case "message":
		gs.Reset = false
		gs.recordMessage(msg.Data)
	case "submission":
		gs.Reset = false
		gs.recordSubmission(msg)
	case "command":
		gs.Reset = true
		gs.restartGame()
	}
	gs.checkWinner()
}

//TODO: Refactor
//recordSubmission takes the response and assigns it to the player
func (gs *gameState) recordSubmission(msg *incomingMessage) {
	if gs.PlayerOne.ID == 0 {
		gs.PlayerOne.ID = msg.ID
		gs.playerOneChoice = msg.Data
		gs.PlayerOne.Status = true
	} else if gs.PlayerOne.ID == msg.ID {
		gs.playerOneChoice = msg.Data
		gs.PlayerOne.Status = true
	} else if gs.PlayerTwo.ID == 0 {
		gs.PlayerTwo.ID = msg.ID
		gs.playerTwoChoice = msg.Data
		gs.PlayerTwo.Status = true
	} else if gs.PlayerTwo.ID == msg.ID {
		gs.playerTwoChoice = msg.Data
		gs.PlayerTwo.Status = true
	}
}

func (gs *gameState) checkWinner() {
	if gs.PlayerOne.Status && gs.PlayerTwo.Status {
		gs.Results = results{
			gs.playerOneChoice,
			gs.playerTwoChoice,
		}
		gs.Over = true
		gs.StatusMessage = over
	}
}

func (gs *gameState) recordMessage(text string) {
	gs.Messages = append(gs.Messages, text)
}

// restartGame sets the gameState to a state so that a new game between the same
// players can begin
func (gs *gameState) restartGame() {
	pi := playerInfo{
		Status: false,
	}
	gs.StatusMessage = gameBegins
	gs.Over = false
	gs.PlayerOne = pi
	gs.PlayerTwo = pi
	gs.Results = results{}
	gs.playerOneChoice = ""
	gs.playerTwoChoice = ""
}

// resetGame is needed, when one player drops out. It sets the gameState to a state so that
// the player who is left can wait for a new opponent.
func (gs *gameState) resetGame() {
	gs.restartGame()
	gs.Started = false
	gs.StatusMessage = resetWaitPaired
}

// gameStateToJSON marshals the gameState struct to JSON represented by a slice of bytes
func (gs *gameState) gameStateToJSON() []byte {
	json, err := json.Marshal(gs)
	if err != nil {
		log.Fatal("Error in marshalling json:", err)
	}
	return json
}
