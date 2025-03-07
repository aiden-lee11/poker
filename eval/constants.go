package eval

// handString: [Seven of Hearts Five of Spades Three of Clubs Four of Spades Eight of Spades Seven of Spades Queen of Diamonds]
// handInts: [
// 0b00000000001000000010010100001101,
// 0b00000000000010000001001100000111,
// 0b00000000000000101000000100000011,
// 0b00000000000001000001001000000101,
// 0b00000000010000000001011000010001,
// 0b00000000001000000001010100001101,
// 0b00000100000000000100101000011111]

//xxxAKQJT 98765432 CDHS_rrrr xxPPPPPP
// from: https://github.com/ihendley/treys/blob/70fbaade2f9b63ee3ea41b54aed7be2c921add7e/treys/card.py
// p = prime number of rank (two->2, three->3, four->5, ..., ace->41)
// r = rank of card (two->0, three->1, ..., ace->12)

var TwoOfSpades = Card{Bits: 0b00000000_00000001_0001_0000_00000010, Name: "Two of Spades"}
var ThreeOfSpades = Card{Bits: 0b00000000_00000010_0001_0001_00000011, Name: "Three of Spades"}
var FourOfSpades = Card{Bits: 0b00000000_00000100_0001_0010_00000101, Name: "Four of Spades"}
var FiveOfSpades = Card{Bits: 0b00000000_00001000_0001_0011_00000111, Name: "Five of Spades"}
var SixOfSpades = Card{Bits: 0b00000000_00010000_0001_0100_00001011, Name: "Six of Spades"}
var SevenOfSpades = Card{Bits: 0b00000000_00100000_0001_0101_00001101, Name: "Seven of Spades"}
var EightOfSpades = Card{Bits: 0b00000000_01000000_0001_0110_00010001, Name: "Eight of Spades"}
var NineOfSpades = Card{Bits: 0b00000000_10000000_0001_0111_00010011, Name: "Nine of Spades"}
var TenOfSpades = Card{Bits: 0b00000001_00000000_0001_1000_00010111, Name: "Ten of Spades"}
var JackOfSpades = Card{Bits: 0b00000010_00000000_0001_1001_00011101, Name: "Jack of Spades"}
var QueenOfSpades = Card{Bits: 0b00000100_00000000_0001_1010_00011111, Name: "Queen of Spades"}
var KingOfSpades = Card{Bits: 0b00001000_00000000_0001_1011_00100101, Name: "King of Spades"}
var AceOfSpades = Card{Bits: 0b00010000_00000000_0001_1100_00101001, Name: "Ace of Spades"}

var TwoOfHearts = Card{Bits: 0b00000000_00000001_0010_0000_00000010, Name: "Two of Hearts"}
var ThreeOfHearts = Card{Bits: 0b00000000_00000010_0010_0001_00000011, Name: "Three of Hearts"}
var FourOfHearts = Card{Bits: 0b00000000_00000100_0010_0010_00000101, Name: "Four of Hearts"}
var FiveOfHearts = Card{Bits: 0b00000000_00001000_0010_0011_00000111, Name: "Five of Hearts"}
var SixOfHearts = Card{Bits: 0b00000000_00010000_0010_0100_00001011, Name: "Six of Hearts"}
var SevenOfHearts = Card{Bits: 0b00000000_00100000_0010_0101_00001101, Name: "Seven of Hearts"}
var EightOfHearts = Card{Bits: 0b00000000_01000000_0010_0110_00010001, Name: "Eight of Hearts"}
var NineOfHearts = Card{Bits: 0b00000000_10000000_0010_0111_00010011, Name: "Nine of Hearts"}
var TenOfHearts = Card{Bits: 0b00000001_00000000_0010_1000_00010111, Name: "Ten of Hearts"}
var JackOfHearts = Card{Bits: 0b00000010_00000000_0010_1001_00011101, Name: "Jack of Hearts"}
var QueenOfHearts = Card{Bits: 0b00000100_00000000_0010_1010_00011111, Name: "Queen of Hearts"}
var KingOfHearts = Card{Bits: 0b00001000_00000000_0010_1011_00100101, Name: "King of Hearts"}
var AceOfHearts = Card{Bits: 0b00010000_00000000_0010_1100_00101001, Name: "Ace of Hearts"}

