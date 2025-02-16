package game

import (
	"zugzwang/sliceutils"
)

var (
	KingDirs   = []Pos{{1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	KnightDirs = []Pos{{1, 2}, {2, 1}, {-1, 2}, {-2, 1}, {1, -2}, {2, -1}, {-1, -2}, {-2, -1}}
	BishopDirs = []Pos{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}}
	RookDirs   = []Pos{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	QueenDirs  = []Pos{{1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
)

var (
	KingSideDir  = 1
	QueenSideDir = -1
)

func IsOnBoard(p Pos) bool {
	return p.X >= 0 && p.Y >= 0 && p.X < 8 && p.Y < 8
}

func (state *State) IsAttacked(p Pos, oppSide Side) bool {
	//pawns
	var pawnDirection int
	if oppSide == White {
		pawnDirection = -1
	} else {
		pawnDirection = 1
	}

	if p.X-pawnDirection >= 0 && p.X-pawnDirection < 8 {
		if p.Y >= 1 && state.Board[p.X-pawnDirection][p.Y-1] == oppSide|Pawn {
			return true
		}
		if p.Y <= 6 && state.Board[p.X-pawnDirection][p.Y+1] == oppSide|Pawn {
			return true
		}
	}

	for _, dir := range KingDirs {
		if IsOnBoard(Pos{p.X + dir.X, p.Y + dir.Y}) && state.Board[p.X+dir.X][p.Y+dir.Y] == oppSide|King {
			return true
		}
	}

	for _, dir := range KnightDirs {
		if IsOnBoard(Pos{p.X + dir.X, p.Y + dir.Y}) && state.Board[p.X+dir.X][p.Y+dir.Y] == oppSide|Knight {
			return true
		}
	}

	for _, dir := range BishopDirs {
		newPos := Pos{p.X + dir.X, p.Y + dir.Y}
		for IsOnBoard(newPos) {
			if state.Board[newPos.X][newPos.Y] == oppSide|Bishop || state.Board[newPos.X][newPos.Y] == oppSide|Queen {
				return true
			} else if state.Board[newPos.X][newPos.Y] != NilPiece {
				break
			}
			newPos.X += dir.X
			newPos.Y += dir.Y
		}
	}

	for _, dir := range RookDirs {
		newPos := Pos{p.X + dir.X, p.Y + dir.Y}
		for IsOnBoard(newPos) {
			if state.Board[newPos.X][newPos.Y] == oppSide|Rook || state.Board[newPos.X][newPos.Y] == oppSide|Queen {
				return true
			} else if state.Board[newPos.X][newPos.Y] != NilPiece {
				break
			}
			newPos.X += dir.X
			newPos.Y += dir.Y
		}
	}

	return false
}

func (state *State) GenPseudoMoves() []*Move { //allows the king to be in check
	moves := []*Move{}

	var selfPieceList []Pos = state.PieceLists[SideToInd[state.SideToMove]]
	var pawnDirection int
	var pawnStartRank int
	var pawnPromotionRank int

	if state.SideToMove == White {
		pawnDirection = -1
		pawnStartRank = 6
		pawnPromotionRank = 0
	} else {
		pawnDirection = 1
		pawnStartRank = 1
		pawnPromotionRank = 7
	}

	selfSide := state.SideToMove
	oppSide := OppSide(state.SideToMove)

	appendPawnMove := func(move *Move) { //handles promotion
		if move.End.X == pawnPromotionRank {
			moves = append(moves, &Move{move.Start, move.End, Queen | selfSide})
			moves = append(moves, &Move{move.Start, move.End, Rook | selfSide})
			moves = append(moves, &Move{move.Start, move.End, Bishop | selfSide})
			moves = append(moves, &Move{move.Start, move.End, Knight | selfSide})
		} else {
			moves = append(moves, move)
		}
	}

	generateJumpMoves := func(p Pos, dirs []Pos) {
		for _, dir := range dirs {
			newPos := Pos{p.X + dir.X, p.Y + dir.Y}
			if IsOnBoard(newPos) && (state.Board[newPos.X][newPos.Y] == NilPiece || state.Board[newPos.X][newPos.Y]&oppSide > 0) {
				moves = append(moves, &Move{Start: p, End: newPos})
			}
		}
	}
	generateSlideMoves := func(p Pos, dirs []Pos) {
		for _, dir := range dirs {
			newPos := Pos{p.X + dir.X, p.Y + dir.Y}
			for IsOnBoard(newPos) && (state.Board[newPos.X][newPos.Y] == NilPiece || state.Board[newPos.X][newPos.Y]&oppSide > 0) {
				moves = append(moves, &Move{Start: p, End: newPos})
				if state.Board[newPos.X][newPos.Y]&oppSide > 0 {
					break
				}
				newPos.X += dir.X
				newPos.Y += dir.Y
			}
		}
	}

	for _, p := range selfPieceList {
		switch state.Board[p.X][p.Y] - selfSide {
		case Pawn:
			if state.Board[p.X+pawnDirection][p.Y] == NilPiece {
				appendPawnMove(&Move{Start: p, End: Pos{p.X + pawnDirection, p.Y}})
				if p.X == pawnStartRank && state.Board[p.X+pawnDirection*2][p.Y] == NilPiece {
					appendPawnMove(&Move{Start: p, End: Pos{p.X + pawnDirection*2, p.Y}})
				}
			}
			if p.Y >= 1 && state.Board[p.X+pawnDirection][p.Y-1]&oppSide > 0 {
				appendPawnMove(&Move{Start: p, End: Pos{p.X + pawnDirection, p.Y - 1}})
			}
			if p.Y <= 6 && state.Board[p.X+pawnDirection][p.Y+1]&oppSide > 0 {
				appendPawnMove(&Move{Start: p, End: Pos{p.X + pawnDirection, p.Y + 1}})
			}

			if state.EnPassantPos != nil && state.EnPassantPos.X == p.X && (state.EnPassantPos.Y-p.Y == 1 || state.EnPassantPos.Y-p.Y == -1) {
				appendPawnMove(&Move{Start: p, End: Pos{p.X + pawnDirection, state.EnPassantPos.Y}})
			}
		case King:
			generateJumpMoves(p, KingDirs)
			if !state.IsAttacked(p, oppSide) {
				var KingSide, QueenSide CastleRight
				if selfSide == White {
					KingSide = WhiteKingSide
					QueenSide = WhiteQueenSide
				} else {
					KingSide = BlackKingSide
					QueenSide = BlackQueenSide
				}
				if state.CastleRights&KingSide > 0 && state.Board[p.X][p.Y+1] == NilPiece && state.Board[p.X][p.Y+2] == NilPiece && !state.IsAttacked(p, oppSide) && !state.IsAttacked(Pos{p.X, p.Y + 1}, oppSide) && !state.IsAttacked(Pos{p.X, p.Y + 2}, oppSide) {
					moves = append(moves, &Move{Start: p, End: Pos{p.X, p.Y + 2}})
				}
				if state.CastleRights&QueenSide > 0 && state.Board[p.X][p.Y-1] == NilPiece && state.Board[p.X][p.Y-2] == NilPiece && state.Board[p.X][p.Y-3] == NilPiece && !state.IsAttacked(p, oppSide) && !state.IsAttacked(Pos{p.X, p.Y - 1}, oppSide) && !state.IsAttacked(Pos{p.X, p.Y - 2}, oppSide) {
					moves = append(moves, &Move{Start: p, End: Pos{p.X, p.Y - 2}})
				}
			}
		case Knight:
			generateJumpMoves(p, KnightDirs)
		case Bishop:
			generateSlideMoves(p, BishopDirs)
		case Rook:
			generateSlideMoves(p, RookDirs)
		case Queen:
			generateSlideMoves(p, QueenDirs)
		default:
		}
	}
	return moves
}

func (state *State) RunMove(move *Move) (Piece, bool, int, *Pos, uint8) {
	var capturedPiece Piece
	var isEnPassant bool = false
	var oldFiftyCount int = state.FiftyCount
	var oldEnPassantPos *Pos = state.EnPassantPos

	var oldCastleRights uint8 = state.CastleRights

	//intializing info
	var pawnDirection int
	var pawnStartRank int
	var backRank int

	if state.SideToMove == White {
		pawnDirection = -1
		pawnStartRank = 6
		backRank = 7
	} else {
		pawnDirection = 1
		pawnStartRank = 1
		backRank = 0
	}

	selfStartPiece := state.Board[move.Start.X][move.Start.Y]
	selfEndPiece := state.Board[move.End.X][move.End.Y]
	selfSide := state.SideToMove
	oppSide := OppSide(state.SideToMove)

	//make move

	capturedPos := move.End
	if selfStartPiece&Pawn > 0 && move.Start.Y-move.End.Y != 0 && selfEndPiece == NilPiece {
		isEnPassant = true
		capturedPos = *state.EnPassantPos
	}
	capturedPiece = state.Board[capturedPos.X][capturedPos.Y]

	state.Board[move.End.X][move.End.Y] = state.Board[move.Start.X][move.Start.Y]
	if move.Promotion != NilPiece {
		state.Board[move.End.X][move.End.Y] = move.Promotion
	}

	state.Board[move.Start.X][move.Start.Y] = NilPiece

	//en passant

	if isEnPassant {
		state.Board[state.EnPassantPos.X][state.EnPassantPos.Y] = NilPiece
	}

	//piece lists
	if capturedPiece != NilPiece {
		for i, p := range state.PieceLists[SideToInd[oppSide]] {
			if p == capturedPos {
				state.PieceLists[SideToInd[oppSide]] = sliceutils.RemoveByIndex(state.PieceLists[SideToInd[oppSide]], i)
				break
			}
		}
	}

	for i, p := range state.PieceLists[SideToInd[selfSide]] {
		if p == move.Start {
			state.PieceLists[SideToInd[selfSide]][i] = move.End
			break
		}
	}

	//move counts
	if selfStartPiece&Pawn > 0 || selfEndPiece != NilPiece {
		state.FiftyCount = 0
	} else {
		state.FiftyCount++
	}

	if selfSide == Black {
		state.MoveCount++
	}

	//new en passant square
	if selfStartPiece&Pawn > 0 && move.Start.X == pawnStartRank && move.End.X == pawnStartRank+pawnDirection*2 {
		state.EnPassantPos = &move.End
	} else {
		state.EnPassantPos = nil
	}

	//king pos
	if selfStartPiece&King > 0 {
		state.KingPos[SideToInd[selfSide]] = move.End
	}

	//castle rights
	if selfStartPiece&King > 0 {
		if selfSide == White {
			state.CastleRights &= ^WhiteKingSide
			state.CastleRights &= ^WhiteQueenSide
		} else {
			state.CastleRights &= ^BlackKingSide
			state.CastleRights &= ^BlackQueenSide
		}
	}
	if selfStartPiece&Rook > 0 {
		if move.Start.X == backRank {
			if move.Start.Y == 7 {
				if selfSide == White {
					state.CastleRights &= ^WhiteKingSide
				} else {
					state.CastleRights &= ^BlackKingSide
				}
			} else if move.Start.Y == 0 {
				if selfSide == White {
					state.CastleRights &= ^WhiteQueenSide
				} else {
					state.CastleRights &= ^BlackQueenSide
				}
			}
		}
	}

	if capturedPiece&Rook > 0 {
		if capturedPos.X == 7-backRank {
			if capturedPos.Y == 7 {
				if oppSide == White {
					state.CastleRights &= ^WhiteKingSide
				} else {
					state.CastleRights &= ^BlackKingSide
				}
			} else if capturedPos.Y == 0 {
				if oppSide == White {
					state.CastleRights &= ^WhiteQueenSide
				} else {
					state.CastleRights &= ^BlackQueenSide
				}
			}
		}
	}

	//moving castled rook

	if selfStartPiece&King > 0 {
		if move.End.Y-move.Start.Y == 2 || move.End.Y-move.Start.Y == -2 {
			var rookPos Pos
			var newRookPos Pos
			if move.End.Y-move.Start.Y == 2 {
				rookPos = Pos{backRank, move.End.Y + 1}
				newRookPos = Pos{backRank, move.End.Y - 1}
			}
			if move.End.Y-move.Start.Y == -2 {
				rookPos = Pos{backRank, move.End.Y - 2}
				newRookPos = Pos{backRank, move.End.Y + 1}
			}
			state.Board[newRookPos.X][newRookPos.Y] = state.Board[rookPos.X][rookPos.Y]
			state.Board[rookPos.X][rookPos.Y] = NilPiece
			for i, p := range state.PieceLists[SideToInd[selfSide]] {
				if p == rookPos {
					state.PieceLists[SideToInd[selfSide]][i] = newRookPos
					break
				}
			}
		}
	}

	//side to move
	state.SideToMove = oppSide

	return capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights
}

func (state *State) ReverseMove(move *Move, capturedPiece Piece, isEnPassant bool, oldFiftyCount int, oldEnPassantPos *Pos, oldCastleRights uint8) {
	selfSide := OppSide(state.SideToMove)
	oppSide := state.SideToMove
	var pawnDirection int
	var backRank int

	if selfSide == White {
		pawnDirection = -1
		backRank = 7
	} else {
		pawnDirection = 1
		backRank = 0
	}

	//move pieces
	selfStartPiece := state.Board[move.End.X][move.End.Y]
	if move.Promotion != NilPiece {
		selfStartPiece = Pawn | selfSide
	}

	state.Board[move.Start.X][move.Start.Y] = selfStartPiece

	if !isEnPassant {
		state.Board[move.End.X][move.End.Y] = capturedPiece
	} else {
		state.Board[move.End.X-pawnDirection][move.End.Y] = capturedPiece
		state.Board[move.End.X][move.End.Y] = NilPiece
	}

	//piece lists

	for i, p := range state.PieceLists[SideToInd[selfSide]] {
		if p == move.End {
			state.PieceLists[SideToInd[selfSide]][i] = move.Start
			break
		}
	}

	var capturedPos Pos
	if capturedPiece != NilPiece {

		if !isEnPassant {
			capturedPos = move.End
		} else {
			capturedPos = Pos{move.End.X - pawnDirection, move.End.Y}
		}

		state.PieceLists[SideToInd[oppSide]] = append(state.PieceLists[SideToInd[oppSide]], capturedPos)
	}

	//move counts
	state.FiftyCount = oldFiftyCount

	if selfSide == Black {
		state.MoveCount--
	}

	//en passant square
	state.EnPassantPos = oldEnPassantPos

	if selfStartPiece&King > 0 {
		state.KingPos[SideToInd[selfSide]] = move.Start
	}

	//castling rights
	state.CastleRights = oldCastleRights

	//reverse castling
	if selfStartPiece&King > 0 {
		if move.End.Y-move.Start.Y == 2 || move.End.Y-move.Start.Y == -2 {
			var rookPos Pos
			var newRookPos Pos
			if move.End.Y-move.Start.Y == 2 {
				rookPos = Pos{backRank, move.End.Y + 1}
				newRookPos = Pos{backRank, move.End.Y - 1}
			}
			if move.End.Y-move.Start.Y == -2 {
				rookPos = Pos{backRank, move.End.Y - 2}
				newRookPos = Pos{backRank, move.End.Y + 1}
			}

			state.Board[rookPos.X][rookPos.Y] = state.Board[newRookPos.X][newRookPos.Y]
			state.Board[newRookPos.X][newRookPos.Y] = NilPiece

			for i, p := range state.PieceLists[SideToInd[selfSide]] {
				if p == newRookPos {
					state.PieceLists[SideToInd[selfSide]][i] = rookPos
					break
				}
			}
		}
	}

	//side to move
	state.SideToMove = selfSide
}

func (state *State) GenMoves() []*Move { //all valid moves
	pseudoMoves := state.GenPseudoMoves()
	validMoves := []*Move{}
	selfSide := state.SideToMove

	for _, pseudoMove := range pseudoMoves {
		capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights := state.RunMove(pseudoMove)

		if !state.IsAttacked(state.KingPos[SideToInd[selfSide]], OppSide(selfSide)) {
			validMoves = append(validMoves, pseudoMove)
		}

		state.ReverseMove(pseudoMove, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights)
	}

	return validMoves
}

func (state *State) GenCaptures() []*Move { //all valid moves
	pseudoMoves := state.GenPseudoMoves()
	validMoves := []*Move{}
	selfSide := state.SideToMove

	for _, pseudoMove := range pseudoMoves {
		capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights := state.RunMove(pseudoMove)

		if !state.IsAttacked(state.KingPos[SideToInd[selfSide]], OppSide(selfSide)) && capturedPiece != NilPiece {
			validMoves = append(validMoves, pseudoMove)
		}

		state.ReverseMove(pseudoMove, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights)
	}

	return validMoves
}