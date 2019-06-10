package main

import (
	"fmt"
	"time"

	"github.com/superhawk610/bar"
)

func main() {
	n := 20
	b := bar.New(30)

	fmt.Println()
	fmt.Println()

	for i := 0; i < n; i++ {
		b.Tick()
		if i%2 == 0 {
			b.Interruptf("%d is even!", i)
		}
		time.Sleep(500 * time.Millisecond)
	}

	b.Done()

	fmt.Println()
	fmt.Println()
}
