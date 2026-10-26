# Shared logscope contract

This file is the required behavior, not illustrative pseudocode. Rust and Go
must agree exactly. Keep the sequential CLI stable through week 12.

## Bytes, records, and malformed lines

Read a file as bytes. Split on LF (byte 10). Remove at most one trailing CR
(byte 13) from each resulting line, including an unterminated final line.
Ignore zero-length lines after this removal; they are not valid or invalid.
An empty file therefore has zero records. A final LF adds no record.

A nonempty line is valid only if it is ASCII and exactly:

```text
LEVEL service duration
```

- Separators are exactly one ASCII space. No leading/trailing spaces, tabs,
  extra fields, comments, or trimming.
- LEVEL is exactly INFO, WARN, or ERROR, case-sensitive.
- service matches `[a-z][a-z0-9-]{0,15}`: 1-16 ASCII bytes.
- duration is 1-5 ASCII decimal digits with value 0 through 60000 inclusive.
  Leading zeros are accepted. Signs, decimal points, and units are not.

Every other nonempty line increments invalid once and contributes nothing else.
Invalid UTF-8, NUL, extra CR bytes, and whitespace-only lines are malformed,
not fatal file errors. Never decode the entire file as UTF-8 before splitting.
Only convert fields to text after ASCII validation. Do not use Unicode word
splitting or language-specific whitespace trimming.

## Aggregation and rendering

Count valid and invalid records, counts per level, and per-service count,
total duration, and maximum duration. Use u64 in Rust and uint64 in Go for
counts and sums. Detect arithmetic overflow and return an analysis error,
never wrap silently. A direct unit test may construct near-limit state to
test overflow; no enormous fixture is needed.

Output ASCII text with LF separators and exactly one final LF:

```text
valid=6 invalid=3
level INFO=3 WARN=1 ERROR=2
service api count=3 total_ms=45 max_ms=20
service auth count=1 total_ms=7 max_ms=7
service worker count=2 total_ms=100 max_ms=100
```

This is the complete expected output for [sample.log](../fixtures/sample.log).
Service lines are sorted by ascending ASCII service name. Each service has one
line; there are no equal-key sorting ambiguities. Include all three levels,
even when zero. Never depend on hash-map iteration order. For
[empty.log](../fixtures/empty.log), output exactly:

```text
valid=0 invalid=0
level INFO=0 WARN=0 ERROR=0
```

For [malformed.log](../fixtures/malformed.log), output exactly:

```text
valid=0 invalid=3
level INFO=0 WARN=0 ERROR=0
```

Expected-output files beside the fixtures contain these exact bytes.

## CLI

The program accepts exactly one positional file path. No flags or stdin mode.
A path containing spaces is one quoted argument. A name beginning with `-`
is still a path. Wrong argument count: exit 2, empty stdout, diagnostic on
stderr starting `error:`. File read, analysis overflow, and stdout write
failures: exit 1, diagnostic starting `error:`. Diagnostic wording after the
prefix is your choice and need not match across languages.

Success, including files containing malformed lines: exit 0, the exact
report on stdout, empty stderr. Build the complete report before writing it;
never emit a success report after a read/analysis failure. An output failure
can leave partial stdout; it cannot be rolled back. Stderr reporting is
best-effort if stderr itself fails. Ordinary invalid input must not panic.
Tests of the pure renderer compare bytes, avoiding shell encoding and newline
conversion. CLI tests invoke a compiled binary, not a build wrapper whose exit
status or build diagnostics may obscure the program's own behavior.

## Required shared case matrix

Test all three fixture outputs; CRLF; no final LF; blank lines; single CR
removal; one-byte service; 16-byte service; rejected 17-byte service; `00001`;
`0`; `60000`; rejected `60001`, six digits, negative values; missing/extra
fields; tabs; doubled spaces; leading/trailing spaces; mixed-case levels;
uppercase services; non-ASCII and invalid UTF-8 bytes; NUL; all levels zero;
multiple services inserted in reverse order; and arithmetic overflow.

No minimum throughput or memory bound is graded. Whole-file reading is
acceptable for these small exercises. Parallel extensions must produce this
same report or return an explicit error, never a silently partial result.
