# Roadmap

`go-ruby-prettyprint/prettyprint` is grown **test-first**, each capability
differential-tested against MRI rather than built in isolation. Ruby's `PrettyPrint` layout
engine — the deterministic, interpreter-independent Wadler/Lindig algorithm — is
**complete**.

| Stage | What | Status |
| --- | --- | --- |
| Groups & width-fitting | Groups print flat when they fit and break at their breakables when they overflow `maxwidth` — the exact `break_outmost_groups` / depth-bucketed `GroupQueue` algorithm, so nested groups break outermost-first. | **Done** |
| Breakables | `breakable(sep, width)` line-break hints that emit their separator when the line is not broken, with the `width` argument for multibyte/proportional separators. | **Done** |
| Nesting & indentation | `nest(indent)` and the `group(indent, …)` indent argument, with a pluggable `genspace` block (default `DefaultGenSpace`). | **Done** |
| Open/close text & fill mode | `group(indent, open, close)` bracketing text counted toward the fit decision, and `fill_breakable`, where each break is decided individually. | **Done** |
| Single-line formatter | `singleline_format`, where breakables become their separator text and nothing ever breaks; custom `newline` and `maxwidth`. | **Done** |
| Differential oracle & coverage | A shared corpus run through this package and a generated `PrettyPrint.format` / `singleline_format` script, output asserted identical; deterministic golden + white-box tests alone hold 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches and three OSes. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **The engine, not the inspector.** This module is the layout engine only. The `pp`
  object inspector — walking an object graph and emitting the `text`/`group`/`breakable`
  calls — stays host-side in rbgo. This library is the backend it binds to.
- **No interpreter.** The library implements the deterministic algorithm; it never runs
  arbitrary Ruby. Anything that needs a live binding is the consumer's job — that is why
  `rbgo` binds this module rather than the reverse.
- **Reference is reference Ruby (MRI).** Byte-for-byte conformance targets MRI's
  `prettyprint.rb`, pinned by the differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime; the
  dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
engine / inspector split.
