package main

import "poker/table"

type Message struct {
	tableID  string
	playerID string
	Type     string
	content  []byte
}

type GameMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type GameEndMessage struct {
	Winners     []string `json:"winners"`
	WinningHand []string `json:"winningHand"`
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

type PrivatePlayerState struct {
	HoleCards []table.Card `json:"holeCards"`
}
