package engine

import (
	"fmt"
	"sort"
	"zugzwang/boolwrapper"
	"zugzwang/deepcopy"
	"zugzwang/game"
	"zugzwang/hash"
	"zugzwang/notation"
	"zugzwang/sliceutils"
	"time"
)

var currNodes = 0

type pvEntry struct {
	move  *game.Move
	eval  int
	depth int
}

var PVS map[uint64]map[uint64]*pvEntry = make(map[uint64]map[uint64]*pvEntry)
var newPVS map[uint64]map[uint64]*pvEntry = make(map[uint64]map[uint64]*pvEntry)

var bitTables1 = hash.NewBitTables()
var bitTables2 = hash.NewBitTables()

func CaptureSearch(state *game.State, alpha, beta int) int {
	currNodes++
	standPat, isDecisive := Eval(state)
	if standPat >= beta {
		return beta
	}
	if standPat > alpha {
		alpha = standPat
	}
	if isDecisive {
		CheckmateEvalSplit := 0
		if standPat < -CheckmateEvalSplit {
			standPat++
		}
		if standPat > CheckmateEvalSplit {
			standPat--
		}
		return standPat
	}
	captures := state.GenPseudoMoves()
	selfSide := state.SideToMove
	for _, capture := range captures {
		capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights := state.RunMove(capture)
		if state.IsAttacked(state.KingPos[game.SideToInd[selfSide]], game.OppSide(selfSide)) || capturedPiece == game.NilPiece {
			state.ReverseMove(capture, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights)
			continue
		}
		eval := -CaptureSearch(state, -beta, -alpha)
		state.ReverseMove(capture, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights)
		if eval >= beta {
			return beta
		}
		if eval > alpha {
			alpha = eval
		}
	}
	// CheckmateEvalSplit := 0
	if alpha < -CheckmateEvalSplit {
		alpha++
	}
	if alpha > CheckmateEvalSplit {
		alpha--
	}
	return alpha
}

func IterativeDeepening(state *game.State, moveChan chan *game.Move, isSearching *boolwrapper.BoolWrapper) {
	currDepth := 2
	moves := state.GenMoves()
	var bestMove *game.Move
	var bestEval int
	for {
		newPVS = make(map[uint64]map[uint64]*pvEntry)
		currNodes = 0
		currDepthStartTime := time.Now()
		BigEval := 0
		bestEval, bestMove = Minimax(state, currDepth, -BigEval, BigEval, isSearching)

		if !isSearching.Val {
			return
		}

		for k1 := range newPVS {
			for k2 := range newPVS[k1] {
				_, hasStoredPV := PVS[k1]
				if !hasStoredPV {
					PVS[k1] = make(map[uint64]*pvEntry)
				}
				PVS[k1][k2] = newPVS[k1][k2]
			}
		}

		var bestMoveInd int
		for i := range moves {
			if moves[i] == bestMove {
				bestMoveInd = i
			}
		}

		moves = sliceutils.RemoveByIndex(moves, bestMoveInd)
		moves = append([]*game.Move{bestMove}, moves...)
		moveChan <- bestMove
		pvString := notation.MoveToUCIString(bestMove)
		stateCopyIface, _ := deepcopy.Anything(state)
		stateCopy := stateCopyIface.(*game.State)
		stateCopy.RunMove(bestMove)
		for {
			pvTable1, ok := PVS[hash.Hash(stateCopy, bitTables1)]
			if !ok {
				break
			}
			pv, ok := pvTable1[hash.Hash(stateCopy, bitTables2)]
			if !ok {
				break
			}
			pvString += " " + notation.MoveToUCIString(pv.move)
			stateCopy.RunMove(pv.move)
		}
		if bestEval > CheckmateEvalSplit {
			fmt.Printf("info depth %v multipv 1 score mate %v nps %v pv %v\n", currDepth, CheckmateEval-bestEval, int(float64(currNodes)/float64(time.Since(currDepthStartTime).Seconds())), pvString)
		} else if bestEval < -CheckmateEvalSplit {
			fmt.Printf("info depth %v multipv 1 score mate %v nps %v pv %v\n", currDepth, bestEval+CheckmateEval, int(float64(currNodes)/float64(time.Since(currDepthStartTime).Seconds())), pvString)
		} else {
			fmt.Printf("info depth %v multipv 1 score cp %v nps %v pv %v\n", currDepth, bestEval, int(float64(currNodes)/float64(time.Since(currDepthStartTime).Seconds())), pvString)
		}
		currDepth++
	}
}

