package mathh

import (
	"fmt"
	"math"
)

func ExampleRound_pi1() {
	rounded := Round(math.Pi, 1)
	fmt.Printf("%.1f", rounded)
	// Output:
	// 3.1
}

func ExampleRound_pi5() {
	rounded := Round(math.Pi, 5)
	fmt.Printf("%.5f", rounded)
	// Output:
	// 3.14159
}

func ExampleRound_n2l() {
	rounded := Round(1.494, 2)
	fmt.Printf("%.2f", rounded)
	// Output:
	// 1.49
}

func ExampleRound_n2h() {
	rounded := Round(1.495, 2)
	fmt.Printf("%.2f", rounded)
	// Output:
	// 1.50
}
