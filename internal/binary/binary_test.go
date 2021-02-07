package binary

import (
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func Test(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSize = 1
	parameters.MaxSize = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("Convert bytes to bits and back", prop.ForAll(
		func(input []byte) bool {
			bits := Bits(input)
			output, err := Bytes(bits)

			inputLength := len(input)
			outputLength := len(output)

			return err == nil && inputLength == outputLength && reflect.DeepEqual(input, output)
		},
		gen.SliceOf(gen.UInt8(), reflect.TypeOf(byte(0))),
	))

	parameters.MinSize = 8 // any value smaller than 8 results in error
	parameters.MaxSize = 50

	properties.Property("Convert bits to bytes and back", prop.ForAll(
		func(input []byte) bool {
			bytes, err := Bytes(input)
			output := Bits(bytes)

			return err == nil && reflect.DeepEqual(input, output[:len(input)])
		},
		gen.SliceOf(gen.UInt8Range(0, 1), reflect.TypeOf(byte(0))),
	))

	properties.TestingRun(t)
}
