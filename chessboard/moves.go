package chessboard

import "math/bits"

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

// MAPS
func (cb *ChessBoard) generateKnightMoveMap() {
	if cb.knightAttackMap == nil {
		cb.knightAttackMap = make(map[int]uint64)
	}
	for i := 0; i < 64; i++ {
		knightBoard := 1 << (63 - i)
		move := knightMove(uint64(knightBoard))
		cb.knightAttackMap[i] = move
	}
}

func (cb *ChessBoard) generateKingMovesMap() {
	if cb.kingAttackMap == nil {
		cb.kingAttackMap = make(map[int]uint64)
	}
	for i := 0; i < 64; i++ {
		kingBoard := 1 << (63 - i)
		cb.kingAttackMap[i] = kingMove(uint64(kingBoard))
	}
}

//
// chess pieces Moves
//

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

func knightMove(knightBoard uint64) uint64 {
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

func kingMove(kingBoard uint64) uint64 {
	return noWeOne(kingBoard) | northOne(kingBoard) | noEaOne(kingBoard) | eastOne(kingBoard) | soEaOne(kingBoard) | southOne(kingBoard) | soWeOne(kingBoard) | westOne(kingBoard)
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