func Minimax(state *game.State, depth int, alpha, beta int, isSearching *boolwrapper.BoolWrapper) (int, *game.Move) {
	currNodes++
	origAlpha, origBeta := alpha, beta
	currSide := state.SideToMove
	currHash1 := hash.Hash(state, bitTables1)
	currHash2 := hash.Hash(state, bitTables2)
	moves := state.GenPseudoMoves()

	bestEval := -BigEval
	var bestMove *game.Move = nil

	sort.Slice(moves, func(a int, b int) bool {
		scoreA := PieceValues[game.PieceOnly(state.Board[moves[a].End.X][moves[a].End.Y])] - PieceValues[game.PieceOnly(state.Board[moves[a].Start.X][moves[a].Start.Y])]/100
		scoreB := PieceValues[game.PieceOnly(state.Board[moves[b].End.X][moves[b].End.Y])] - PieceValues[game.PieceOnly(state.Board[moves[b].Start.X][moves[b].Start.Y])]/100
		return scoreA > scoreB
	})
	pvHash1, hasStoredPV := PVS[currHash1]
	var pv *pvEntry
	if hasStoredPV {
		pv, hasStoredPV = pvHash1[currHash2]
		if hasStoredPV {
			if pv.depth >= depth {
				return pv.eval, pv.move
			} else {
				moves = append([]*game.Move{pv.move}, moves...)
			}
		}
	} else {
		pvHash1, hasStoredPV := newPVS[currHash1]
		if hasStoredPV {
			pv, hasStoredPV = pvHash1[currHash2]
			if hasStoredPV {
				if pv.depth >= depth {
					return pv.eval, pv.move
				} else {
					moves = append([]*game.Move{pv.move}, moves...)
				}
			}
		}
	}
	selfSide := state.SideToMove
	numValidMoves := 0
	for _, m := range moves {
		capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights := state.RunMove(m)
		if state.IsAttacked(state.KingPos[game.SideToInd[selfSide]], game.OppSide(selfSide)) {
			state.ReverseMove(m, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights)
			continue
		}
		numValidMoves++
		var currOppEval int
		if depth == 1 {
			currOppEval = CaptureSearch(state, -beta, -alpha) //TODO: fix capture search
		} else {
			currOppEval, _ = Minimax(state, depth-1, -beta, -alpha, isSearching)
		}
		if !isSearching.Val {
			return 0, nil
		}
		currEval := -currOppEval
		if currEval > bestEval {
			bestEval = currEval
			bestMove = m
		}
		if bestEval > alpha {
			alpha = bestEval
		}
		state.ReverseMove(m, capturedPiece, isEnPassant, oldFiftyCount, oldEnPassantPos, oldCastleRights)
		if alpha >= beta {
			break
		}
	}
	if numValidMoves == 0 {
		if state.IsAttacked(state.KingPos[game.SideToInd[currSide]], game.OppSide(currSide)) {
			bestEval = -CheckmateEval
		} else {
			bestEval = 0
		}
		return bestEval, nil
	}
	if bestEval > CheckmateEvalSplit {
		bestEval--
	}
	if bestEval < -CheckmateEvalSplit {
		bestEval++
	}
	if bestMove != nil && bestEval > origAlpha && bestEval < origBeta && (!hasStoredPV || depth > pv.depth) {
		_, hasStoredInNewPV := newPVS[currHash1]
		if !hasStoredInNewPV {
			newPVS[currHash1] = make(map[uint64]*pvEntry)
		}
		newPVS[currHash1][currHash2] = &pvEntry{bestMove, bestEval, depth}
	}
	return bestEval, bestMove
}