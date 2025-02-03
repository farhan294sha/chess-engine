package main

import (
	"chess-engine/chessboard"
	"fmt"
	"os"
	"strconv"
)

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

func print64bitBinary(data uint64, pos int) {
	binaryString := fmt.Sprintf("%064b", data)
	printBinaryWithNewlines(binaryString, pos)
}

// "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func main() {
	var cb chessboard.ChessBoard
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
	cb.Initialize(fen)

	Args := os.Args
	var pos int
	pos, err := strconv.Atoi(Args[1])
	if err != nil {
		return
	}
	fmt.Println(pos)

	//	mask := 1 << (63 - 61)
	//
	//	f1 := cb.whiteBishop & uint64(mask)
	//
	//	setBit(&f1, pos)
	//	clearBit(&f1, 61)
	//
	//	moves := cb.queenMove(f1)
	//
	//	cb.canCastleQueenSide(WHITE)

	fmt.Println(" ")

	fmt.Println(0b0011)

	//print64bitBinary(moves, pos)

}

//https://www.chessprogramming.org/Efficient_Generation_of_Sliding_Piece_Attacks
