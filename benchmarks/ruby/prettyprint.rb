# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
#
# Ruby driver for the go-ruby-prettyprint library-level benchmark. It builds the
# SAME four representative documents as go/main.go and drives the reference
# runtime's own `prettyprint` stdlib (PrettyPrint.format) at maxwidth 79. This is
# the oracle the pure-Go engine is byte-checked against and timed beside.
#
# With VERIFY set in the environment it prints each rendered document (framed by
# the same header line as the Go driver) instead of timing, so run.sh can diff
# the two and prove byte-identical output before any number is taken.
require "prettyprint"
require_relative "_harness"

# --- document builders (mirrored exactly by go/main.go) ---------------------

def emit_flat(q)
  q.group do
    q.text("[")
    8.times do |i|
      if i > 0
        q.text(",")
        q.breakable
      end
      q.text(i.to_s)
    end
    q.text("]")
  end
end

def emit_nested(q)
  q.group(2, "[", "]", 1, 1) do
    12.times do |i|
      if i > 0
        q.text(",")
        q.breakable
      end
      q.group(2, "{", "}", 1, 1) do
        3.times do |k|
          if k > 0
            q.text(",")
            q.breakable
          end
          q.text(":key#{k}=>")
          q.group(1, "[", "]", 1, 1) do
            (i % 4).times do |v|
              if v > 0
                q.text(",")
                q.breakable
              end
              q.text((i * k + v).to_s)
            end
          end
        end
      end
    end
  end
end

def emit_fill(q)
  q.group do
    40.times do |i|
      q.fill_breakable if i > 0
      q.text(word(i))
    end
  end
end

def emit_deep(q)
  rec = lambda do |d|
    if d.zero?
      q.text("leaf")
      next
    end
    q.group(2, "(", ")", 1, 1) do
      q.text("node#{d}")
      q.breakable
      rec.call(d - 1)
    end
  end
  rec.call(20)
end

# word matches go/main.go: stem letter repeated (i%5)+2 times, then the index.
def word(i)
  stem = ("a".ord + i % 5).chr
  (stem * (i % 5 + 2)) + i.to_s
end

def render(&block)
  PrettyPrint.format("".dup, 79, "\n", &block)
end

DOCS = [
  ["format-flat",   method(:emit_flat)],
  ["format-nested", method(:emit_nested)],
  ["format-fill",   method(:emit_fill)],
  ["format-deep",   method(:emit_deep)],
].freeze

if ENV["VERIFY"] && !ENV["VERIFY"].empty?
  DOCS.each do |label, fn|
    printf("---8<--- %s\n%s\n", label, render(&fn))
  end
else
  DOCS.each do |label, fn|
    bench(label, 200) { render(&fn) }
  end
end
