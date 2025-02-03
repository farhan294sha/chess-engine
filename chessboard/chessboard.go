package chessboard

import "unicode"

type ChessBoard struct {
	whitePawn   uint64
	whiteRook   uint64
	whiteKnight uint64
	whiteBishop uint64
	whiteQueen  uint64
	whiteKing   uint64

	blackPawn   uint64
	blackRook   uint64
	blackKnight uint64
	blackBishop uint64
	blackQueen  uint64
	blackKing   uint64

	// 0 0 0 0 for bit for each bit represnt the castel right of q side k side and for black and white
	// 1 0 0 0 -> blackQueen side allowed
	// 1 1 0 0 -> blackQueen side allowed and blackKing side allowed lise so
	castlingRights uint

	boardMap        map[string]int
	knightAttackMap map[int]uint64
	kingAttackMap   map[int]uint64
	pieceMap        map[rune]*uint64

	blackBoard uint64
	whiteBoard uint64

	isWhitePlay bool
}

func (cb *ChessBoard) Initialize(fen string) {

	cb.generateKnightMoveMap()
	cb.generateKingMovesMap()
	cb.initializeBoardMap()
	cb.fentoPos(fen)
	cb.update()
}
func (cb *ChessBoard) initializeBoardMap() {
	cb.boardMap = map[string]int{
		"a8": 0, "b8": 1, "c8": 2, "d8": 3, "e8": 4, "f8": 5, "g8": 6, "h8": 7,
		"a7": 8, "b7": 9, "c7": 10, "d7": 11, "e7": 12, "f7": 13, "g7": 14, "h7": 15,
		"a6": 16, "b6": 17, "c6": 18, "d6": 19, "e6": 20, "f6": 21, "g6": 22, "h6": 23,
		"a5": 24, "b5": 25, "c5": 26, "d5": 27, "e5": 28, "f5": 29, "g5": 30, "h5": 31,
		"a4": 32, "b4": 33, "c4": 34, "d4": 35, "e4": 36, "f4": 37, "g4": 38, "h4": 39,
		"a3": 40, "b3": 41, "c3": 42, "d3": 43, "e3": 44, "f3": 45, "g3": 46, "h3": 47,
		"a2": 48, "b2": 49, "c2": 50, "d2": 51, "e2": 52, "f2": 53, "g2": 54, "h2": 55,
		"a1": 56, "b1": 57, "c1": 58, "d1": 59, "e1": 60, "f1": 61, "g1": 62, "h1": 63,
	}
}

func (cb *ChessBoard) fentoPos(fen string) {
	cb.pieceMap = map[rune]*uint64{
		'P': &cb.whitePawn,
		'p': &cb.blackPawn,
		'R': &cb.whiteRook,
		'r': &cb.blackRook,
		'N': &cb.whiteKnight,
		'n': &cb.blackKnight,
		'B': &cb.whiteBishop,
		'b': &cb.blackBishop,
		'Q': &cb.whiteQueen,
		'q': &cb.blackQueen,
		'K': &cb.whiteKing,
		'k': &cb.blackKing,
	}
	index := 0
	for _, char := range fen {

		if char == '/' {
			continue
		}
		if unicode.IsDigit(char) {
			index += int(char - '0')
			continue
		}
		if piece, ok := cb.pieceMap[char]; ok {
			setBit(piece, index)
			index++
		}

	}

}

func (cb *ChessBoard) update() {
	cb.blackBoard = cb.blackPawn | cb.blackRook | cb.blackKnight | cb.blackBishop | cb.blackQueen | cb.blackKing
	cb.whiteBoard = cb.whitePawn | cb.whiteRook | cb.whiteKnight | cb.whiteBishop | cb.whiteQueen | cb.whiteKing
	cb.castlingRights = WKCR | WQCR | BQCR | BKCR
}
