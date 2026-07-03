<!-- SPDX-License-Identifier: BSD-3-Clause -->
# `go-ruby-prettyprint` library-level benchmark harness

Reproducible, cross-runtime benchmark of the **pure-Go `go-ruby-prettyprint`
library** against the reference Ruby runtimes (MRI, MRI + YJIT, JRuby,
TruffleRuby). It measures the **library primitive** — the Wadler/Lindig layout
engine — through its Go API, isolated from the `rbgo` interpreter, so the numbers
answer: *is the pure-Go implementation as fast as the reference runtime's own
`prettyprint`?*

## Layout

- `go/`                  — self-contained Go driver; `go.mod` pins the published
  library by pseudo-version (no `replace`).
- `ruby/prettyprint.rb`  — the equivalent workload driving MRI's `prettyprint`
  stdlib; `ruby/_harness.rb` is the shared timer.
- `run.sh`               — runs every available runtime and prints one Markdown
  table per sub-benchmark (ns/op + ratio vs MRI).

## Run

```sh
GOWORK=off bash benchmarks/run.sh
```

Environment knobs: `OUTER` (timed passes, default 25), `WARM` (untimed warm-up
passes, default 3), and `RUBY`/`JRUBY`/`TRUFFLERUBY` to select runtime binaries.

## Workload

Four representative documents drive the engine at MRI's default `maxwidth` of 79,
each forcing a different mix of **flat** and **broken** layout decisions over the
core primitives `text`, `breakable`, `group`, `nest` and `fill_breakable`:

- **format-flat**   — a short list `[0, 1, …, 7]` that fits on one line, so every
  breakable stays flat (the no-break fast path).
- **format-nested** — an array of 12 keyed-hash records that overflows and breaks;
  inner groups break or stay flat by their own width (the `break_outmost_groups`
  path on every line).
- **format-fill**   — a 40-word paragraph laid out with `fill_breakable`, wrapping
  across several lines (fill mode / `group_sub`).
- **format-deep**   — 20 levels of nested groups that force indentation and a break
  at every level (nest + deep group nesting).

## Method

Each process runs `WARM` untimed passes (to let the JVM / GraalVM JITs warm up),
then `OUTER` timed passes of a fixed inner loop, timed with a monotonic clock;
the **best** pass is reported as **ns/op**. Interpreter start-up is outside the
timed region. The Go driver and the Ruby script build **identical documents** and
their rendered output is verified **byte-identical to MRI** before timing —

```sh
diff <(cd go && VERIFY=1 GOWORK=off go run .) <(VERIFY=1 ruby ruby/prettyprint.rb)
```

— which is empty (the Wadler engine is byte-exact). Results are published, dated,
in `../docs/performance.md`.

## Notes

- `go/bench` (the compiled driver) is git-ignored; no binary is committed.
- JRuby and TruffleRuby carry JVM / Graal warm-up; the harness warms them before
  timing, but sub-microsecond ops still carry the most relative noise.
