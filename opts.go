package bar

// this package supports the functional opts style for incrementally
// adding customization options when creating a new instance - you can
// read more about this here: https://halls-of-valhalla.org/beta/articles/functional-options-pattern-in-go,54/

import (
	"fmt"
	"time"
)

type barOpts struct {
	total, width               int
	start, end                 string
	complete, head, incomplete string
	formatString               string
	callback                   func()
}

// NewWithFormat creates a new instance of bar.Bar with the given total
// and format and returns a reference to it
func NewWithFormat(t int, f string) *Bar {
	return &Bar{
		progress:     0,
		total:        t,
		width:        20,
		start:        "(",
		complete:     "█",
		head:         "█",
		incomplete:   " ",
		end:          ")",
		closed:       false,
		startedAt:    time.Now(),
		rate:         0,
		formatString: f,
		format:       tokenize(f, nil),
		callback:     noop,
	}
}

// NewWithOpts creates a new instance of bar.Bar with the provided options
// and returns a reference to it
func NewWithOpts(opts ...func(o *barOpts)) *Bar {
	o := &barOpts{
		start:        "(",
		complete:     "█",
		head:         "█",
		incomplete:   " ",
		end:          ")",
		formatString: defaultFormat,
		callback:     noop,
	}

	for _, augment := range opts {
		augment(o)
	}

	if o.width <= 0 {
		panic(fmt.Sprintf("a bar may not have a zero or negative width (received: %d)", o.width))
	}

	return &Bar{
		progress:     0,
		total:        o.total,
		width:        o.width,
		start:        o.start,
		complete:     o.complete,
		head:         o.head,
		incomplete:   o.incomplete,
		end:          o.end,
		closed:       false,
		startedAt:    time.Now(),
		rate:         0,
		formatString: o.formatString,
		format:       tokenize(o.formatString, nil),
		callback:     o.callback,
	}
}

// WithDisplay augments an options constructor by customizing terminal
// output characters
//
// [ xxx >    ]
// | |   | |  |- end
// | |   | |- incomplete
// | |   |- head
// | |- complete
// |- start
func WithDisplay(start, complete, head, incomplete, end string) func(*barOpts) {
	return func(o *barOpts) {
		o.start = start
		o.complete = complete
		o.head = head
		o.incomplete = incomplete
		o.end = end
	}
}

// WithDimensions augments an options constructor by customizing the
// bar's width and total
func WithDimensions(total, width int) func(*barOpts) {
	return func(o *barOpts) {
		o.total = total
		o.width = width
	}
}

// WithFormat augments an options constructor by customizing the bar's
// output format
func WithFormat(f string) func(*barOpts) {
	return func(o *barOpts) {
		o.formatString = f
	}
}

// WithCallback augments an options constructor by setting a callback
func WithCallback(cb func()) func(*barOpts) {
	return func(o *barOpts) {
		o.callback = cb
	}
}
