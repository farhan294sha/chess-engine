# Chess Engine in Go
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A chess engine implementation in Go using bitboards for efficient move generation and board representation. Currently supports basic FEN parsing (piece positions only) and piece movement logic.

---

## Features

- **Bitboard-based representation** for all chess pieces
- **FEN string parsing** (partial implementation for piece positions)
- Move generation for:
  - Pawns (basic moves and attacks)
  - Knights (precomputed attack tables)
  - Bishops/Rooks/Queens (sliding piece logic)
  - Kings (precomputed attack tables)
- Attack detection and check validation
- Terminal-based board visualization

## Installation

```bash
git clone https://github.com/your-username/chess-engine-go.git
cd chess-engine-go
go build -o chess-engine-+
```

## Usage
```bash
./chess-engine
./chess-engine 61  # Shows possible moves for piece at position 61 (A2)
```

## Supported FEN Format (Partial)
```
rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
```
Note: *Only the piece positions (first field) are currently parsed.*

## Current Limitations

- ❌ **Full FEN Support**: Missing support for castling rights, en passant, half-move clock, and full-move number.
- ❌ **Castling and En Passant Moves**: Castling and en passant moves are not yet implemented.
- ❌ **Move Validation**: No validation for checkmate, stalemate, or illegal moves (e.g., moving into check).
- ❌ **Endgame Detection**: No detection of endgame conditions (e.g., insufficient material, threefold repetition).
- ❌ **UCI Protocol Support**: The engine does not yet support the Universal Chess Interface (UCI) protocol for integration with GUIs.
- ❌ **Move prediction**: No move prediction or search algorithm (e.g., minimax, alpha-beta pruning) implemented.

## Development Roadmap

- [ ] **Complete FEN Parsing Implementation**: Add support for all FEN fields (castling rights, en passant, half-move clock, and full-move number).
- [ ] **Castling Rights Validation**: Implement castling move generation and validation.
- [ ] **En Passant Move Generation**: Add support for en passant captures.
- [ ] **Perft Testing Suite**: Implement a performance testing suite to validate move generation accuracy.
- [ ] **Basic Alpha-Beta Search Implementation**: Add a basic search algorithm with alpha-beta pruning for move prediction.
- [ ] **UCI Protocol Support**: Integrate the Universal Chess Interface (UCI) protocol for compatibility with chess GUIs.



