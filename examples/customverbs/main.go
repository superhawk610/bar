package main

import (
	"fmt"
	"time"

	"github.com/superhawk610/bar"
)

func main() {
	n := 20
	b := bar.NewWithOpts(
		bar.WithDimensions(n, 30),
		bar.WithFormat(" loading... :percent :bar :rate :hello :world "),
	)

	fmt.Println()
	fmt.Println()

	for i := 0; i < n; i++ {
		b.TickAndUpdate(bar.Context{
			bar.Ctx(":hello", "Hello,"),
			bar.Ctx(":world", "world!"),
		})
		time.Sleep(500 * time.Millisecond)
	}

	b.Done()

	fmt.Println()
	fmt.Println()
}
