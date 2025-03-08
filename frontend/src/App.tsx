"use client"

import { useState, useEffect } from "react"
import useWebSocket, { ReadyState } from "react-use-websocket"
import { PokerTable } from "./components/poker-table"
import { ActionPanel } from "./components/action-panel"
import type { Card } from "./types/card"
import { Button } from "./components/ui/button"
import { Input } from "./components/ui/input"
import { Label } from "./components/ui/label"
import { ScrollArea } from "./components/ui/scroll-area"

// Types based on the provided JSON structures
interface PublicPlayerState {
  playerID: string
  stackSize: number
  active: boolean
}

interface PrivatePlayerState {
  holeCards: Card[]
}

interface PublicGameState {
  potSize: number
  communityCards: Card[]
  players: PublicPlayerState[]
  currentTurn: string
}

interface GameEndMessage {
  winners: string[]
  winningHand: string[]
}

function App() {
  // WebSocket URL - make sure this is correct
  const WS_URL = "ws://localhost:8080/ws"

  const [tableID, setTableID] = useState("")
  const [betAmount, setBetAmount] = useState("")
  const [logs, setLogs] = useState<string[]>([])
  const [playerID, setPlayerID] = useState<string>("")
  const [gameState, setGameState] = useState<PublicGameState>({
    potSize: 0,
    communityCards: [],
    players: [],
    currentTurn: "",
  })
  const [privateState, setPrivateState] = useState<PrivatePlayerState>({
    holeCards: [],
  })

  // Set up WebSocket connection using react-use-websocket
  const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(WS_URL, {
    shouldReconnect: () => true,
    reconnectAttempts: 10,
    reconnectInterval: 3000,
    share: false,
    onOpen: () => {
      console.log("Connected to WebSocket server")
      addToLog("Connected to server")
    },
    onClose: () => {
      console.log("Disconnected from WebSocket server")
      addToLog("Disconnected from server")
    },
    onError: (event) => {
      console.error("WebSocket Error:", event)
      addToLog("WebSocket error occurred")
    },
    onMessage: (event) => {
      console.log("Message from server:", event.data)
      addToLog(`Received: ${event.data}`)

      try {
        const data = JSON.parse(event.data)
        console.log(data)

        // Handle different message types
        if (data.type === "join_response") {
          setPlayerID(data.payload)
        } else if (data.type === "public_state") {
          console.log(data.payload)
          setGameState(data.payload)
        } else if (data.type === "private_state") {
          setPrivateState(data.payload)
        } else if (data.type === "game_end") {
          handleGameEnd(data.payload)
        }
      } catch (error) {
        console.error("Error parsing message:", error)
      }
    },
  })

  // Connection status based on readyState
  const connectionStatus = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Connected",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Disconnected",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState]

  const isConnected = readyState === ReadyState.OPEN

  const addToLog = (message: string) => {
    setLogs((prev) => [message, ...prev])
  }

  const handleGameEnd = (payload: GameEndMessage) => {
    const winnersMessage = `Winners: ${payload.winners.join(", ")}`
    const handMessage = `Winning hand: ${payload.winningHand.join(", ")}`
    addToLog(winnersMessage)
    addToLog(handMessage)
  }

  const joinTable = () => {
    if (!tableID) {
      alert("Please enter a Table ID")
      return
    }

    const message = {
      type: "join",
      payload: {
        tableID: tableID,
      },
    }

    sendJsonMessage(message)
    addToLog(`Joining table: ${tableID}`)
  }

  const initGame = () => {
    if (!tableID) {
      alert("You must join a table first!")
      return
    }

    const message = {
      type: "init",
      payload: {
        tableID: tableID,
      },
    }

    sendJsonMessage(message)
    addToLog("Initializing game")
  }

  const placeBet = () => {
    if (!tableID) {
      alert("You must join a table first!")
      return
    }

    if (!betAmount || Number.parseFloat(betAmount) <= 0) {
      alert("Enter a valid bet amount!")
      return
    }

    const message = {
      type: "bet",
      payload: {
        tableID: tableID,
        amount: Number.parseFloat(betAmount),
      },
    }

    sendJsonMessage(message)
    addToLog(`Placing bet: ${betAmount}`)
  }

  const sendAction = (action: "fold" | "check" | "call") => {
    if (!tableID) {
      alert("You must join a table first!")
      return
    }

    const message = {
      type: action,
      payload: {
        tableID: tableID,
      },
    }

    sendJsonMessage(message)
    addToLog(`Action: ${action}`)
  }

  // Process lastJsonMessage when it changes
  useEffect(() => {
    if (lastJsonMessage) {
      console.log("Last JSON message:", lastJsonMessage)
      // Note: We're already handling messages in the onMessage callback,
      // so we don't need to process lastJsonMessage here
    }
  }, [lastJsonMessage])

  return (
    <div className="min-h-screen bg-gray-900 text-white p-4">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-3xl font-bold mb-6">Poker Table</h1>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2">
            <div className="bg-gray-800 rounded-lg p-4 mb-6">
              <div className="flex flex-col sm:flex-row gap-4 mb-4">
                <div className="flex-1">
                  <Label htmlFor="tableID">Table ID</Label>
                  <div className="flex gap-2">
                    <Input
                      id="tableID"
                      value={tableID}
                      onChange={(e) => setTableID(e.target.value)}
                      placeholder="Enter Table ID"
                      className="bg-gray-700"
                    />
                    <Button onClick={joinTable} disabled={!isConnected}>
                      Join Table
                    </Button>
                  </div>
                </div>
                <Button onClick={initGame} disabled={!isConnected || !tableID} className="mt-6 sm:mt-0">
                  Start Game
                </Button>
              </div>

              <div className="mb-4">
                <p className="text-sm mb-1">
                  Connection Status:
                  <span className={isConnected ? "text-green-500" : "text-red-500"}>{` ${connectionStatus}`}</span>
                </p>
                {playerID && <p className="text-sm">Your Player ID: {playerID}</p>}
              </div>
            </div>

            <PokerTable gameState={gameState} privateState={privateState} playerID={playerID} />

            <ActionPanel
              betAmount={betAmount}
              setBetAmount={setBetAmount}
              onBet={placeBet}
              onFold={() => sendAction("fold")}
              onCheck={() => sendAction("check")}
              onCall={() => sendAction("call")}
              isActive={gameState.currentTurn === playerID}
              disabled={!isConnected || !tableID}
            />
          </div>

          <div>
            <div className="bg-gray-800 rounded-lg p-4">
              <h3 className="text-xl font-semibold mb-2">Message Log</h3>
              <ScrollArea className="h-[500px] rounded border border-gray-700 p-2 bg-gray-900">
                {logs.map((log, index) => (
                  <div key={index} className="py-1 border-b border-gray-800 text-sm font-mono">
                    {log}
                  </div>
                ))}
              </ScrollArea>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default App

