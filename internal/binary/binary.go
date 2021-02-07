package binary

import (
	"errors"
	"fmt"
	"math"
)

// Converts a bytes to bits
func Bits(bytes []byte) (bits []byte) {
	bits = make([]byte, len(bytes)*8)
	for i, oneByte := range bytes {
		for j := 0; j < 8; j++ {
			bits[i*8+j] = byte(oneByte>>j) & 0x01
		}
	}
	return bits
}

// Converts a bits to bytes
func Bytes(bits []byte) (bytes []byte, err error) {
	slices, err := slicesOf(bits, 8)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert bits to bytes: %v", err)
	}

	bytes = make([]byte, len(slices))
	for i, _byte := range slices {
		var res byte
		for j, bit := range _byte {
			res |= bit << j
		}
		bytes[i] = res
	}

	return bytes, nil
}

// Splits input in slices of length
func slicesOf(input []byte, length int) (output [][]byte, err error) {
	inputLen := len(input)
	if length < 1 {
		return nil, errors.New("Length of sub-slices should be greater than zero")
	}

	if inputLen < length {
		return nil, errors.New("Length of sub-slices should be smaller than then length of the input slice")
	}

	nrSlices := int(math.Ceil(float64(inputLen) / float64(length)))
	output = make([][]byte, nrSlices)

	for i := 0; i < nrSlices; i++ {
		begin := i * length
		end := (i + 1) * length
		if end < inputLen {
			output[i] = input[begin:end]
		} else {
			output[i] = input[begin:]
		}
	}

	return output, nil
}
