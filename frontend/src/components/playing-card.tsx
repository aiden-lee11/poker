import type { Card } from "../types/card"

interface PlayingCardProps {
  card?: Card
  faceUp: boolean
  small?: boolean
}

export function PlayingCard({ card, faceUp, small = false }: PlayingCardProps) {
  // Helper function to get the suit symbol
  const getSuitSymbol = (suit?: string) => {
    switch (suit) {
      case "hearts":
        return "♥"
      case "diamonds":
        return "♦"
      case "clubs":
        return "♣"
      case "spades":
        return "♠"
      default:
        return ""
    }
  }

  // Helper function to get the color based on suit
  const getSuitColor = (suit?: string) => {
    if (suit === "hearts" || suit === "diamonds") {
      return "text-red-600"
    }
    return "text-black"
  }

  const cardSize = small ? "w-10 h-14" : "w-16 h-24"

  return (
    <div className={`${cardSize} rounded-md overflow-hidden relative`}>
      {faceUp && card ? (
        <div className="absolute inset-0 bg-white flex flex-col justify-between p-1">
          <div className={`text-left font-bold ${getSuitColor(card.suit)}`}>
            {card.rank}
            <span className="ml-1">{getSuitSymbol(card.suit)}</span>
          </div>
          <div className={`text-center text-2xl ${getSuitColor(card.suit)}`}>{getSuitSymbol(card.suit)}</div>
          <div className={`text-right font-bold ${getSuitColor(card.suit)}`}>
            {card.rank}
            <span className="ml-1">{getSuitSymbol(card.suit)}</span>
          </div>
        </div>
      ) : (
        <div className="absolute inset-0 bg-blue-800 flex items-center justify-center">
          <div className="absolute inset-1 border-2 border-blue-600 rounded-sm"></div>
          <div className="absolute inset-0 bg-[url('/card-back-pattern.png')] opacity-30"></div>
        </div>
      )}
    </div>
  )
}

