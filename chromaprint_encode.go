package sa

import (
	"encoding/base64"
	"encoding/binary"
)

// EncodeFingerprint converts a raw uint32 slice into a compressed Chromaprint string.
// This implements the custom bit-packing algorithm used by AcoustID.
func EncodeFingerprint(fp []uint32) string {
	if len(fp) == 0 {
		return ""
	}

	// 1. Create the header (4 bytes, Little Endian)
	// Algorithm version 1 is standard.
	header := make([]byte, 4)
	binary.LittleEndian.PutUint32(header, uint32(1))

	// 2. Bit-packing the XOR differences
	var data []byte
	var currentByte byte
	var bitPos int

	// Helper to write bits to the byte slice (LSB first)
	writeBits := func(val uint32, bits int) {
		for i := 0; i < bits; i++ {
			if (val>>uint(i))&1 != 0 {
				currentByte |= (1 << uint(bitPos))
			}
			bitPos++
			if bitPos == 8 {
				data = append(data, currentByte)
				currentByte = 0
				bitPos = 0
			}
		}
	}

	lastV := uint32(0)
	for _, v := range fp {
		diff := v ^ lastV
		if diff == 0 {
			// Write '0' bit for no change
			writeBits(0, 1)
		} else {
			// Write '1' bit followed by range and value
			writeBits(1, 1)
			if diff < (1 << 3) {
				writeBits(0, 2) // Range 00: 3 bits
				writeBits(diff, 3)
			} else if diff < (1 << 6) {
				writeBits(1, 2) // Range 01: 6 bits
				writeBits(diff, 6)
			} else if diff < (1 << 12) {
				writeBits(2, 2) // Range 10: 12 bits
				writeBits(diff, 12)
			} else {
				writeBits(3, 2) // Range 11: 32 bits
				writeBits(diff, 32)
			}
		}
		lastV = v
	}

	// Flush remaining bits
	if bitPos > 0 {
		data = append(data, currentByte)
	}

	// 3. Combine header and data, then Base64 encode
	fullData := append(header, data...)
	return base64.StdEncoding.EncodeToString(fullData)
}
