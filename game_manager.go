package main

import (
	"encoding/json"
	"fmt"
	"log"
	"poker/eval"
	"poker/table"
	"reflect"
	"strconv"
	"sync"
)

type GameManager struct {
	tables map[string]*table.Table
	hub    *Hub
	mu     sync.Mutex
}

func NewGameManager(hub *Hub) *GameManager {
	return &GameManager{
		tables: make(map[string]*table.Table),
		hub:    hub,
	}
}

func (gm *GameManager) CreateTable(tableID string) *table.Table {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	if _, exists := gm.tables[tableID]; exists {
		return gm.tables[tableID]
	}

	newTable := &table.Table{
		ID:               tableID,
		Players:          []*table.Player{},
		PotSize:          0,
		CommunityCards:   []table.Card{},
		MostRecentRaise:  table.Bet{PlayerID: "", BetAmount: 0},
		Round:            table.PreFlop,
		CurrentTurnIndex: 0,
	}

	gm.tables[tableID] = newTable
	return newTable
}

func (gm *GameManager) GetTable(tableID string) (*table.Table, bool) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	table, exists := gm.tables[tableID]
	return table, exists
}

func (gm *GameManager) AdvanceTurn(table *table.Table) {
	numPlayers := len(table.Players)
	if numPlayers < 2 {
		log.Println("Need at least two players to play... loner")
		return
	}

	for i := 1; i < numPlayers; i++ {
		nextIndex := (table.CurrentTurnIndex + i) % numPlayers
		fmt.Printf("nextIndex: %v\n", nextIndex)
		if table.Players[nextIndex].PlayingHand {
			table.CurrentTurnIndex = nextIndex
			if table.MostRecentRaise.PlayerID != "" && table.Players[nextIndex].PlayerID == table.MostRecentRaise.PlayerID {
				gm.advanceBettingRound(table)
			}
			return
		}
	}

	log.Println("No active players left, ending hand.")
}

func (gm *GameManager) BroadcastState(tableID string) {
	table, exists := gm.GetTable(tableID)

	if !exists {
		log.Println("table not found:", tableID)
		return
	}

	// Remake the array to ensure no data mismatches
	publicPlayers := make([]PublicPlayerState, len(table.Players))

	for i, p := range table.Players {
		publicPlayers[i] = PublicPlayerState{
			PlayerID:  p.PlayerID,
			StackSize: p.StackSize,
			Active:    p.PlayingHand,
		}
	}

	publicState := PublicGameState{
		PotSize:        table.PotSize,
		CommunityCards: table.CommunityCards,
		Players:        publicPlayers,
		CurrentTurn:    publicPlayers[table.CurrentTurnIndex].PlayerID,
	}

	msg := GameMessage{
		Type:    "stateUpdate",
		Payload: publicState,
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("json marshal error in broadcast state:", err)
		return
	}

	gm.hub.broadcast <- Message{tableID: tableID, Type: "public", content: jsonMsg}
}

func (gm *GameManager) HandleJoin(client *Client, payload interface{}) {
	params, ok := payload.(map[string]interface{})

	if !ok {
		log.Println("invalid payload in handle join")
		return
	}

	tableID, ok := params["tableID"].(string)
	if !ok {
		log.Println("missing tableID in handle join")
		return
	}

	// create or grab the table
	playerTable := gm.CreateTable(tableID)

	playerID := "player-" + strconv.Itoa(len(playerTable.Players)+1)

	stackSizeFloat, ok := params["stackSize"].(float64)

	// No stacksize in params so set to default
	if !ok {
		stackSizeFloat = 1000.0
	}

	stackSize := int(stackSizeFloat)

	player := &table.Player{
		PlayerID:    playerID,
		StackSize:   stackSize,
		PlayingHand: true,
	}

	playerTable.Players = append(playerTable.Players, player)

	client.tableID = tableID
	client.playerID = playerID

	log.Printf("Player %s joined table %s", playerID, tableID)

	gm.BroadcastState(tableID)
}

