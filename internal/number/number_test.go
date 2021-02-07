package number

import (
	"math"
	"math/big"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func Test(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("Super increasing sequence for length >= 2", prop.ForAll(
		func(n int) bool {
			// Generate super increasing sequence of length n
			seq, err := GenerateSuperIncreasingSequence(n)

			isSuperIncreasing := true
			diff := big.NewInt(0) // Difference between two consecutive elements of the super increasing sequence
			two := big.NewInt(2)
			twoExpN := big.NewInt(0)
			twoExpN.Exp(two, big.NewInt(int64(n)), nil)

			// The first element must be larger that 2^n
			if seq[0].Cmp(twoExpN) < 1 {
				isSuperIncreasing = false
			}

			length := len(seq)
			for i := 1; i < length; i++ {
				// The difference must be larger than the vaue of the previous element
				diff.Sub(seq[i], seq[i-1])
				if diff.Cmp(seq[i-1]) < 1 {
					isSuperIncreasing = false
				}

				if !isSuperIncreasing {
					break
				}
			}

			return err == nil && length == int(n) && isSuperIncreasing == true
		},
		gen.IntRange(2, math.MaxUint16/16),
	))

	properties.TestingRun(t)
}
