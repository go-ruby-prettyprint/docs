// Copyright (c) the go-ruby-prettyprint authors
// SPDX-License-Identifier: BSD-3-Clause
//
// Go driver for the go-ruby-prettyprint library-level benchmark. It builds the
// SAME four representative PrettyPrint documents as ruby/prettyprint.rb and
// drives the pure-Go Wadler/Lindig layout engine at a fixed maxwidth (79, MRI's
// default). Each document exercises the engine's core primitives — text,
// breakable, group, nest and fill — while forcing a different mix of flat and
// broken layout decisions:
//
//   format-flat   — a short list that fits on one line (flat path only)
//   format-nested — an array of hashes that overflows and breaks (group/break)
//   format-fill   — a word paragraph laid out with fill_breakable (fill mode)
//   format-deep   — deeply nested groups that force indentation and breaks
//
// With VERIFY set in the environment the driver prints each rendered document
// (framed by a header line) instead of timing, so run.sh can diff it against
// MRI's PrettyPrint and prove the Go output is byte-identical before any number
// is taken.
package main

import (
	"fmt"
	"os"
	"strconv"

	pp "github.com/go-ruby-prettyprint/prettyprint"
)

// --- document builders (mirrored exactly by ruby/prettyprint.rb) ------------

// emitFlat lays out a short list "[0, 1, ..., 7]" that fits within maxwidth, so
// every breakable stays flat and renders as its separator. Exercises the
// no-break fast path of the engine.
func emitFlat(q *pp.PrettyPrint) {
	q.GroupDefault(func() {
		q.TextString("[")
		for i := 0; i < 8; i++ {
			if i > 0 {
				q.TextString(",")
				q.BreakableString()
			}
			q.TextString(strconv.Itoa(i))
		}
		q.TextString("]")
	})
}

// emitNested lays out an array of 12 hash records, each with three keyed short
// array values. The outer array overflows maxwidth and breaks; inner groups
// break or stay flat depending on their own width, so the document forces the
// break_outmost_groups machinery on every line.
func emitNested(q *pp.PrettyPrint) {
	q.Group(2, "[", 1, "]", 1, func() {
		for i := 0; i < 12; i++ {
			if i > 0 {
				q.TextString(",")
				q.BreakableString()
			}
			q.Group(2, "{", 1, "}", 1, func() {
				for k := 0; k < 3; k++ {
					if k > 0 {
						q.TextString(",")
						q.BreakableString()
					}
					q.TextString(fmt.Sprintf(":key%d=>", k))
					q.Group(1, "[", 1, "]", 1, func() {
						m := i % 4
						for v := 0; v < m; v++ {
							if v > 0 {
								q.TextString(",")
								q.BreakableString()
							}
							q.TextString(strconv.Itoa(i*k + v))
						}
					})
				}
			})
		}
	})
}

// emitFill lays out a 40-word paragraph with fill_breakable between words. Fill
// mode makes the break decision individually at each separator, so the total
// (~150 columns) wraps across several lines and drives the group_sub/fill path.
func emitFill(q *pp.PrettyPrint) {
	q.GroupDefault(func() {
		for i := 0; i < 40; i++ {
			if i > 0 {
				q.FillBreakableString()
			}
			q.TextString(word(i))
		}
	})
}

// emitDeep builds 20 levels of nested groups, each adding indent and a
// breakable, so the whole structure overflows and breaks at every level —
// stressing nest, group nesting and indentation generation.
func emitDeep(q *pp.PrettyPrint) {
	var rec func(d int)
	rec = func(d int) {
		if d == 0 {
			q.TextString("leaf")
			return
		}
		q.Group(2, "(", 1, ")", 1, func() {
			q.TextString(fmt.Sprintf("node%d", d))
			q.BreakableString()
			rec(d - 1)
		})
	}
	rec(20)
}

// word is the deterministic word generator shared with the Ruby script: a stem
// letter repeated (i%5)+2 times followed by the index, e.g. "aa0", "bbb1".
func word(i int) string {
	stem := byte('a' + i%5)
	n := i%5 + 2
	b := make([]byte, 0, n+3)
	for j := 0; j < n; j++ {
		b = append(b, stem)
	}
	b = append(b, []byte(strconv.Itoa(i))...)
	return string(b)
}

func render(fn func(*pp.PrettyPrint)) string {
	return pp.FormatDefault(fn)
}

// --- entry point ------------------------------------------------------------

var docs = []struct {
	label string
	fn    func(*pp.PrettyPrint)
}{
	{"format-flat", emitFlat},
	{"format-nested", emitNested},
	{"format-fill", emitFill},
	{"format-deep", emitDeep},
}

func main() {
	if os.Getenv("VERIFY") != "" {
		for _, d := range docs {
			fmt.Printf("---8<--- %s\n%s\n", d.label, render(d.fn))
		}
		return
	}
	for _, d := range docs {
		fn := d.fn
		bench(d.label, 200, func() { sink = render(fn) })
	}
}
