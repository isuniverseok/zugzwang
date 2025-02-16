package notation

import "zugzwang/game"

var ByteToPiece map[byte]game.Piece = map[byte]game.Piece{
	'p': game.BlackPawn,
	'n': game.BlackKnight,
	'b': game.BlackBishop,
	'r': game.BlackRook,
	'q': game.BlackQueen,
	'k': game.BlackKing,
	'P': game.WhitePawn,
	'N': game.WhiteKnight,
	'B': game.WhiteBishop,
	'R': game.WhiteRook,
	'Q': game.WhiteQueen,
	'K': game.WhiteKing,
}

var PieceToByte map[byte]game.Piece = map[byte]game.Piece{
	game.BlackPawn:   'p',
	game.BlackKnight: 'n',
	game.BlackBishop: 'b',
	game.BlackRook:   'r',
	game.BlackQueen:  'q',
	game.BlackKing:   'k',
	game.WhitePawn:   'P',
	game.WhiteKnight: 'N',
	game.WhiteBishop: 'B',
	game.WhiteRook:   'R',
	game.WhiteQueen:  'Q',
	game.WhiteKing:   'K',
}