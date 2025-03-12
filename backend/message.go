package main

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
	PlayerID  string  `json:"playerID"`
	StackSize int     `json:"stackSize"`
	Active    bool    `json:"active"`
	WinOdds   float32 `json:"winOdds"`
	SplitOdds float32 `json:"splitOdds"`
}

type PublicGameState struct {
	PotSize        int                 `json:"potSize"`
	CommunityCards []string            `json:"communityCards"`
	Players        []PublicPlayerState `json:"players"`
	CurrentTurn    string              `json:"currentTurn"`
}

type PrivatePlayerState struct {
	HoleCards []string `json:"holeCards"`
}
