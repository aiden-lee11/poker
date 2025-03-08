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
	BaseStringDeck[0]:  baseBitDeck[0],
	BaseStringDeck[1]:  baseBitDeck[1],
	BaseStringDeck[2]:  baseBitDeck[2],
	BaseStringDeck[3]:  baseBitDeck[3],
	BaseStringDeck[4]:  baseBitDeck[4],
	BaseStringDeck[5]:  baseBitDeck[5],
	BaseStringDeck[6]:  baseBitDeck[6],
	BaseStringDeck[7]:  baseBitDeck[7],
	BaseStringDeck[8]:  baseBitDeck[8],
	BaseStringDeck[9]:  baseBitDeck[9],
	BaseStringDeck[10]: baseBitDeck[10],
	BaseStringDeck[11]: baseBitDeck[11],
	BaseStringDeck[12]: baseBitDeck[12],

	BaseStringDeck[13]: baseBitDeck[13],
	BaseStringDeck[14]: baseBitDeck[14],
	BaseStringDeck[15]: baseBitDeck[15],
	BaseStringDeck[16]: baseBitDeck[16],
	BaseStringDeck[17]: baseBitDeck[17],
	BaseStringDeck[18]: baseBitDeck[18],
	BaseStringDeck[19]: baseBitDeck[19],
	BaseStringDeck[20]: baseBitDeck[20],
	BaseStringDeck[21]: baseBitDeck[21],
	BaseStringDeck[22]: baseBitDeck[22],
	BaseStringDeck[23]: baseBitDeck[23],
	BaseStringDeck[24]: baseBitDeck[24],
	BaseStringDeck[25]: baseBitDeck[25],

	BaseStringDeck[26]: baseBitDeck[26],
	BaseStringDeck[27]: baseBitDeck[27],
	BaseStringDeck[28]: baseBitDeck[28],
	BaseStringDeck[29]: baseBitDeck[29],
	BaseStringDeck[30]: baseBitDeck[30],
	BaseStringDeck[31]: baseBitDeck[31],
	BaseStringDeck[32]: baseBitDeck[32],
	BaseStringDeck[33]: baseBitDeck[33],
	BaseStringDeck[34]: baseBitDeck[34],
	BaseStringDeck[35]: baseBitDeck[35],
	BaseStringDeck[36]: baseBitDeck[36],
	BaseStringDeck[37]: baseBitDeck[37],
	BaseStringDeck[38]: baseBitDeck[38],

	BaseStringDeck[39]: baseBitDeck[39],
	BaseStringDeck[40]: baseBitDeck[40],
	BaseStringDeck[41]: baseBitDeck[41],
	BaseStringDeck[42]: baseBitDeck[42],
	BaseStringDeck[43]: baseBitDeck[43],
	BaseStringDeck[44]: baseBitDeck[44],
	BaseStringDeck[45]: baseBitDeck[45],
	BaseStringDeck[46]: baseBitDeck[46],
	BaseStringDeck[47]: baseBitDeck[47],
	BaseStringDeck[48]: baseBitDeck[48],
	BaseStringDeck[49]: baseBitDeck[49],
	BaseStringDeck[50]: baseBitDeck[50],
	BaseStringDeck[51]: baseBitDeck[51],
}

var BitsToCards = map[int32]Card{
	baseBitDeck[0]:  BaseStringDeck[0],
	baseBitDeck[1]:  BaseStringDeck[1],
	baseBitDeck[2]:  BaseStringDeck[2],
	baseBitDeck[3]:  BaseStringDeck[3],
	baseBitDeck[4]:  BaseStringDeck[4],
	baseBitDeck[5]:  BaseStringDeck[5],
	baseBitDeck[6]:  BaseStringDeck[6],
	baseBitDeck[7]:  BaseStringDeck[7],
	baseBitDeck[8]:  BaseStringDeck[8],
	baseBitDeck[9]:  BaseStringDeck[9],
	baseBitDeck[10]: BaseStringDeck[10],
	baseBitDeck[11]: BaseStringDeck[11],
	baseBitDeck[12]: BaseStringDeck[12],

	baseBitDeck[13]: BaseStringDeck[13],
	baseBitDeck[14]: BaseStringDeck[14],
	baseBitDeck[15]: BaseStringDeck[15],
	baseBitDeck[16]: BaseStringDeck[16],
	baseBitDeck[17]: BaseStringDeck[17],
	baseBitDeck[18]: BaseStringDeck[18],
	baseBitDeck[19]: BaseStringDeck[19],
	baseBitDeck[20]: BaseStringDeck[20],
	baseBitDeck[21]: BaseStringDeck[21],
	baseBitDeck[22]: BaseStringDeck[22],
	baseBitDeck[23]: BaseStringDeck[23],
	baseBitDeck[24]: BaseStringDeck[24],
	baseBitDeck[25]: BaseStringDeck[25],

	baseBitDeck[26]: BaseStringDeck[26],
	baseBitDeck[27]: BaseStringDeck[27],
	baseBitDeck[28]: BaseStringDeck[28],
	baseBitDeck[29]: BaseStringDeck[29],
	baseBitDeck[30]: BaseStringDeck[30],
	baseBitDeck[31]: BaseStringDeck[31],
	baseBitDeck[32]: BaseStringDeck[32],
	baseBitDeck[33]: BaseStringDeck[33],
	baseBitDeck[34]: BaseStringDeck[34],
	baseBitDeck[35]: BaseStringDeck[35],
	baseBitDeck[36]: BaseStringDeck[36],
	baseBitDeck[37]: BaseStringDeck[37],
	baseBitDeck[38]: BaseStringDeck[38],

	baseBitDeck[39]: BaseStringDeck[39],
	baseBitDeck[40]: BaseStringDeck[40],
	baseBitDeck[41]: BaseStringDeck[41],
	baseBitDeck[42]: BaseStringDeck[42],
	baseBitDeck[43]: BaseStringDeck[43],
	baseBitDeck[44]: BaseStringDeck[44],
	baseBitDeck[45]: BaseStringDeck[45],
	baseBitDeck[46]: BaseStringDeck[46],
	baseBitDeck[47]: BaseStringDeck[47],
	baseBitDeck[48]: BaseStringDeck[48],
	baseBitDeck[49]: BaseStringDeck[49],
	baseBitDeck[50]: BaseStringDeck[50],
	baseBitDeck[51]: BaseStringDeck[51],
}
