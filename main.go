package main

import (
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"unicode"
)

const (
	notAfile  uint64 = 0x7F7F7F7F7F7F7F7F // For left shifts (avoids wraparound)
	notABfile uint64 = 0x3F3F3F3F3F3F3F3F // For two left shifts
	notHfile  uint64 = 0xFEFEFEFEFEFEFEFE // For right shifts
	notHGfile uint64 = 0xFCFCFCFCFCFCFCFC // For two right shifts
)

const (
	WKCR = 1 // 0 0 0 1
	WQCR = 2 // 0 0 1 0
	BKCR = 4 // 0 1 0 0
	BQCR = 8 // 1 0 0 0
)

const (
	BLACK = iota
	WHITE
)

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

func (cb *ChessBoard) initialize(fen string) {

	generateKnightMoveMap(cb)
	generateKingMoves(cb)
	initializeBoardMap(cb)
	generateKingMoves(cb)
	fentoPos(fen, cb)
	cb.update()
}

func westAttacks(p, empty uint64) uint64 {

	flood := p
	empty &= notAfile

	for i := 0; i < 6; i++ {
		p = (p << 1) & empty
		flood |= p
	}
	return (flood << 1) & notAfile
}

func eastAttacks(p, empty uint64) uint64 {

	flood := p
	empty &= notHfile

	for i := 0; i < 6; i++ {
		p = (p >> 1) & empty
		flood |= p
	}
	return (flood >> 1) & notHfile
}

func northAttacks(p, empty uint64) uint64 {

	flood := p

	for i := 0; i < 6; i++ {
		p = (p << 8) & empty
		flood |= p
	}
	return (flood << 8)
}

func southAttacks(p, empty uint64) uint64 {

	flood := p

	for i := 0; i < 6; i++ {
		p = (p >> 8) & empty
		flood |= p
	}
	return (flood >> 8)
}

func southOne(b uint64) uint64 { return b >> 8 }
func northOne(b uint64) uint64 { return b << 8 }
func eastOne(b uint64) uint64  { return (b >> 1) & notAfile }
func noEaOne(b uint64) uint64  { return (b << 7) & notAfile }
func soEaOne(b uint64) uint64  { return (b >> 9) & notAfile }
func westOne(b uint64) uint64  { return (b << 1) & notHfile }
func soWeOne(b uint64) uint64  { return (b >> 7) & notHfile }
func noWeOne(b uint64) uint64  { return (b << 9) & notHfile }

func setBit(p *uint64, shift int) {
	shift = 63 - shift
	*p |= 1 << shift
}

func clearBit(board *uint64, position int) {
	position = 63 - position
	*board &= ^(1 << position)
}

func pawnMove(pawnBoard uint64, color int) uint64 {
	var moves uint64

	if color < 0 || color > 1 {
		panic("Invalid color provided for pawn")
	}

	moves |= noEaOne(pawnBoard)
	moves |= noWeOne(pawnBoard)
	if color == WHITE {
		moves |= northOne(pawnBoard)
	} else {
		moves |= southOne(pawnBoard)
	}

	return moves
}