var TwoOfDiamonds = Card{Bits: 0b00000000_00000001_0100_0000_00000010, Name: "Two of Diamonds"}
var ThreeOfDiamonds = Card{Bits: 0b00000000_00000010_0100_0001_00000011, Name: "Three of Diamonds"}
var FourOfDiamonds = Card{Bits: 0b00000000_00000100_0100_0010_00000101, Name: "Four of Diamonds"}
var FiveOfDiamonds = Card{Bits: 0b00000000_00001000_0100_0011_00000111, Name: "Five of Diamonds"}
var SixOfDiamonds = Card{Bits: 0b00000000_00010000_0100_0100_00001011, Name: "Six of Diamonds"}
var SevenOfDiamonds = Card{Bits: 0b00000000_00100000_0100_0101_00001101, Name: "Seven of Diamonds"}
var EightOfDiamonds = Card{Bits: 0b00000000_01000000_0100_0110_00010001, Name: "Eight of Diamonds"}
var NineOfDiamonds = Card{Bits: 0b00000000_10000000_0100_0111_00010011, Name: "Nine of Diamonds"}
var TenOfDiamonds = Card{Bits: 0b00000001_00000000_0100_1000_00010111, Name: "Ten of Diamonds"}
var JackOfDiamonds = Card{Bits: 0b00000010_00000000_0100_1001_00011101, Name: "Jack of Diamonds"}
var QueenOfDiamonds = Card{Bits: 0b00000100_00000000_0100_1010_00011111, Name: "Queen of Diamonds"}
var KingOfDiamonds = Card{Bits: 0b00001000_00000000_0100_1011_00100101, Name: "King of Diamonds"}
var AceOfDiamonds = Card{Bits: 0b00010000_00000000_0100_1100_00101001, Name: "Ace of Diamonds"}

var TwoOfClubs = Card{Bits: 0b00000000_00000001_1000_0000_00000010, Name: "Two of Clubs"}
var ThreeOfClubs = Card{Bits: 0b00000000_00000010_1000_0001_00000011, Name: "Three of Clubs"}
var FourOfClubs = Card{Bits: 0b00000000_00000100_1000_0010_00000101, Name: "Four of Clubs"}
var FiveOfClubs = Card{Bits: 0b00000000_00001000_1000_0011_00000111, Name: "Five of Clubs"}
var SixOfClubs = Card{Bits: 0b00000000_00010000_1000_0100_00001011, Name: "Six of Clubs"}
var SevenOfClubs = Card{Bits: 0b00000000_00100000_1000_0101_00001101, Name: "Seven of Clubs"}
var EightOfClubs = Card{Bits: 0b00000000_01000000_1000_0110_00010001, Name: "Eight of Clubs"}
var NineOfClubs = Card{Bits: 0b00000000_10000000_1000_0111_00010011, Name: "Nine of Clubs"}
var TenOfClubs = Card{Bits: 0b00000001_00000000_1000_1000_00010111, Name: "Ten of Clubs"}
var JackOfClubs = Card{Bits: 0b00000010_00000000_1000_1001_00011101, Name: "Jack of Clubs"}
var QueenOfClubs = Card{Bits: 0b00000100_00000000_1000_1010_00011111, Name: "Queen of Clubs"}
var KingOfClubs = Card{Bits: 0b00001000_00000000_1000_1011_00100101, Name: "King of Clubs"}
var AceOfClubs = Card{Bits: 0b00010000_00000000_1000_1100_00101001, Name: "Ace of Clubs"}

var UnshuffledDeck = []Card{
	TwoOfSpades,
	ThreeOfSpades,
	FourOfSpades,
	FiveOfSpades,
	SixOfSpades,
	SevenOfSpades,
	EightOfSpades,
	NineOfSpades,
	TenOfSpades,
	JackOfSpades,
	QueenOfSpades,
	KingOfSpades,
	AceOfSpades,

	TwoOfHearts,
	ThreeOfHearts,
	FourOfHearts,
	FiveOfHearts,
	SixOfHearts,
	SevenOfHearts,
	EightOfHearts,
	NineOfHearts,
	TenOfHearts,
	JackOfHearts,
	QueenOfHearts,
	KingOfHearts,
	AceOfHearts,

	TwoOfDiamonds,
	ThreeOfDiamonds,
	FourOfDiamonds,
	FiveOfDiamonds,
	SixOfDiamonds,
	SevenOfDiamonds,
	EightOfDiamonds,
	NineOfDiamonds,
	TenOfDiamonds,
	JackOfDiamonds,
	QueenOfDiamonds,
	KingOfDiamonds,
	AceOfDiamonds,

	TwoOfClubs,
	ThreeOfClubs,
	FourOfClubs,
	FiveOfClubs,
	SixOfClubs,
	SevenOfClubs,
	EightOfClubs,
	NineOfClubs,
	TenOfClubs,
	JackOfClubs,
	QueenOfClubs,
	KingOfClubs,
	AceOfClubs,
}
