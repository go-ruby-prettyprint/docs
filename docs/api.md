# Usage & API

The public API lives at the module root (`github.com/go-ruby-prettyprint/prettyprint`). It
is **Ruby-shaped but Go-idiomatic**: the methods mirror Ruby's `PrettyPrint`
(`text`/`breakable`/`group`/`nest`/`fill_breakable`), while the surface follows Go
conventions — value types, callbacks for blocks, no global state.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-prettyprint/prettyprint`, bound into `rbgo` as the layout backend; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-prettyprint/prettyprint
```

## Worked example

```go
package main

import (
	"fmt"

	"github.com/go-ruby-prettyprint/prettyprint"
)

func main() {
	// PrettyPrint.format(''.dup, 10) { |q| ... }
	out := prettyprint.Format(10, "\n", nil, func(q *prettyprint.PrettyPrint) {
		q.Group(2, "[", 1, "]", 1, func() {
			q.TextString("111")
			q.BreakableString()
			q.TextString("222")
			q.BreakableString()
			q.TextString("333")
		})
	})
	fmt.Printf("%q\n", out) // "[111\n  222\n  333]"

	// The single-line formatter never breaks: breakables become their separator.
	flat := prettyprint.SingleLineFormat(func(q *prettyprint.SingleLine) {
		q.GroupDefault(func() {
			q.TextString("[")
			q.Breakable("", 0)
			q.TextString("1")
			q.BreakableString()
			q.TextString("2")
			q.TextString("]")
		})
	})
	fmt.Printf("%q\n", flat) // "[1 2]"
}
```

## API map (MRI → Go)

| Ruby (`PrettyPrint`)                       | Go                                                              |
| ------------------------------------------ | -------------------------------------------------------------- |
| `PrettyPrint.new(out, maxwidth, nl, &gs)`  | `New(maxwidth, newline, genspace)` / `NewDefault()`            |
| `PrettyPrint.format(...)`                  | `Format(maxwidth, newline, genspace, fn)` / `FormatDefault`    |
| `PrettyPrint.singleline_format(...)`       | `SingleLineFormat(fn)`                                          |
| `#text(obj, width)`                        | `Text(obj, width)` / `TextString(obj)`                         |
| `#breakable(sep, width)`                   | `Breakable(sep, width)` / `BreakableString()`                  |
| `#group(indent, open, close, ow, cw)`      | `Group(indent, open, ow, close, cw, fn)` / `GroupDefault(fn)`  |
| `#group_sub`                               | `GroupSub(...)` / `NestedGroup(fn)`                            |
| `#nest(indent)`                            | `Nest(indent, fn)`                                             |
| `#fill_breakable(sep, width)`              | `FillBreakable(sep, width)` / `FillBreakableString()`          |
| `#current_group`                           | `CurrentGroup()`                                               |
| `#break_outmost_groups`                    | `BreakOutmostGroups()`                                         |
| `#flush`                                   | `Flush()`                                                      |
| `PrettyPrint::VERSION`                     | `VERSION`                                                       |

The default `genspace` is `DefaultGenSpace` (`n` ASCII spaces); pass `nil` to `New` /
`Format` to use it.

## Capabilities

- **Groups** that print flat when they fit and break at their breakables when they
  overflow `maxwidth` — the exact `break_outmost_groups` / depth-bucketed `GroupQueue`
  algorithm MRI uses, so nested groups break outermost-first.
- **Breakables** — `breakable(sep, width)` line-break hints that emit their separator when
  the line is not broken, with the `width` argument for multibyte or proportional
  separators.
- **Nesting / indentation** — `nest(indent)` and the `group(indent, …)` indent argument,
  with a pluggable `genspace` block for the indentation string.
- **Open/close text** — `group(indent, open, close)` wrapping the block in bracketing text
  counted toward the fit decision.
- **Fill mode** — `fill_breakable`, where each break is decided individually.
- **The single-line formatter** — `singleline_format`, where breakables become their
  separator text and nothing ever breaks.
- **Custom `newline`** and **custom `maxwidth`**, matching MRI's `PrettyPrint.new` /
  `PrettyPrint.format` signatures.

## MRI conformance

Correctness is defined by reference Ruby. A **differential oracle** runs a shared corpus
through both this package and a generated `PrettyPrint.format` / `singleline_format` script
and asserts the output is **identical**, not approximated from memory. It binmodes
stdout/stdin (so Windows text-mode never rewrites the bytes), skips when `ruby` is absent,
and gates on `RUBY_VERSION >= "4.0"`. The deterministic golden + white-box tests alone hold
the 100% gate on every platform, including Windows and the qemu cross-arch lanes.

## Relationship to Ruby

`go-ruby-prettyprint/prettyprint` is **standalone and reusable**, and is the layout backend
bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` — the
same way [go-ruby-yaml](https://github.com/go-ruby-yaml) and
[go-ruby-marshal](https://github.com/go-ruby-marshal) are bound. The `pp` object inspector
that drives it stays host-side; the dependency runs the other way, with no dependency on
the Ruby runtime.
