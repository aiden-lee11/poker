package main

import (
	"encoding/json"
	"fmt"
	"log"
	"poker/eval"
	"poker/table"
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
		ID:              tableID,
		Players:         []*table.Player{},
		PotSize:         0,
		CommunityCards:  []eval.Card{},
		MostRecentRaise: table.Bet{PlayerID: "", BetAmount: 0, Start: true},
		Round:           table.PreFlop,
		// assume table needs two players to be valid
		BigBlindIndex:    0,
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

	fmt.Printf("in advance turn, table.CurrentTurnIndex: %v\n", table.CurrentTurnIndex)
	fmt.Printf("table.MostRecentRaise: %v\n", table.MostRecentRaise)

	// If we have wrapped around to the person that raised initially the round is over
	if !table.MostRecentRaise.Start &&
		table.MostRecentRaise.PlayerID != "" &&
		table.Players[table.CurrentTurnIndex].PlayerID == table.MostRecentRaise.PlayerID {
		gm.advanceBettingRound(table)
		return
	}

	for i := 1; i < numPlayers; i++ {
		nextIndex := (table.CurrentTurnIndex + i) % numPlayers
		fmt.Printf("nextIndex: %v\n", nextIndex)
		if table.Players[nextIndex].PlayingHand {
			table.CurrentTurnIndex = nextIndex
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

	// for now will append the odds of the player in the public state cause why not :)
	for i, p := range table.Players {
		if p == nil {
			continue
		}
		publicPlayers[i] = PublicPlayerState{
			PlayerID:  p.PlayerID,
			StackSize: p.StackSize,
			Active:    p.PlayingHand,
			WinOdds:   p.WinOdds,
			SplitOdds: p.SplitOdds,
		}
	}

	publicState := PublicGameState{
		PotSize:        table.PotSize,
		CommunityCards: table.CommunityCardNames(),
		Players:        publicPlayers,
		CurrentTurn:    publicPlayers[table.CurrentTurnIndex].PlayerID,
	}

	msg := GameMessage{
		Type:    "public_state",
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

	// to fix when a player in the lower slot leaves but the last player doesnt
	// could resolve in a conflict where two players have the same player id
	// possible fixes are to keep track of leavers and store array of useable ids
	// if no useable ids then and only then assign playerID like below
	fmt.Printf("playerTable.Players: %v\n", playerTable.Players)

	n := len(playerTable.Players)
	playerID := "player-" + strconv.Itoa(len(playerTable.Players)+1)
	for i := 0; i < n; i++ {
		if playerTable.Players[i] == nil {
			playerID = "player-" + strconv.Itoa(i)
		}
	}

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
	gm.BroadcastJoin(client)
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
		if p == nil {
			continue
		}
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

	if !playerTable.ValidBet(amount) {
		log.Println("need to bet either 2x most recent raise", player.PlayerID)
		return
	}

	player.StackSize -= amount
	playerTable.PotSize += amount
	playerTable.MostRecentRaise = table.Bet{PlayerID: player.PlayerID, BetAmount: amount, Round: playerTable.Round, Start: false}

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
		if p == nil {
			continue
		}
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

func (gm *GameManager) HandleCheck(client *Client) {
	playerTable, exists := gm.GetTable(client.tableID)

	if !exists {
		log.Printf("player %s is trying to fold in a non existent table %s", client.playerID, client.tableID)
		return
	}

	if !gm.ValidAction(client) {
		log.Println("not the clients turn to act", client.playerID)
		return
	}

	// if someone has bet this round disallow the client from checking
	if playerTable.MostRecentRaise.Round == playerTable.Round && playerTable.MostRecentRaise.BetAmount != 0 {
		log.Println("cannot check when there is a bet this round, most recent bet: ", playerTable.MostRecentRaise)
		return
	}

	// default most recent raise is small blind so if they check disregard that a "raise" was this round
	// start is the flag of being set as default
	if playerTable.MostRecentRaise.PlayerID == client.playerID && playerTable.MostRecentRaise.Start {
		playerTable.MostRecentRaise.Start = false
	}

	log.Printf("Player %s checks", client.playerID)

	gm.AdvanceTurn(playerTable)

	gm.BroadcastState(client.tableID)
}

func (gm *GameManager) HandleFold(client *Client) {
	playerTable, exists := gm.GetTable(client.tableID)

	if !exists {
		log.Printf("player %s is trying to fold in a non existent table %s", client.playerID, client.tableID)
		return
	}

	if !gm.ValidAction(client) {
		log.Println("not the clients turn to act", client.playerID)
		return
	}

	for _, p := range playerTable.Players {
		if p == nil {
			continue
		}
		if p.PlayerID == client.playerID {
			p.PlayingHand = false
			break
		}
	}

	log.Printf("Player %s folds", client.playerID)

	gm.AdvanceTurn(playerTable)
	gm.BroadcastState(client.tableID)
}

func (gm *GameManager) HandleLeave(client *Client) {
	playerTable, exists := gm.GetTable(client.tableID)

	if !exists {
		log.Println("table does not exist", client.tableID)
		return
	}

	// deregister the client from the players
	for i, player := range playerTable.Players {
		if player == nil {
			continue
		}
		if player.PlayerID == client.playerID {
			playerTable.Players[i] = nil
		}
	}

	log.Println("we left successfully")

	gm.BroadcastState(playerTable.ID)
}

func (gm *GameManager) HandleInitGame(client *Client) {
	if client.playerID != "player-1" {
		log.Println("only player one can init the game not ", client.playerID)
		return
	}
	playerTable, exists := gm.GetTable(client.tableID)

	if !exists {
		log.Println("table does not exist", client.tableID)
		return
	}

	playerTable.SetBigBlindIndex()
	playerTable.SetPositions()
	playerTable.Deck = playerTable.ShuffleDeck(eval.UnshuffledDeck)
	playerTable.DistributeCards()

	playerTable.PrintTableDetails()

	gm.BroadcastPrivate(playerTable.ID)
	gm.BroadcastState(playerTable.ID)

}

func (gm *GameManager) BroadcastJoin(client *Client) {
	playerTable, exists := gm.GetTable(client.tableID)

	if !exists {
		log.Println("table not found:", client.tableID)
		return
	}

	msg := GameMessage{
		Type:    "join_response",
		Payload: client.playerID,
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("json marshal error in broadcast state:", err)
		return
	}

	gm.hub.broadcast <- Message{tableID: playerTable.ID, content: jsonMsg, Type: "private", playerID: client.playerID}
}

func (gm *GameManager) BroadcastPrivate(tableID string) {
	playerTable, exists := gm.GetTable(tableID)

	if !exists {
		log.Println("table not found:", tableID)
		return
	}

	for _, player := range playerTable.Players {
		if player == nil {
			continue
		}
		privateState := PrivatePlayerState{
			HoleCards: player.HoleCardNames(),
		}

		msg := GameMessage{
			Type:    "private_state",
			Payload: privateState,
		}

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Println("json marshal error in broadcast state:", err)
			return
		}

		gm.hub.broadcast <- Message{tableID: playerTable.ID, content: jsonMsg, Type: "private", playerID: player.PlayerID}
	}
}

func (gm *GameManager) NewGame(tableID string) {
	playerTable, exists := gm.GetTable(tableID)

	if !exists {
		log.Println("table does not exist", tableID)
		return
	}

	playerTable.SetPositions()
	playerTable.Deck = playerTable.ShuffleDeck(eval.UnshuffledDeck)
	playerTable.DistributeCards()
	playerTable.SetBigBlindIndex()

	playerTable.CommunityCards = []eval.Card{}
	playerTable.PotSize = 0
	playerTable.Round = table.PreFlop

	playerTable.PrintTableDetails()

	gm.BroadcastState(playerTable.ID)
	gm.BroadcastPrivate(playerTable.ID)
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

func (gm *GameManager) advanceBettingRound(t *table.Table) {
	switch t.Round {
	case table.PreFlop:
		t.Round = table.Flop
		t.ShowFlopCards()
		t.SimulateOdds()
		log.Println("Flop dealt")
	case table.Flop:
		t.Round = table.Turn
		t.ShowTurnCard()
		t.SimulateOdds()
		log.Println("Turn dealt")
	case table.Turn:
		t.Round = table.River
		t.ShowRiverCard()
		t.SimulateOdds()
		log.Println("River dealt")
	case table.River:
		log.Println("Betting round complete, ready to evaluate hands")
		hand, winners := t.HandleEvaluateHands()
		log.Println("The winners are.... :", winners)
		gm.BroadcastWinners(hand, winners, t.ID)
		gm.NewGame(t.ID)
		return
	}
	// Set the default MostRecentRaise for the new round:
	t.SetPositions()

	// Optionally, reset CurrentTurnIndex to the first active player for the new round.
	// (This depends on your game rules.)
	gm.BroadcastState(t.ID)
}

func (gm *GameManager) BroadcastWinners(hand, winners []string, tableID string) {
	msg := GameMessage{
		Type: "game_end",
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
