"use client"

import { Button } from "./ui/button"
import { Input } from "./ui/input"
import { Label } from "./ui/label"

interface ActionPanelProps {
  betAmount: string
  setBetAmount: (amount: string) => void
  onBet: () => void
  onFold: () => void
  onCheck: () => void
  onCall: () => void
  isActive: boolean
  disabled: boolean
}

export function ActionPanel({
  betAmount,
  setBetAmount,
  onBet,
  onFold,
  onCheck,
  onCall,
  isActive,
  disabled,
}: ActionPanelProps) {
  return (
    <div className={`bg-gray-800 rounded-lg p-4 ${isActive ? "ring-2 ring-blue-500" : ""}`}>
      <div className="flex flex-col sm:flex-row gap-4 items-end">
        <div className="flex-1">
          <Label htmlFor="betAmount">Bet Amount</Label>
          <Input
            id="betAmount"
            type="number"
            value={betAmount}
            onChange={(e) => setBetAmount(e.target.value)}
            placeholder="Enter amount"
            className="bg-gray-700"
            disabled={disabled || !isActive}
          />
        </div>
        <Button
          onClick={onBet}
          disabled={disabled || !isActive || !betAmount}
          className="bg-yellow-600 hover:bg-yellow-700"
        >
          Bet
        </Button>
        <Button onClick={onCall} disabled={disabled || !isActive} className="bg-blue-600 hover:bg-blue-700">
          Call
        </Button>
        <Button onClick={onCheck} disabled={disabled || !isActive} className="bg-green-600 hover:bg-green-700">
          Check
        </Button>
        <Button onClick={onFold} disabled={disabled || !isActive} variant="destructive">
          Fold
        </Button>
      </div>

      {!isActive && (
        <div className="mt-2 text-sm text-gray-400">
          {disabled ? "Join a table to play" : "Waiting for your turn..."}
        </div>
      )}
    </div>
  )
}

