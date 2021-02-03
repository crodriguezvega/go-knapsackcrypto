package number

import (
	"crypto/rand"
	"math/big"
)

func GenerateSuperIncreasingSequence(n uint) (seq []*big.Int) {
	// Range for random number generation
	min := big.NewInt(2)
	max := big.NewInt(100)

	// Set first element of sequence (>= 2^n)
	var low *big.Int
	seq = make([]*big.Int, n)
	mul := generateRandomNumber(min, max)
	low.Exp(big.NewInt(2), big.NewInt(int64(n)), nil)
	seq[0] = mul.Mul(mul, low)

	for i := uint(1); i < n; i++ {
		mul = generateRandomNumber(min, max)
		seq[i] = mul.Mul(mul, seq[i-1])
	}

	return
}

// Generates a random number in the range [min, max).
func GenerateRandomNumber(min *big.Int, max *big.Int) (num *big.Int) {
	if num, err := rand.Int(rand.Reader, max.Sub(max, min)); err != nil {
		panic(err)
	} else {
		num.Add(num, min)
	}

	return
}
