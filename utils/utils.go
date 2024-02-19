package utils

func GetNibbleAt(value uint16, at int) uint16 {
	return uint16((value & (0xf000 >> (at * 4))) >> ((3 - at) * 4))
}

func GetBitAt(value uint8, at int) uint8 {
	return (value & (0x80 >> at)) >> (7 - at)
}
