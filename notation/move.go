package notation

import (
	"fmt"
	"zugzwang/game"
)

func ParsePosString(posString string) *game.Pos {
	return &game.Pos{X: RankToInt(posString[1]), Y: FileToInt(posString[0])}
}

func ParseMoveString(moveString string, side game.Side) *game.Move {
	move := game.Move{Start: game.Pos{X: RankToInt(moveString[1]), Y: FileToInt(moveString[0])}, End: game.Pos{X: RankToInt(moveString[3]), Y: FileToInt(moveString[2])}}
	if len(moveString) == 5 {
		move.Promotion = ByteToPiece[moveString[4]] //black by default
		move.Promotion -= game.Black
		move.Promotion += side
	}
	return &move
}

func MoveToUCIString(move *game.Move) string {
	result := fmt.Sprintf("%v%v%v%v", string(IntToFile(move.Start.Y)), string(IntToRank(move.Start.X)), string(IntToFile(move.End.Y)), string(IntToRank(move.End.X)))

	if move.Promotion != game.NilPiece {
		promotionPiece := move.Promotion
		if promotionPiece&game.White > 0 {
			promotionPiece -= game.White
			promotionPiece += game.Black
		}
		result += string(PieceToByte[promotionPiece])
	}
	return result
}