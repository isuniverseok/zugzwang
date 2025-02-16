package game

type Pos struct {
	X int
	Y int
}

type Move struct {
	Start     Pos
	End       Pos
	Promotion Piece
}