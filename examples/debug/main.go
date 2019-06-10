package main

import (
	"time"

	"github.com/superhawk610/bar"
)

func main() {
	n := 20
	b := bar.NewWithOpts(
		bar.WithDimensions(20, 20),
		bar.WithDebug(),
		bar.WithFormat(" loading... :bar :percent :rate ops/s :custom "),
		bar.WithContext(bar.Context{
			bar.Ctx("custom", "val"),
		}),
	)

	for i := 0; i < n; i++ {
		b.Tick()
		time.Sleep(500 * time.Millisecond)
	}

	b.Done()
}