// func for checking if player is the action player
func (gm *GameManager) ValidAction(client *Client) bool {
	table, exists := gm.GetTable(client.tableID)

	if !exists {
		log.Println("table not found:", client.tableID)
		return false
	}

	return table.Players[table.CurrentTurnIndex].PlayerID == client.playerID

}

func (gm *GameManager) HandleBet(client *Client, payload interface{}) {
	playerTable, exists := gm.GetTable(client.tableID)

	if !exists {
		log.Println("table not found for client in handle bet:", client.tableID)
		return
	}

	if !gm.ValidAction(client) {
		log.Println("not the clients turn to act", client.playerID)
		return
	}

	params, ok := payload.(map[string]interface{})
	if !ok {
		log.Println("invalid payload in handle bet")
		return
	}

	amountFloat, ok := params["amount"].(float64)
	if !ok {
		log.Println("invalid or missing bet amount in handle bet")
		return
	}

	amount := int(amountFloat)
	var player *table.Player
	for _, p := range playerTable.Players {
		if p.PlayerID == client.playerID {
			player = p
			break
		}
	}

	if player == nil {
		log.Println("player not found in handle bet based on client id")
		return
	}

	if player.StackSize < amount {
		log.Println("insufficient stack in handle bet for bet size for player", player.PlayerID)
		return
	}

	player.StackSize -= amount
	playerTable.PotSize += amount
	playerTable.MostRecentRaise = table.Bet{PlayerID: player.PlayerID, BetAmount: amount, Round: playerTable.Round}

	log.Printf("Player %s bet %d at table %s", player.PlayerID, amount, playerTable.ID)

	gm.AdvanceTurn(playerTable)

	gm.BroadcastState(playerTable.ID)

}

func (gm *GameManager) HandleCall(client *Client) {
	playerTable, exists := gm.GetTable(client.tableID)

	if !exists {
		log.Println("table not found for client in handle bet:", client.tableID)
		return
	}

	if !gm.ValidAction(client) {
		log.Println("not the clients turn to act", client.playerID)
		return
	}

	var player *table.Player
	for _, p := range playerTable.Players {
		if p.PlayerID == client.playerID {
			player = p
			break
		}
	}

	if player == nil {
		log.Println("player not found in handle bet based on client id")
		return
	}

	amount := playerTable.MostRecentRaise.BetAmount

	if player.StackSize < amount {
		log.Println("insufficient stack in handle bet for bet size for player", player.PlayerID)
		return
	}

	player.StackSize -= amount
	playerTable.PotSize += amount

	log.Printf("Player %s bet %d at table %s", player.PlayerID, amount, playerTable.ID)

	gm.AdvanceTurn(playerTable)

	gm.BroadcastState(playerTable.ID)

}

// Need to figure out turn logic here but for now doesn't do anything
func (gm *GameManager) HandleCheck(client *Client) {
	playerTable, exists := gm.GetTable(client.tableID)

	if !exists {
		log.Println("player %s is trying to fold in a non existent table %s", client.playerID, client.tableID)
		return
	}

	if !gm.ValidAction(client) {
		log.Println("not the clients turn to act", client.playerID)
		return
	}

	if playerTable.MostRecentRaise.Round == playerTable.Round && playerTable.MostRecentRaise.BetAmount != 0 {
		log.Println("cannot check when there is a bet this round, most recent bet: ", playerTable.MostRecentRaise)
		return
	}

	log.Printf("Player %s checks", client.playerID)

	gm.AdvanceTurn(playerTable)

	gm.BroadcastState(client.tableID)
}

func (gm *GameManager) HandleFold(client *Client) {
	playerTable, exists := gm.GetTable(client.tableID)

	if !exists {
		log.Println("player %s is trying to fold in a non existent table %s", client.playerID, client.tableID)
		return
	}

	if !gm.ValidAction(client) {
		log.Println("not the clients turn to act", client.playerID)
		return
	}

	for _, p := range playerTable.Players {
		if p.PlayerID == client.playerID {
			p.PlayingHand = false
			break
		}
	}

	log.Printf("Player %s folds", client.playerID)

	gm.AdvanceTurn(playerTable)
	gm.BroadcastState(client.tableID)
}

