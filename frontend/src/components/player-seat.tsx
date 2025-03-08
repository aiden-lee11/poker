import type { Card } from "../types/card"
import { PlayingCard } from "./playing-card"

interface PublicPlayerState {
  playerID: string
  stackSize: number
  active: boolean
}

interface PlayerSeatProps {
  player: PublicPlayerState
  position: {
    top: string
    left: string
    transform: string
  }
  isCurrentPlayer: boolean
  isCurrentTurn: boolean
  holeCards: Card[]
}

export function PlayerSeat({ player, position, isCurrentPlayer, isCurrentTurn, holeCards }: PlayerSeatProps) {
  return (
    <div
      className={`absolute w-32 h-32 flex flex-col items-center justify-center ${isCurrentTurn ? "animate-pulse" : ""}`}
      style={{
        top: position.top,
        left: position.left,
        transform: position.transform,
      }}
    >
      <div
        className={`
        w-24 h-24 rounded-full flex items-center justify-center
        ${player.active ? "bg-gray-800" : "bg-gray-800/50"}
        ${isCurrentPlayer ? "border-4 border-yellow-400" : "border-2 border-gray-600"}
        ${isCurrentTurn ? "ring-4 ring-blue-500" : ""}
      `}
      >
        <div className="text-center">
          <div className="text-xs text-gray-300 mb-1">
            {isCurrentPlayer ? "You" : `Player ${player.playerID.slice(0, 4)}`}
          </div>
          <div className="font-bold">${player.stackSize}</div>
          {!player.active && <div className="text-xs text-red-400 mt-1">Folded</div>}
        </div>
      </div>

      {/* Player's cards */}
      <div className="flex -mt-4 gap-1">
        {isCurrentPlayer && holeCards.length > 0 ? (
          // Show actual cards for current player
          holeCards.map((card, index) => <PlayingCard key={index} card={card} faceUp={true} small />)
        ) : // Show card backs for other players
          player.active ? (
            <>
              <PlayingCard faceUp={false} small />
              <PlayingCard faceUp={false} small />
            </>
          ) : null}
      </div>
    </div>
  )
}

