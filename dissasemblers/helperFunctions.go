package main

// processes 2's compliment nunbers
func Imm_to_32bit_converter(num uint32, bitsize uint) int32 {

	var negBitMask uint32
	var extendMask uint32

	if bitsize == 12 {
		negBitMask = 0x800 // figure out if 12 bit num is neg
		extendMask = 0xFFFFF000

	} else if bitsize == 16 {
		negBitMask = 0x8000 // figure out if 16 bit num is neg
		extendMask = 0xFFFF0000

	} else if bitsize == 19 {
		negBitMask = 0x40000 // figure out if 19 bit num is neg
		extendMask = 0xFFF80000

	} else if bitsize == 26 {
		negBitMask = 0x2000000 // figure out if 26 bit num is neg
		extendMask = 0xFC000000

	} else if bitsize == 32 {
		negBitMask = 0x10000000 // figure out if 26 bit num is neg
		extendMask = 0x00000000

	} else {
		print(" You ARE USING AN INVALID BIT LENGTH")
	}

	var snum int32
	snum = int32(num)
	if (negBitMask & num) > 0 { // is it?
		num = num | extendMask // if so extend with 1's
		num = num ^ 0xFFFFFFFF // 2s comp
		snum = int32(num + 1)
		snum = snum * -1 // add neg sign
	}
	return snum

}
