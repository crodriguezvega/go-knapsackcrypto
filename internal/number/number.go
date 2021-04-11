package number

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// Number of tries to generate coprime pair in generateCoprimes
const tries = 5

// Generates a super increasing sequence of length n.
// This function reports an error if n < 2.
func GenerateSuperIncreasingSequence(n int) (r []*big.Int, err error) {
	if n < 2 {
		return nil, errors.New("Length of super increasing sequence should be greater than 1")
	}

	aux := big.NewInt(0)
	one := big.NewInt(1)
	two := big.NewInt(2)
	twoExpN := big.NewInt(0)

	r = make([]*big.Int, n)
	// Set first element of sequence to a value >= 2^n
	twoExpN.Exp(two, big.NewInt(int64(n)), nil)
	offset := generateRandomNumber(one, aux.Sqrt(twoExpN))
	r[0] = offset.Add(offset, twoExpN)

	for i := 1; i < n; i++ {
		offset = generateRandomNumber(one, aux.Sqrt(r[i-1]))
		// Next element of the sequence must be at least 2 times larger than the previous one
		aux.Mul(two, r[i-1])
		r[i] = offset.Add(offset, aux)
	}

	return r, nil
}

// Generates two coprime numbers A and B. Random number A will be such that
// 1 <= A <= B - 1 . Prime B will have the same bit length as the bit length
// of n + 1.
func GenerateCoprimes(n *big.Int) (a *big.Int, b *big.Int, err error) {
	gcd := big.NewInt(0)
	one := big.NewInt(1)
	bitLen := n.BitLen()

	b, err = rand.Prime(rand.Reader, bitLen+1)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate prime number: %v", err)
	}

	// Try several times to find random number A that is coprime with B
	for i := tries; i > 0; i-- {
		a = generateRandomNumber(one, b)

		// Check to make sure A and B are coprime
		if gcd.GCD(nil, nil, a, b).Cmp(one) == 0 {
			return a, b, nil
		}
	}

	return nil, nil, fmt.Errorf("Failed to generate coprime number after %d attempts", tries)
}

// Generates a random number in the range [min, max).
// This function panics if a random number cannot be generated.
func generateRandomNumber(min *big.Int, max *big.Int) (num *big.Int) {
	diff := big.NewInt(0)
	if num, err := rand.Int(rand.Reader, diff.Sub(max, min)); err != nil {
		panic(err)
	} else {
		num.Add(num, min)
		return num
	}
}

// Generates a sequence m where each element m[i] is calculated as r[i]*a mod b.
func MulMod(r []*big.Int, a *big.Int, b *big.Int) (m []*big.Int) {
	mul := big.NewInt(0)
	m = make([]*big.Int, len(r))

	for i, ri := range r {
		mul.Mul(a, ri)
		m[i] = big.NewInt(0).Mod(mul, b)
	}

	return m
}