func (gm *GameManager) HandleInitGame(client *Client) {
	if client.playerID != "player-1" {
		log.Println("only player one can init the game not ", client.playerID)
		return
	}
	table, exists := gm.GetTable(client.tableID)

	if !exists {
		log.Println("table does not exist", client.tableID)
		return
	}

	table.SetDefaultMostRecentRaise()
	table.ShuffleDeck()
	table.DistributeCards()
	fmt.Printf("table.Deck: %v\n", table.Deck)

	for _, player := range table.Players {
		privateState := PrivatePlayerState{
			HoleCards: player.HoleCards[:],
		}

		msg := GameMessage{
			Type:    "stateUpdate",
			Payload: privateState,
		}

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Println("json marshal error in broadcast state:", err)
			return
		}

		gm.hub.broadcast <- Message{tableID: table.ID, content: jsonMsg, Type: "private", playerID: player.PlayerID}
	}
}

// These two functions might become intertwined because if were doing the simulations
// then each new card should trigger

// honestly this should not be a client send for now ill do it bc easy
// but should be something server side where the game starts and then checks
// if a loop has been done after flop then turn then river
func (gm *GameManager) HandleDealCards(table *table.Table) {
	n := len(table.CommunityCards)
	switch n {
	case 0:
		table.ShowFlopCards()
	case 3:
		table.ShowTurnCard()
	case 4:
		table.ShowRiverCard()
	}

	gm.BroadcastState(table.ID)
}

// similar to above this should be handled by server eventually
// would be triggered when a table has done the full round after river
// or when all active players cannot bet further IN TERMS OF DETERMINING WINNER
// returns stringified hand and list of winners

func (gm *GameManager) HandleEvaluateHands(t *table.Table) ([]string, []string) {
	winners := []string{t.Players[0].PlayerID}
	resHand, highest := t.Players[0].EvalHand(t)

	for i := 1; i < len(t.Players); i++ {
		player := t.Players[i]

		if !player.PlayingHand {
			continue
		}

		hand, val := player.EvalHand(t)
		// val here means like two pair one pair etc

		if val > highest {
			log.Printf("New highest val of %d with hand %v", val, hand)
			highest = val
			resHand = hand
			winners = []string{player.PlayerID}
		} else if val == highest {
			log.Printf("Same val of %d with hand %v", val, hand)
			best := eval.HandleTie(hand, resHand)
			if !reflect.DeepEqual(resHand, best) {
				winners = []string{player.PlayerID}
				resHand = best
			} else {
				winners = append(winners, player.PlayerID)
			}
		}
	}

	return table.StringifyHand(resHand), winners
}

func (gm *GameManager) advanceBettingRound(t *table.Table) {
	switch t.Round {
	case table.PreFlop:
		t.Round = table.Flop
		t.ShowFlopCards()
		log.Println("Flop dealt")
	case table.Flop:
		t.Round = table.Turn
		t.ShowTurnCard()
		log.Println("Turn dealt")
	case table.Turn:
		t.Round = table.River
		t.ShowRiverCard()
		log.Println("River dealt")
	case table.River:
		log.Println("Betting round complete, ready to evaluate hands")
		hand, winners := gm.HandleEvaluateHands(t)
		log.Println("The winners are.... :", winners)
		gm.BroadcastWinners(hand, winners, t.ID)
		return
	}
	// Set the default MostRecentRaise for the new round:
	t.SetDefaultMostRecentRaise()

	// Optionally, reset CurrentTurnIndex to the first active player for the new round.
	// (This depends on your game rules.)
	gm.BroadcastState(t.ID)
}

func (gm *GameManager) BroadcastWinners(hand, winners []string, tableID string) {
	msg := GameMessage{
		Type: "stateUpdate",
		Payload: GameEndMessage{
			Winners:     winners,
			WinningHand: hand,
		},
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("json marshal error in broadcast state:", err)
		return
	}

	gm.hub.broadcast <- Message{tableID: tableID, Type: "public", content: jsonMsg}
}
