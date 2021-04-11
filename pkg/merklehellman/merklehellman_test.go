package merklehellman

import (
	"math"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func Test(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSize = 8
	parameters.MaxSize = 64
	properties := gopter.NewProperties(parameters)

	properties.Property("Encryption and decryption", prop.ForAll(
		func(p []byte, iterations int) bool {
			byteSize := uint(math.Ceil(math.Pow(float64(len(p)*8), 2) / 4))

			privKey, pubKey, err := GenerateKeys(byteSize, 1)
			c, err := pubKey.Encrypt(p)
			m, err := privKey.Decrypt(c)

			return err == nil && reflect.DeepEqual(p, m[:len(p)])
		},
		gen.SliceOf(gen.UInt8Range(0, 255), reflect.TypeOf(byte(0))),
		gen.IntRange(1, math.MaxUint8/8),
	))

	properties.TestingRun(t)
}
