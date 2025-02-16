package notation

import (
	"strconv"
	"strings"
	"zugzwang/game"
)

func ParseFenString(fen string) *game.State {
	state := &game.State{}
	args := strings.Split(fen, " ")
	ranks := strings.Split(args[0], "/")

	for i, rank := range ranks {
		currFile := 0
		for _, c := range rank {
			if byte(c) >= '1' && byte(c) <= '8' {
				currFile += int(c-'1') + 1
			} else {
				state.Board[i][currFile] = ByteToPiece[byte(c)]
				currFile += 1
			}
		}
	}

	state.SideToMove = ByteToSide[args[1][0]]

	state.CastleRights = 0

	if strings.Contains(args[2], "K") {
		state.CastleRights |= game.WhiteKingSide
	}
	if strings.Contains(args[2], "Q") {
		state.CastleRights |= game.WhiteQueenSide
	}
	if strings.Contains(args[2], "k") {
		state.CastleRights |= game.BlackKingSide
	}
	if strings.Contains(args[2], "q") {
		state.CastleRights |= game.BlackQueenSide
	}

	state.EnPassantPos = nil
	if args[3] != "-" {
		state.EnPassantPos = ParsePosString(args[3])
		if state.SideToMove == game.White {
			state.EnPassantPos.X += 1
		} else {
			state.EnPassantPos.X -= 1
		}
	}

	state.FiftyCount, _ = strconv.Atoi(args[4])

	state.MoveCount, _ = strconv.Atoi(args[5])

	state.GenPieceLists()
	return state
}