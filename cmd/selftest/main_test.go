package selftest

import (
	"crypto/sha256"
	"fmt"
)

func ExampleCountDiffLetters() {

	sumx := sha256.Sum256([]byte("x"))
	sumX := sha256.Sum256([]byte("X"))

	fmt.Println(countDiffLetters(sumx, sumX))
	fmt.Println(countDiffLetters(sumx, sumx))
	// Output:
	// 31
	// 0
}
