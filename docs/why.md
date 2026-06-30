# Why pure Go

`go-ruby-prettyprint/prettyprint` reimplements Ruby's `prettyprint` engine in **pure Go,
with cgo disabled**. The layout algorithm it covers is **deterministic and
interpreter-independent**: given a stream of `text` / `breakable` / `group` calls and a
`maxwidth`, the laid-out document is a pure function of those inputs — no live binding, no
evaluation of arbitrary Ruby. That is exactly the part that can — and should — live as a
standalone Go library, separate from the interpreter.

## What it is — and isn't

This is the **layout engine** only — the deterministic group/breakable/indent/width-fitting
algorithm (the buffer, the group stack and the depth-bucketed group queue). The `pp`
object inspector that *uses* it — walking an object graph and emitting
`text`/`group`/`breakable` calls — is the host's job and stays in
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby). This library is the
standalone Go backend `pp` binds to.

## Extracted from rbgo, reusable by anyone

It is the layout backend bound into rbgo, but is a **standalone, reusable library** so
that:

- any Go program can import `github.com/go-ruby-prettyprint/prettyprint` directly, with no
  Ruby runtime;
- the dependency runs the *other* way — `rbgo` binds this module as a native module (the
  same pattern as [go-ruby-yaml](https://github.com/go-ruby-yaml),
  [go-ruby-regexp](https://github.com/go-ruby-regexp) and
  [go-ruby-marshal](https://github.com/go-ruby-marshal)), rather than this module
  depending on the interpreter;
- the behaviour is pinned by a **differential oracle** against the system `ruby`,
  independent of any one consumer.

## Why pure Go matters here

Because the library is CGO-free and dependency-free, it:

- cross-compiles to every Go target with no C toolchain, and links into a single static
  binary;
- has **no dependency on the Ruby runtime** — the dependency runs the other way;
- can be differentially tested against the `ruby` binary wherever one is on `PATH`, while
  the cross-arch and Windows lanes (where `ruby` is absent) still validate the library via
  the deterministic golden + white-box tests.

See [Usage & API](api.md) for the surface and [Roadmap](roadmap.md) for what is in scope.
