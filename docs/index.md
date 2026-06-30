# go-ruby-prettyprint documentation

**Ruby's `PrettyPrint` Wadler/Lindig layout engine in pure Go — MRI byte-exact, no cgo.**

`go-ruby-prettyprint/prettyprint` is a faithful, pure-Go (zero cgo) reimplementation of
Ruby's `prettyprint` standard library — the Wadler/Lindig pretty-printing *engine* that
lays out a stream of text, breakable separators and groups into a width-constrained,
nicely indented document — matching reference Ruby (MRI) byte-for-byte. The module path is
`github.com/go-ruby-prettyprint/prettyprint`.

It is the **layout engine** only: the deterministic group/breakable/indent/width-fitting
algorithm (the buffer, the group stack and the depth-bucketed group queue). The `pp`
object inspector that *uses* it — walking an object graph and emitting
`text`/`group`/`breakable` calls — is the host's job and stays in
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby). This library is the
standalone Go backend `pp` binds to, just like
[go-ruby-yaml](https://github.com/go-ruby-yaml),
[go-ruby-regexp](https://github.com/go-ruby-regexp) and
[go-ruby-marshal](https://github.com/go-ruby-marshal). The dependency runs the other way:
this library has **no dependency on the Ruby runtime**.

!!! success "Status: layout engine complete — MRI byte-exact"
    Faithful port of `PrettyPrint`: **groups** that break outermost-first when they overflow `maxwidth`, **breakables** with the `width` argument, **`nest`** and **`group(indent, …)`** indentation with a pluggable `genspace`, **open/close text**, **fill mode** (`fill_breakable`), the **single-line formatter** (`singleline_format`), and **custom `newline` / `maxwidth`**. Validated by a **differential oracle** against the system `ruby` — output compared byte-for-byte — at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets and three OSes.

## What it is — and isn't

This is the *layout engine* only — the deterministic `break_outmost_groups` /
depth-bucketed `GroupQueue` algorithm. The `pp` object inspector that walks an object graph
and emits the `text`/`group`/`breakable` calls is the host's job and stays in rbgo. This
library is the standalone Go backend `pp` binds to.

## Quick taste

```go
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
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`prettyprint`](https://github.com/go-ruby-prettyprint/prettyprint) | the library — Ruby's `PrettyPrint` layout engine in pure Go |
| [`docs`](https://github.com/go-ruby-prettyprint/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-prettyprint.github.io`](https://github.com/go-ruby-prettyprint/go-ruby-prettyprint.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-prettyprint/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **The engine, not the inspector** — only the deterministic Wadler/Lindig layout
  algorithm; the `pp` object walker stays host-side in rbgo.
- **MRI byte-exact.** Output matches reference Ruby's `prettyprint.rb` exactly,
  validated by a differential oracle against the `ruby` binary.
- **Standalone & reusable.** No dependency on the Ruby runtime — the dependency
  runs the other way.
- **100% test coverage** is the target, enforced as a CI gate, across 6 arches
  and 3 OSes.

## Where to go next

- [Why pure Go](why.md) — why this layout engine is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface, the MRI→Go method map, and worked
  examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-prettyprint/prettyprint](https://github.com/go-ruby-prettyprint/prettyprint).
