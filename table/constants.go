package table

// handString: [Seven of Hearts Five of Spades Three of Clubs Four of Spades Eight of Spades Seven of Spades Queen of Diamonds]
// handInts: [
// 0b00000000001000000010010100001101,
// 0b00000000000010000001001100000111,
// 0b00000000000000101000000100000011,
// 0b00000000000001000001001000000101,
// 0b00000000010000000001011000010001,
// 0b00000000001000000001010100001101,
// 0b00000100000000000100101000011111]

var baseBitDeck = []int32{
	//xxxAKQJT 98765432 CDHS_rrrr xxPPPPPP
	0b00000000_00000001_0001_0000_00000010, // Two of Spades
	0b00000000_00000010_0001_0001_00000011, // Three of Spades
	0b00000000_00000100_0001_0010_00000101, // Four of Spades
	0b00000000_00001000_0001_0011_00000111, // Five of Spades
	0b00000000_00010000_0001_0100_00001011, // Six of Spades
	0b00000000_00100000_0001_0101_00001101, // Seven of Spades
	0b00000000_01000000_0001_0110_00010001, // Eight of Spades
	0b00000000_10000000_0001_0111_00010011, // Nine of Spades
	0b00000001_00000000_0001_1000_00010111, // Ten of Spades
	0b00000010_00000000_0001_1001_00011101, // Jack of Spades
	0b00000100_00000000_0001_1010_00011111, // Queen of Spades
	0b00001000_00000000_0001_1011_00100101, // King of Spades
	0b00010000_00000000_0001_1100_00101001, // Ace of Spades

	0b00000000_00000001_0010_0000_00000010, // Two of Hearts
	0b00000000_00000010_0010_0001_00000011, // Three of Hearts
	0b00000000_00000100_0010_0010_00000101, // Four of Hearts
	0b00000000_00001000_0010_0011_00000111, // Five of Hearts
	0b00000000_00010000_0010_0100_00001011, // Six of Hearts
	0b00000000_00100000_0010_0101_00001101, // Seven of Hearts
	0b00000000_01000000_0010_0110_00010001, // Eight of Hearts
	0b00000000_10000000_0010_0111_00010011, // Nine of Hearts
	0b00000001_00000000_0010_1000_00010111, // Ten of Hearts
	0b00000010_00000000_0010_1001_00011101, // Jack of Hearts
	0b00000100_00000000_0010_1010_00011111, // Queen of Hearts
	0b00001000_00000000_0010_1011_00100101, // King of Hearts
	0b00010000_00000000_0010_1100_00101001, // Ace of Hearts

	0b00000000_00000001_0100_0000_00000010, // Two of Diamonds
	0b00000000_00000010_0100_0001_00000011, // Three of Diamonds
	0b00000000_00000100_0100_0010_00000101, // Four of Diamonds
	0b00000000_00001000_0100_0011_00000111, // Five of Diamonds
	0b00000000_00010000_0100_0100_00001011, // Six of Diamonds
	0b00000000_00100000_0100_0101_00001101, // Seven of Diamonds
	0b00000000_01000000_0100_0110_00010001, // Eight of Diamonds
	0b00000000_10000000_0100_0111_00010011, // Nine of Diamonds
	0b00000001_00000000_0100_1000_00010111, // Ten of Diamonds
	0b00000010_00000000_0100_1001_00011101, // Jack of Diamonds
	0b00000100_00000000_0100_1010_00011111, // Queen of Diamonds
	0b00001000_00000000_0100_1011_00100101, // King of Diamonds
	0b00010000_00000000_0100_1100_00101001, // Ace of Diamonds

	0b00000000_00000001_1000_0000_00000010, // Two of Clubs
	0b00000000_00000010_1000_0001_00000011, // Three of Clubs
	0b00000000_00000100_1000_0010_00000101, // Four of Clubs
	0b00000000_00001000_1000_0011_00000111, // Five of Clubs
	0b00000000_00010000_1000_0100_00001011, // Six of Clubs
	0b00000000_00100000_1000_0101_00001101, // Seven of Clubs
	0b00000000_01000000_1000_0110_00010001, // Eight of Clubs
	0b00000000_10000000_1000_0111_00010011, // Nine of Clubs
	0b00000001_00000000_1000_1000_00010111, // Ten of Clubs
	0b00000010_00000000_1000_1001_00011101, // Jack of Clubs
	0b00000100_00000000_1000_1010_00011111, // Queen of Clubs
	0b00001000_00000000_1000_1011_00100101, // King of Clubs
	0b00010000_00000000_1000_1100_00101001, // Ace of Clubs
}

