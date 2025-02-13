package main

import "poker/table"

type Message struct {
	tableID string
	content []byte
}

type GameMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type PublicPlayerState struct {
	PlayerID  string `json:"playerID"`
	StackSize int    `json:"stackSize"`
	Active    bool   `json:"active"`
}

type PublicGameState struct {
	PotSize        int                 `json:"potSize"`
	CommunityCards []table.Card        `json:"communityCards"`
	Players        []PublicPlayerState `json:"players"`
	CurrentTurn    string              `json:"currentTurn"`
}
