package merklehellman

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/crodriguezvega/go-knapsackcrypto/internal/binary"
	"github.com/crodriguezvega/go-knapsackcrypto/internal/number"
)

type PublicKey struct {
	m []*big.Int
}

type PrivateKey struct {
	r    []*big.Int
	a, b []*big.Int
}

// Generates private and public keys of byte length byteSize.
// This function reports an error if it cannot generate the keys.
func GenerateKeys(byteSize uint, iterations int) (privKey *PrivateKey, pubKey *PublicKey, err error) {
	if byteSize < 2 {
		return nil, nil, errors.New("Length of public key should be greater than 1 byte")
	}
	if iterations < 1 {
		return nil, nil, errors.New("Number of iterations must be greater than 0")
	}

	bitSize := 8 * byteSize
	n := int(math.Ceil(math.Sqrt(float64(bitSize) / 2.0)))

	r, err := number.GenerateSuperIncreasingSequence(n)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate super increasing sequence: %v", err)
	}

	a, b, m, err := GenerateKeyParameters(r, iterations)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate pub/priv key pair: %v", err)
	}

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
		return nil, errors.New("Failed to encrypt: length of message in bits should less or equal than length of public key")
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
// This function reports an error if the conversion of
// the decrypted bits slice to bytes fails.
func (privKey *PrivateKey) Decrypt(c *big.Int) (p []byte, err error) {
	mul := big.NewInt(0)
	invA := big.NewInt(0)
	s := new(big.Int).Set(c)

	length := len(privKey.a)
	for i := length - 1; i >= 0; i-- {
		invA.ModInverse(privKey.a[i], privKey.b[i])
		mul.Mul(invA, s)
		s.Mod(mul, privKey.b[i])
	}

	length = len(privKey.r)
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

// Generate parameters for the public and private keys
// This function reports an error if it fails to generate a pair of coprime numbers.
func GenerateKeyParameters(r []*big.Int, iterations int) (a []*big.Int, b []*big.Int, m []*big.Int, err error) {
	m = r
	a = make([]*big.Int, iterations)
	b = make([]*big.Int, iterations)

	for i := 0; i < iterations; i++ {
		sum := big.NewInt(0)
		for _, mi := range m {
			sum.Add(sum, mi)
		}

		ai, bi, err := number.GenerateCoprimes(sum)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Failed to generate key parameters: %v", err)
		}

		a[i] = ai
		b[i] = bi
		m = number.MulMod(m, ai, bi)
	}

	return a, b, m, nil
}