var BaseStringDeck = []Card{
	"Two of Spades",
	"Three of Spades",
	"Four of Spades",
	"Five of Spades",
	"Six of Spades",
	"Seven of Spades",
	"Eight of Spades",
	"Nine of Spades",
	"Ten of Spades",
	"Jack of Spades",
	"Queen of Spades",
	"King of Spades",
	"Ace of Spades",

	"Two of Hearts",
	"Three of Hearts",
	"Four of Hearts",
	"Five of Hearts",
	"Six of Hearts",
	"Seven of Hearts",
	"Eight of Hearts",
	"Nine of Hearts",
	"Ten of Hearts",
	"Jack of Hearts",
	"Queen of Hearts",
	"King of Hearts",
	"Ace of Hearts",

	"Two of Diamonds",
	"Three of Diamonds",
	"Four of Diamonds",
	"Five of Diamonds",
	"Six of Diamonds",
	"Seven of Diamonds",
	"Eight of Diamonds",
	"Nine of Diamonds",
	"Ten of Diamonds",
	"Jack of Diamonds",
	"Queen of Diamonds",
	"King of Diamonds",
	"Ace of Diamonds",

	"Two of Clubs",
	"Three of Clubs",
	"Four of Clubs",
	"Five of Clubs",
	"Six of Clubs",
	"Seven of Clubs",
	"Eight of Clubs",
	"Nine of Clubs",
	"Ten of Clubs",
	"Jack of Clubs",
	"Queen of Clubs",
	"King of Clubs",
	"Ace of Clubs",
}

var CardToBits = map[Card]int32{
	"Two of Spades":   baseBitDeck[0],
	"Three of Spades": baseBitDeck[1],
	"Four of Spades":  baseBitDeck[2],
	"Five of Spades":  baseBitDeck[3],
	"Six of Spades":   baseBitDeck[4],
	"Seven of Spades": baseBitDeck[5],
	"Eight of Spades": baseBitDeck[6],
	"Nine of Spades":  baseBitDeck[7],
	"Ten of Spades":   baseBitDeck[8],
	"Jack of Spades":  baseBitDeck[9],
	"Queen of Spades": baseBitDeck[10],
	"King of Spades":  baseBitDeck[11],
	"Ace of Spades":   baseBitDeck[12],

	"Two of Hearts":   baseBitDeck[13],
	"Three of Hearts": baseBitDeck[14],
	"Four of Hearts":  baseBitDeck[15],
	"Five of Hearts":  baseBitDeck[16],
	"Six of Hearts":   baseBitDeck[17],
	"Seven of Hearts": baseBitDeck[18],
	"Eight of Hearts": baseBitDeck[19],
	"Nine of Hearts":  baseBitDeck[20],
	"Ten of Hearts":   baseBitDeck[21],
	"Jack of Hearts":  baseBitDeck[22],
	"Queen of Hearts": baseBitDeck[23],
	"King of Hearts":  baseBitDeck[24],
	"Ace of Hearts":   baseBitDeck[25],

	"Two of Diamonds":   baseBitDeck[26],
	"Three of Diamonds": baseBitDeck[27],
	"Four of Diamonds":  baseBitDeck[28],
	"Five of Diamonds":  baseBitDeck[29],
	"Six of Diamonds":   baseBitDeck[30],
	"Seven of Diamonds": baseBitDeck[31],
	"Eight of Diamonds": baseBitDeck[32],
	"Nine of Diamonds":  baseBitDeck[33],
	"Ten of Diamonds":   baseBitDeck[34],
	"Jack of Diamonds":  baseBitDeck[35],
	"Queen of Diamonds": baseBitDeck[36],
	"King of Diamonds":  baseBitDeck[37],
	"Ace of Diamonds":   baseBitDeck[38],

	"Two of Clubs":   baseBitDeck[39],
	"Three of Clubs": baseBitDeck[40],
	"Four of Clubs":  baseBitDeck[41],
	"Five of Clubs":  baseBitDeck[42],
	"Six of Clubs":   baseBitDeck[43],
	"Seven of Clubs": baseBitDeck[44],
	"Eight of Clubs": baseBitDeck[45],
	"Nine of Clubs":  baseBitDeck[46],
	"Ten of Clubs":   baseBitDeck[47],
	"Jack of Clubs":  baseBitDeck[48],
	"Queen of Clubs": baseBitDeck[49],
	"King of Clubs":  baseBitDeck[50],
	"Ace of Clubs":   baseBitDeck[51],
}
