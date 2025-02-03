package chessboard

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
