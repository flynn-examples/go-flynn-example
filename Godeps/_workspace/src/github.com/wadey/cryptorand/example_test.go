package cryptorand_test

import (
	"fmt"
	"math/rand"

	"github.com/flynn-examples/go-flynn-example/Godeps/_workspace/src/github.com/wadey/cryptorand"
)

func Example() {
	r := rand.New(cryptorand.Source)
	fmt.Println(r.Float64() == r.Float64())

	// Output:
	// false
}