func fentoPos(fen string, cb *ChessBoard) {
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

func initializeBoardMap(cb *ChessBoard) {
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

func printBinaryWithNewlines(binaryStr string, picePos int) {
	// Validate the piece position
	if picePos < 0 || picePos >= len(binaryStr) {
		fmt.Println("Invalid piece position")
		return
	}

	var chunks string
	for i := 0; i < 64; i++ {
		// Colorize the piece at the given position
		if i == picePos {
			chunks += "\033[31m" + string(binaryStr[i]) + "\033[0m" // Red color for the piece
		} else {
			chunks += string(binaryStr[i])
		}

		if (i+1)%8 == 0 {
			fmt.Println(chunks)
			chunks = ""
		}
	}
}

func GenerateKnightMove(knightBoard uint64) uint64 {
	var moves uint64
	moves |= (knightBoard << 15) & notAfile
	moves |= (knightBoard << 17) & notHfile
	moves |= (knightBoard << 10) & notHGfile
	moves |= (knightBoard << 6) & notABfile
	moves |= (knightBoard >> 15) & notHfile
	moves |= (knightBoard >> 17) & notAfile
	moves |= (knightBoard >> 10) & notABfile
	moves |= (knightBoard >> 6) & notHGfile

	//masking for overflow
	return moves
}

func generateKingMoves(cb *ChessBoard) {
	if cb.kingAttackMap == nil {
		cb.kingAttackMap = make(map[int]uint64)
	}
	for i := 0; i < 64; i++ {
		kingBoard := 1 << (63 - i)
		cb.kingAttackMap[i] = kingMove(uint64(kingBoard))
	}
}

func kingMove(kingBoard uint64) uint64 {
	return noWeOne(kingBoard) | northOne(kingBoard) | noEaOne(kingBoard) | eastOne(kingBoard) | soEaOne(kingBoard) | southOne(kingBoard) | soWeOne(kingBoard) | westOne(kingBoard)
}

func generateKnightMoveMap(cb *ChessBoard) {
	if cb.knightAttackMap == nil {
		cb.knightAttackMap = make(map[int]uint64)
	}
	for i := 0; i < 64; i++ {
		knightBoard := 1 << (63 - i)
		move := GenerateKnightMove(uint64(knightBoard))
		cb.knightAttackMap[i] = move
	}
}

func (cb *ChessBoard) bishopMove(bishopBoard uint64) uint64 {

	squareIndex := bits.LeadingZeros64(bishopBoard)
	fileIndex := squareIndex & 7

	occupiedPieces := cb.whiteBoard | cb.blackBoard

	//for West side
	var westSideMoves uint64
	for i := 0; i < fileIndex; i++ {
		westSideMoves |= noWeOne(bishopBoard << (8*i + i))

		westSideMoves |= soWeOne(bishopBoard >> (8*i - i))
	}

	attackPices := occupiedPieces & westSideMoves

	reverse := bits.ReverseBytes64(attackPices) // for loking backwards

	reverse -= bits.ReverseBytes64(bishopBoard) // mask upto the first piece it found from bishop pos

	attackPices -= bishopBoard // mask upto the first piece it found from bishop pos

	attackPices ^= bits.ReverseBytes64(reverse) // this will create mask on both side include the blocker

	westSideMoves &= attackPices & westSideMoves

	// for east side
	var eastSideMoves uint64
	for i := 0; i <= (7 - fileIndex); i++ {
		eastSideMoves |= noEaOne(bishopBoard << (8*i - i))

		eastSideMoves |= soEaOne(bishopBoard >> (8*i + i))
	}

	//same above step on east side (doing spereate to avoid collion on blockers)
	attackPices = occupiedPieces & eastSideMoves

	reverse = bits.ReverseBytes64(attackPices)

	reverse -= bits.ReverseBytes64(bishopBoard)

	attackPices -= bishopBoard

	attackPices ^= bits.ReverseBytes64(reverse)

	eastSideMoves &= attackPices & eastSideMoves

	return westSideMoves | eastSideMoves

}

func (cb *ChessBoard) rookMove(rookBoard uint64) uint64 {

	occupied := cb.whiteBoard | cb.blackBoard

	west := westAttacks(rookBoard, ^occupied)
	east := eastAttacks(rookBoard, ^occupied)
	north := northAttacks(rookBoard, ^occupied)
	south := southAttacks(rookBoard, ^occupied)

	return west | east | north | south

}

func (cb *ChessBoard) queenMove(queen uint64) uint64 {
	dialognalMoves := cb.bishopMove(queen)
	horVrtiMoves := cb.rookMove(queen)
	return dialognalMoves | horVrtiMoves
}

func print64bitBinary(data uint64, pos int) {
	binaryString := fmt.Sprintf("%064b", data)
	printBinaryWithNewlines(binaryString, pos)
}

func (cb *ChessBoard) canCastleKingSide(color int) {

}

func isCheak(cb *ChessBoard) bool {
	var king uint64
	var oppRangePices uint64
	var knight uint64
	var pawn uint64
	if cb.isWhitePlay {
		king = *cb.pieceMap['K']
		oppRangePices = cb.blackBishop | cb.blackQueen | cb.blackRook
		knight = cb.blackKnight
		pawn = cb.blackPawn
	} else {
		king = *cb.pieceMap['k']
		oppRangePices = cb.whiteBishop | cb.whiteQueen | cb.whiteRook
		knight = cb.whiteKnight
		pawn = cb.whitePawn
	}
	// pawn cheak
	kingRayAttacks := noEaOne(king) | noWeOne(king)
	if kingRayAttacks&pawn != 0 {
		return true
	}

	// knight cheak
	indexPos := bits.TrailingZeros64(king)
	kingRayAttacks = cb.kingAttackMap[indexPos]
	if kingRayAttacks&knight != 0 {
		return true
	}

	// ranged pices
	kingRayAttacks = cb.queenMove(king) & oppRangePices
	return kingRayAttacks != 0
}

func (cb *ChessBoard) canCastleQueenSide(color int) bool {
	if color == WHITE {
		if cb.castlingRights&WQCR == 0 {
			return false
		}
		if isCheak(cb) {
			return false
		}

		moveSquare := uint64(1<<(63-cb.boardMap["b1"]) | 1<<(63-cb.boardMap["c1"]) | 1<<(63-cb.boardMap["d1"]))
		rayOfSquare := cb.queenMove(moveSquare)
		if rayOfSquare&(cb.blackBishop|cb.blackQueen|cb.blackRook) == 0 {
			rayOfSquare = cb.knightAttackMap[57] | cb.kingAttackMap[58] | cb.kingAttackMap[59]
			return rayOfSquare&cb.blackKnight == 0
		}
	}
	if color == BLACK {
		if true {
			return false
		}
	}
	return true

}

// "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func main() {
	var cb ChessBoard
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
	cb.initialize(fen)

	Args := os.Args
	var pos int
	pos, err := strconv.Atoi(Args[1])
	if err != nil {
		return
	}

	mask := 1 << (63 - 61)

	f1 := cb.whiteBishop & uint64(mask)

	setBit(&f1, pos)
	clearBit(&f1, 61)

	moves := cb.queenMove(f1)

	cb.canCastleQueenSide(WHITE)

	fmt.Println(" ")

	fmt.Println(0b0011)

	print64bitBinary(moves, pos)

}

//https://www.chessprogramming.org/Efficient_Generation_of_Sliding_Piece_Attacks
