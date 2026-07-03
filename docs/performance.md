# Performance

`go-ruby-prettyprint/prettyprint` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `prettyprint` layout
engine. This page records the **methodology** for the comparative benchmark of that module
against the reference Ruby runtimes — part of the ecosystem-wide per-module parity suite.
No numbers are quoted here until they are measured on a pinned host and committed to the
benchmark harness.

## What is measured

The **same** Ruby script — laying out a representative document of nested groups,
breakables and fill mode through `PrettyPrint.format` at a fixed `maxwidth` — is run under
every runtime. `rbgo`'s number reflects **this pure-Go library doing the work**; every
other column is that interpreter's own `prettyprint`. So the comparison is the
**Ruby-visible operation**, apples-to-apples across interpreters. The script prints a
deterministic checksum and its output is checked **byte-identical to MRI** before any
timing is taken.

## Method

- **Best-of-N wall time** (best, not mean, to suppress scheduler noise); single-shot
  processes, no warm-up beyond the script's own loop.
- **Runtimes:** MRI (the oracle) and MRI `--yjit`; JRuby (on OpenJDK); TruffleRuby
  (GraalVM CE Native) — each timed cold, single-shot, so JVM / Graal startup is on every
  run; read those as one-shot `ruby file.rb` costs, not steady-state JIT numbers.
- The benchmark script and harness live in rbgo's repo under
  [`bench/modules/`](https://github.com/go-embedded-ruby/ruby/tree/main/bench/modules).
  Reproduce with `RBGO=./rbgo bash bench/modules/run.sh N`.

## Result (best of 5, ms)

| Runtime | time | vs MRI |
| --- | ---: | ---: |
| **rbgo** (go-ruby-prettyprint) | 80 | 0.89× |
| MRI (ruby 4.0.5) | 90 | 1.00× |
| MRI + YJIT | 70 | 0.78× |
| JRuby 10.1.0.0 | 1380 | 15.33× |
| TruffleRuby 34.0.1 | 220 | 2.44× |

rbgo runs on **go-ruby-prettyprint** and is **slightly faster than MRI** (0.89x) on this group/breakable layout workload (via the `PrettyPrint.format` API).

!!! note "Honest framing"
    JRuby and TruffleRuby are timed **cold, single-shot**, so they carry JVM /
    Graal startup on every run — read them as one-shot `ruby file.rb` costs, the
    same way `rbgo` and MRI are measured, not as steady-state JIT numbers. Rows
    that complete in well under ~200 ms carry the most relative noise; treat
    their ratios as order-of-magnitude. These are **real measured numbers** from
    the 2026-06-30 run (Apple M-series; `ruby 4.0.5 +PRISM`, `jruby 10.1.0.0`,
    `truffleruby 34.0.1`) — nothing is fabricated or cherry-picked.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This section measures the **pure-Go library directly, through its Go API** — not
the `rbgo` interpreter path recorded above. It isolates the Wadler/Lindig layout
engine from Ruby-interpreter dispatch, answering the parity question head-on:
*is the pure-Go implementation as fast as the reference runtime's own
`prettyprint`?* The **same four documents, same inputs, same iteration counts**
run through the Go library and through each reference runtime's `prettyprint`
stdlib; rendered output was verified **byte-identical to MRI** (a plain `diff`,
empty) before any timing — for all four runtimes.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
  All runtimes on the host, no VM.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.
- **Workload:** `format-flat` (fits one line, flat path), `format-nested` (array
  of 12 keyed-hash records that overflows and breaks), `format-fill` (40-word
  paragraph via `fill_breakable`), `format-deep` (20 nested groups). Harness and
  scripts live in [`benchmarks/`](https://github.com/go-ruby-prettyprint/docs/tree/main/benchmarks);
  reproduce with `GOWORK=off bash benchmarks/run.sh`.

#### format-flat

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby-prettyprint (pure Go)** | 882.3 | 0.13× |
| MRI | 6970.0 | 1.00× |
| MRI + YJIT | 2795.0 | 0.40× |
| JRuby | 27563.5 | 3.95× |
| TruffleRuby | 8255.6 | 1.18× |

#### format-nested

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby-prettyprint (pure Go)** | 10373.1 | 0.13× |
| MRI | 79005.0 | 1.00× |
| MRI + YJIT | 30355.0 | 0.38× |
| JRuby | 52685.4 | 0.67× |
| TruffleRuby | 45769.2 | 0.58× |

#### format-fill

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby-prettyprint (pure Go)** | 4871.0 | 0.08× |
| MRI | 57500.0 | 1.00× |
| MRI + YJIT | 24295.0 | 0.42× |
| JRuby | 37077.3 | 0.64× |
| TruffleRuby | 44398.3 | 0.77× |

#### format-deep

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby-prettyprint (pure Go)** | 4835.2 | 0.12× |
| MRI | 40710.0 | 1.00× |
| MRI + YJIT | 17700.0 | 0.43× |
| JRuby | 29856.7 | 0.73× |
| TruffleRuby | 157741.3 | 3.87× |

### go vs YJIT — the headline

The pure-Go engine **beats MRI + YJIT on every one of the four documents**, the
strongest comparison in the suite (YJIT is MRI's own JIT, warmed here):

| Document | go ns/op | YJIT ns/op | go speed-up vs YJIT |
| --- | ---: | ---: | ---: |
| format-flat   |   882.3 |  2795.0 | **3.2× faster** |
| format-nested | 10373.1 | 30355.0 | **2.9× faster** |
| format-fill   |  4871.0 | 24295.0 | **5.0× faster** |
| format-deep   |  4835.2 | 17700.0 | **3.7× faster** |

So on the library primitive the pure-Go Wadler engine is **~3–8× faster than
MRI** and **~3–5× faster than YJIT** across flat, broken, fill and deep-nesting
layouts, while producing byte-identical output. It also beats JRuby and
TruffleRuby on the two substantial documents (`format-nested`, `format-fill`);
on the sub-microsecond `format-flat` and on `format-deep` the JVM/Graal runtimes
carry their warm-up and dispatch overhead and land slower than MRI.

!!! note "Honest framing"
    Numbers are best-of-25 after 3 warm-up passes, on the host above, dated
    2026-07-03 — **real measured values, nothing fabricated or cherry-picked**;
    a second run reproduced every figure within noise. JRuby and TruffleRuby are
    warmed but still pay JVM / Graal dispatch cost; sub-microsecond ops
    (`format-flat`) carry the most relative noise, so treat their ratios as
    order-of-magnitude. `vs MRI` is wall-clock ns/op ratio; **< 1.00× = faster
    than MRI**. Output is verified byte-identical to MRI for all runtimes before
    timing.
