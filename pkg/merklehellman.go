package merklehellman

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/crodriguezvega/go-knapsackcrypto/number"
)

type PublicKey struct {
	M []*big.Int
}

func newPublicKey(n uint) (pubKey *PublicKey) {
	pubKey = new(PublicKey)

	r := number.GenerateSuperIncreasingSequence(n)
	a, b, err := generateAandB(r[n-1])
	if err != nil {

	}
	pubKey.M = generateM(a, b, r)

	return pubKey
}

type PrivateKey struct {
	M []big.Int
}

func GenerateKeys(byteSize uint) (pubKey *PublicKey, privKey *PrivateKey) {
	bitSize := 8 * byteSize
	n := uint(math.Round(math.Sqrt(float64(bitSize) / 2.0)))

	pubKey = newPublicKey(n)
	return
}

func (pubKey *PublicKey) Encrypt(m []byte) (c *big.Int, err error) {
	// Check length of the plaintext is <= length of public key
	if len(m) > len(pubKey.M) {
		return nil, errors.New("error")
	}

	bits := bits(m)
	var mul, c *big.Int
	for i, mi := range pubKey.M {
		bit := big.NewInt(bits[i])
		mul.Mul(bit, mi)
		c.Add(c, mul)
	}

	return c, nil
}

func generateAandB(rn *big.Int) (a *big.Int, b *big.Int, err error) {
	var gcd *big.Int
	bitLen := rn.BitLen()
	one := big.NewInt(1)

	a, err = rand.Prime(rand.Reader, bitLen)
	if err != nil {
		return nil, nil, err
	}
	// Try up to 5 times to find prime B that is coprime with A
	for i := 5; i > 0; i-- {
		b, err = rand.Prime(rand.Reader, bitLen+2)
		if err != nil {
			return nil, nil, err
		}

		if gcd.GCD(a, b, nil, nil).Cmp(one) == 1 {
			return a, b, nil
		}
	}

	return nil, nil, errors.New("error")
}

func generateM(a *big.Int, b *big.Int, r []*big.Int) (m []*big.Int) {
	var mul, mod *big.Int
	m = make([]*big.Int, len(r))

	for i, ri := range r {
		mul.Mul(a, ri)
		m[i] = mod.Mod(mul, b)
	}

	return m
}

func bits(bytes []byte) (bits []int) {
	bits = make([]int, len(bytes)*8)
	for i, b := range bytes {
		for j := 0; j < 8; j++ {
			bits[i*8+j] = int(b>>j) & 0x01
		}
	}
	return bits
}

func Hello() {
	fmt.Println("hello")

	if num, err := rand.Prime(rand.Reader, 300); err != nil {
		fmt.Println("error")
	} else {
		fmt.Println(num)
	}
}
