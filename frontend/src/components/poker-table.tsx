import type { Card } from "../types/card"
import { PlayingCard } from "./playing-card"
import { PlayerSeat } from "./player-seat"

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

interface PokerTableProps {
  gameState: PublicGameState
  privateState: PrivatePlayerState
  playerID: string
}

export function PokerTable({ gameState, privateState, playerID }: PokerTableProps) {
  const { potSize, communityCards, players, currentTurn } = gameState

  // Calculate positions for up to 9 players around an oval table
  const getPlayerPositions = (totalPlayers: number) => {
    // Default positions for a 9-player table
    const positions = [
      { top: "50%", left: "50%", transform: "translate(-50%, -50%)" }, // Center (for community cards)
      { top: "85%", left: "50%", transform: "translate(-50%, 0)" }, // Bottom center (player's position)
      { top: "75%", left: "20%", transform: "translate(0, 0)" }, // Bottom left
      { top: "75%", left: "80%", transform: "translate(0, 0)" }, // Bottom right
      { top: "50%", left: "10%", transform: "translate(0, -50%)" }, // Middle left
      { top: "50%", left: "90%", transform: "translate(0, -50%)" }, // Middle right
      { top: "25%", left: "20%", transform: "translate(0, 0)" }, // Top left
      { top: "25%", left: "80%", transform: "translate(0, 0)" }, // Top right
      { top: "15%", left: "50%", transform: "translate(-50%, 0)" }, // Top center
    ]

    // Return only the positions we need based on player count
    return positions.slice(0, totalPlayers + 1) // +1 for the center position
  }

  // Find the current player's index in the players array
  const currentPlayerIndex = players.findIndex((p) => p.playerID === playerID)

  // Reorder players so the current player is at the bottom
  const reorderedPlayers = [...players]
  if (currentPlayerIndex > 0) {
    const currentPlayer = reorderedPlayers.splice(currentPlayerIndex, 1)[0]
    reorderedPlayers.unshift(currentPlayer)
  }

  const positions = getPlayerPositions(players.length)
  const centerPosition = positions[0]

  return (
    <div className="relative w-full aspect-[16/9] bg-green-800 rounded-full mb-6 overflow-hidden border-8 border-brown-800">
      {/* Felt texture overlay */}
      <div className="absolute inset-0 bg-[url('/felt-texture.png')] opacity-30"></div>

      {/* Table border */}
      <div className="absolute inset-4 rounded-full border-4 border-brown-700"></div>

      {/* Pot size */}
      <div
        className="absolute text-white font-bold text-xl bg-black/50 px-4 py-2 rounded-full"
        style={{
          top: `calc(${centerPosition.top} - 60px)`,
          left: centerPosition.left,
          transform: centerPosition.transform,
        }}
      >
        Pot: ${potSize}
      </div>

      {/* Community cards */}
      <div
        className="absolute flex gap-2 justify-center"
        style={{
          top: centerPosition.top,
          left: centerPosition.left,
          transform: centerPosition.transform,
          width: "300px",
        }}
      >
        {communityCards.length > 0 ? (
          communityCards.map((card, index) => <PlayingCard key={index} card={card} faceUp={true} />)
        ) : (
          <div className="text-white/70 text-center">Waiting for cards...</div>
        )}
      </div>

      {/* Player seats */}
      {reorderedPlayers.map((player, index) => {
        const position = positions[index + 1] // +1 because position[0] is for community cards
        const isCurrentPlayer = player.playerID === playerID
        const isCurrentTurn = player.playerID === currentTurn

        // Get the player's hole cards if this is the current player
        const holeCards = isCurrentPlayer ? privateState.holeCards : []

        return (
          <PlayerSeat
            key={player.playerID}
            player={player}
            position={position}
            isCurrentPlayer={isCurrentPlayer}
            isCurrentTurn={isCurrentTurn}
            holeCards={holeCards}
          />
        )
      })}
    </div>
  )
}

