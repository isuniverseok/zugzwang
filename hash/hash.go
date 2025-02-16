package hash

import (
	"math/rand"
	"zugzwang/game"
)

type BitTables struct {
	pieceTable       [2][6][8][8]uint64
	castleRightTable [4]uint64
	enPassantTable   [8][8]uint64
	sideToMoveTable  [2]uint64
}

func (bitTables *BitTables) Init() {
	for i := 0; i < 2; i++ {
		for j := 0; j < 6; j++ {
			for x := 0; x < 8; x++ {
				for y := 0; y < 8; y++ {
					bitTables.pieceTable[i][j][x][y] = rand.Uint64()
				}
			}
		}
	}

	for i := 0; i < 4; i++ {
		bitTables.castleRightTable[i] = rand.Uint64()
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			bitTables.enPassantTable[i][j] = rand.Uint64()
		}
	}
	for i := 0; i < 2; i++ {
		bitTables.sideToMoveTable[i] = rand.Uint64()
	}
}

func NewBitTables() *BitTables {
	newBitTables := &BitTables{}
	newBitTables.Init()
	return newBitTables
}

func Hash(state *game.State, bitTables *BitTables) uint64 { //TODO: return two hashes for less collisions
	var result uint64 = 0

	//pieces
	for i := 0; i < 2; i++ {
		for _, p := range state.PieceLists[i] {
			result ^= bitTables.pieceTable[i][game.PieceToInd[game.PieceOnly(state.Board[p.X][p.Y])]][p.X][p.Y]
		}
	}

	//castle rights
	for _, c := range game.CastleRights {
		if state.CastleRights&c > 0 {
			result ^= bitTables.castleRightTable[game.CastleRightToInd[c]]
		}
	}

	//en passant
	if state.EnPassantPos != nil {
		result ^= bitTables.enPassantTable[state.EnPassantPos.X][state.EnPassantPos.Y]
	}

	//side to move
	for i := 0; i < 4; i++ {
		if state.CastleRights&(1<<i) > 0 {
			result ^= bitTables.castleRightTable[i]
		}
	}
	result ^= bitTables.sideToMoveTable[game.SideToInd[state.SideToMove]]
	//TODO: move count, fifty move count

	return result
}