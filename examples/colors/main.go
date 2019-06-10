package main

import (
	"fmt"
	"time"

	"github.com/superhawk610/bar"
	"github.com/ttacon/chalk"
)

func main() {
	n := 20
	b := bar.NewWithOpts(
		bar.WithDimensions(n, 30),
		bar.WithFormat(
			fmt.Sprintf(
				"   %sloading...%s :percent :bar %s:rate ops/s%s ",
				chalk.Blue,
				chalk.Reset,
				chalk.Green,
				chalk.Reset,
			),
		),
	)

	fmt.Println()
	fmt.Println()

	for i := 0; i < n; i++ {
		b.Tick()
		time.Sleep(500 * time.Millisecond)
	}

	b.Done()

	fmt.Println()
	fmt.Println()
}
