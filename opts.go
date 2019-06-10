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
	output                     Output
	context                    Context
}

type augment func(*barOpts)

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
		output:       &stdout{},
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
		output:       &stdout{},
	}

	for _, aug := range opts {
		aug(o)
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
		format:       tokenize(o.formatString, o.context.customVerbs()),
		callback:     o.callback,
		output:       o.output,
		context:      o.context,
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
func WithDisplay(start, complete, head, incomplete, end string) augment {
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
func WithDimensions(total, width int) augment {
	return func(o *barOpts) {
		o.total = total
		o.width = width
	}
}

// WithFormat augments an options constructor by customizing the bar's
// output format
func WithFormat(f string) augment {
	return func(o *barOpts) {
		o.formatString = f
	}
}

// WithCallback augments an options constructor by setting a callback
func WithCallback(cb func()) augment {
	return func(o *barOpts) {
		o.callback = cb
	}
}

// WithOutput augments an options constructor by setting the output stream
func WithOutput(out Output) augment {
	return func(o *barOpts) {
		o.output = out
	}
}

// WithContext augments an options constructor by setting the initial values
// for the bar's context
func WithContext(ctx Context) augment {
	return func(o *barOpts) {
		o.context = ctx
	}
}
