package merklehellman

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/crodriguezvega/go-knapsackcrypto/internal/binary"
	"github.com/crodriguezvega/go-knapsackcrypto/internal/number"
)

// Number of tries to generate coprime in generateAandB
const tries = 5

type PublicKey struct {
	m []*big.Int
}

type PrivateKey struct {
	r    []*big.Int
	a, b *big.Int
}

// Generates private and public keys of byte length byteSize.
// This function reports an error if it cannot generate the keys.
func GenerateKeys(byteSize uint) (privKey *PrivateKey, pubKey *PublicKey, err error) {
	if byteSize < 2 {
		return nil, nil, errors.New("Length of public key should be greater than 1 byte")
	}

	bitSize := 8 * byteSize
	n := int(math.Ceil(math.Sqrt(float64(bitSize) / 2.0)))

	r, err := number.GenerateSuperIncreasingSequence(n)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate super increasing sequence: %v", err)
	}

	a, b, err := generateAandB(r[n-1])
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate coprime numbers A and B: %v", err)
	}

	m := generateM(a, b, r)

	privKey = &PrivateKey{r, a, b}
	pubKey = &PublicKey{m}

	return privKey, pubKey, nil
}

// Encrypts the plaintext message p.
// This function reports an error if the bit length of the plaintext
// message is greater than the length of the public key.
func (pubKey *PublicKey) Encrypt(p []byte) (c *big.Int, err error) {
	bits := binary.Bits(p)

	// Check length of the plaintext is <= length of public key
	if len(bits) > len(pubKey.m) {
		return nil, errors.New("Failed to encrypt: length of message in bits should less or equal than length of punlic key")
	}

	c = big.NewInt(0)
	mul := big.NewInt(0)
	for i, b := range bits {
		bit := big.NewInt(int64(b))
		mul.Mul(bit, pubKey.m[i])
		c.Add(c, mul)
	}

	return c, nil
}

// Decrypts the ciphertext message c.
func (privKey *PrivateKey) Decrypt(c *big.Int) (p []byte, err error) {
	s := big.NewInt(0)
	mul := big.NewInt(0)
	invA := big.NewInt(0)

	invA.ModInverse(privKey.a, privKey.b)
	mul.Mul(invA, c)
	s.Mod(mul, privKey.b)

	length := len(privKey.r)
	bits := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		ri := privKey.r[i]

		if s.Cmp(ri) > -1 {
			bits[i] = 1
			s.Sub(s, ri)
		} else {
			bits[i] = 0
		}
	}

	p, err = binary.Bytes(bits)
	if err != nil {
		return nil, fmt.Errorf("Failed to decrypt: %v", err)
	}

	return p, nil
}

// Generates two coprime numbers A and B. Prime A will have the same
// bit length as the bit length of r[n]. Prime B will have the same bit
// length as the bit length of r[n] + 2.
func generateAandB(rn *big.Int) (a *big.Int, b *big.Int, err error) {
	gcd := big.NewInt(0)
	one := big.NewInt(1)
	bitLen := rn.BitLen()

	a, err = rand.Prime(rand.Reader, bitLen)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate prime number A: %v", err)
	}

	// Try up several times to find prime B that is coprime with A
	for i := tries; i > 0; i-- {
		b, err = rand.Prime(rand.Reader, bitLen+2)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to generate prime number B: %v", err)
		}

		// Extra check to make sure A and B are coprime
		if gcd.GCD(nil, nil, a, b).Cmp(one) == 0 {
			return a, b, nil
		}
	}

	return nil, nil, fmt.Errorf("Failed to generate coprime number B after %d attempts", tries)
}

// Generates a sequence m where each element m[i] is calculated as a*b mod r[i].
func generateM(a *big.Int, b *big.Int, r []*big.Int) (m []*big.Int) {
	mul := big.NewInt(0)
	m = make([]*big.Int, len(r))

	for i, ri := range r {
		mul.Mul(a, ri)
		m[i] = big.NewInt(0).Mod(mul, b)
	}

	return m
}
