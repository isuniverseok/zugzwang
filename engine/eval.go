package engine

import (
	"math"
	"zugzwang/game"
)

const CheckmateEval = 1000000
const CheckmateEvalSplit = 900000 //if absolute value is above this, then it is a checkmate eval. otherwise, normal or stalemate eval
const BigEval = math.MaxInt / 2

var PieceValues map[game.Piece]int = map[game.Piece]int{
	game.Pawn:   100,
	game.Bishop: 300,
	game.Knight: 300,
	game.Rook:   500,
	game.Queen:  900,
}

func GetKingStatus(state *game.State, side game.Side) (bool, bool) {
	kingPos := state.KingPos[game.SideToInd[side]]
	oppSide := game.OppSide(game.PieceSide(state.Board[kingPos.X][kingPos.Y]))
	isInCheck := state.IsAttacked(kingPos, oppSide)
	isSurrounded /*by check*/ := true
	for _, dir := range game.KingDirs {
		newPos := game.Pos{kingPos.X + dir.X, kingPos.Y + dir.Y}
		if game.IsOnBoard(newPos) && state.Board[newPos.X][newPos.Y]&side == 0 { //square is empty or contains an opponent's piece
			if !state.IsAttacked(newPos, oppSide) {
				isSurrounded = false
				break
			}
		}
	}
	return isInCheck, isSurrounded
}

// TODO: turn off piece tables for endgame
func Eval(state *game.State) (int, bool) { //int is eval, bool is if decisive
	side := state.SideToMove
	numMoves := len(state.GenMoves())
	if state.FiftyCount == 50 {
		return 0, true
	}

	if len(state.PieceLists[0]) == 1 && len(state.PieceLists[1]) == 1 {
		return 0, true
	}

	if numMoves == 0 {
		isOppInCheck := state.IsAttacked(state.KingPos[game.SideToInd[game.OppSide(side)]], side)

		if !isOppInCheck { //stalemate
			return 0, true
		} else { //checkmate
			return CheckmateEval, true
		}
	}

	materialEval := 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			pieceValue := PieceValues[game.PieceOnly(state.Board[i][j])]
			if state.Board[i][j]&side > 0 {
				materialEval += pieceValue
			} else {
				materialEval -= pieceValue
			}
		}
	}

	pieceMapEval := 0
	for _, side := range game.Sides {
		for _, p := range state.PieceLists[game.SideToInd[side]] {
			rank := p.X
			if side == game.Black {
				rank = 7 - p.X
			}
			posEval := pieceMaps[game.PieceOnly(state.Board[p.X][p.Y])][rank][p.Y]
			if side == state.SideToMove {
				pieceMapEval += posEval
			} else {
				pieceMapEval -= posEval
			}
		}
	}

	return materialEval + pieceMapEval/2, false
}

//position startpos moves b1c3 e7e6 e2e4 b8c6 d2d4 f8b4 g1e2 g8f6 e4e5 f6e4 c1e3 d7d5 d1d3 a7a5 e1c1 c8d7 c3e4 d5e4 d3e4 h7h5 c1b1 a5a4 a2a3 f7f5 e5f6 b4d6 e4g6 e8f8 g6g7 f8e8 f6f7 e8e7