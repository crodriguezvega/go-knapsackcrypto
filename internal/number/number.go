package number

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// Generates a super increasing sequence of length n.
// This function reports an error if n < 2.
func GenerateSuperIncreasingSequence(n int) (seq []*big.Int, err error) {
	if n < 2 {
		return nil, errors.New("Length of super increasing sequence should be greater than 1")
	}

	aux := big.NewInt(0)
	two := big.NewInt(2)
	twoExpN := big.NewInt(0)

	seq = make([]*big.Int, n)
	// Set first element of sequence to a value >= 2^n
	twoExpN.Exp(two, big.NewInt(int64(n)), nil)
	offset := GenerateRandomNumber(two, aux.Sqrt(twoExpN))
	seq[0] = offset.Add(offset, twoExpN)

	for i := 1; i < n; i++ {
		offset = GenerateRandomNumber(two, aux.Sqrt(seq[i-1]))
		// Next element of the sequence must be at least 2 times larger than the previous one
		aux := aux.Mul(two, seq[i-1])
		seq[i] = offset.Add(offset, aux)
	}

	return seq, nil
}

// Generates a random number in the range [min, max).
// This function panics if a random number cannot be generated.
func GenerateRandomNumber(min *big.Int, max *big.Int) (num *big.Int) {
	diff := big.NewInt(0)
	if num, err := rand.Int(rand.Reader, diff.Sub(max, min)); err != nil {
		panic(err)
	} else {
		num.Add(num, min)
		return num
	}
}
