# A practical bridge from C++

Read Rust sections in weeks 1-6 and Go sections in weeks 7-11. Snippets in
this guide are illustrative only: they are not required assignment APIs or
complete implementations. Use stable official documentation when a compiler
diagnostic contradicts your intuition.

## Learning sequence and source

This is an independent companion to Microsoft's
[Rust for C/C++ Programmers](https://microsoft.github.io/RustTraining/c-cpp-book/).
Use its C++ comparisons to form hypotheses, then test them with small original
exercises, compiler diagnostics, and the stable standard-library docs.

| Week | C++ starting point | Rust concept to practice |
| --- | --- | --- |
| 1 / ch1-4 | Declarations, expressions, build tools | Bindings, types, blocks, Cargo |
| 2 / ch5-6 | vector/span, string/string_view, variant | Collections, views, structs, exhaustive enums |
| 3 / ch7 | RAII, moves, reference validity | Ownership, Copy/Clone/Drop, lifetime relationships |
| 4 / ch8-9 | Headers, test frameworks, exceptions/expected | Modules, built-in tests, Option/Result and `?` |
| 5 / ch10-12 | Constrained templates, conversions, lambdas/ranges | Traits, generics, From/Into/Default, closures/iterators |
| 6 / ch13 | Thread lifetimes and synchronization | Send/Sync, scoped borrows, explicit joining |

Week 2 previews Option/Result only for tiny stats cases; week 3 constructs
already-typed records without parsing. Week 4 adds byte validation and errors.
Do not build the parser while still learning ownership in week 3.

### Precision matters when comparing languages

Treat high-level source comparisons as motivation, not blanket guarantees:

- A moved-from C++ standard-library value is generally valid but unspecified;
  using it is not automatically undefined behavior. Specific operations have
  their own preconditions. Rust instead rejects use of the moved binding.
- Safe Rust does not rule out leaks, deadlocks, panics, or all concurrency bugs.
  Rc/Arc cycles can leak; Arc is shared ownership, not a lock.
- A borrow must not outlive the referenced storage. An annotation does not
  extend storage duration or imply that all borrows last to a closing brace.
- Copy is opt-in for eligible types; Clone does not universally mean a deep
  allocation. String::clone copies text; Rc::clone shares an allocation.
- Ignoring a Result normally produces a must_use warning, not an unconditional
  compile error. Our checks and error contract require deliberate handling.
- `chunks_exact` retains a remainder; it does not panic merely because the
  length is not divisible. Consult std docs before relying on a sidebar.

The ASCII grammar, checked arithmetic, and explicit-error contract override
context-specific source examples that discard errors or choose default values.

## Mental model

| C++ habit | Rust | Go |
| --- | --- | --- |
| RAII resource owner | Ownership plus Drop | Explicit Close, usually paired with defer |
| Copy or move depends on type/expression | Moves by default for non-Copy values | Assignment copies values, sometimes containing shared references |
| Reference lifetime is a programmer obligation | Borrow checker enforces reference validity | GC tracks reachability; does not prevent data races |
| Exceptions or status objects | Result and `?` | Explicit error return and `if err != nil` |
| Templates and virtual interfaces | Generics, trait bounds, trait objects | Generics and implicitly satisfied interfaces |
| std::vector | Vec<T> | Slice descriptor over a backing array |
| std::map / unordered_map | BTreeMap / HashMap | map, unspecified iteration order |

Neither language is "C++ with different punctuation". Separate allocation,
aliasing, mutation, and lifetime questions before choosing a type.

## Rust: ownership is not merely unique_ptr syntax

```rust
let name = String::from("api");
let moved = name;
// println!("{name}"); // Invalid: ownership moved, unlike copying std::string.
let borrowed: &str = &moved;
println!("{borrowed}");
```

A String owns UTF-8 storage; `&str` borrows a string slice. Vec owns elements;
`&[T]` borrows a contiguous sequence. Passing a String by value transfers
ownership unless it is explicitly cloned. Cloning allocates/copies here; it
is not a universal fix for borrow errors. Small Copy types such as u32 copy
implicitly; Copy types cannot implement Drop.

One mutable borrow or multiple shared borrows can be active at a time, subject
to reborrowing rules. Shared references do not generally permit mutation
without special interior-mutability types, which this course does not need.
Borrow scope is often shortened to the last use, rather than the closing brace.
A vector mutation that reallocates cannot occur while a reference into its
storage remains in use. C++ often allows the equivalent dangling-reference
bug to compile.

```rust
fn first(values: &[u32]) -> Option<&u32> {
    values.first()
}
```

Lifetime elision connects the returned reference to the only input reference.
An explicit `'a` annotation expresses that relationship; it does not extend
storage duration or allocate memory. Returning a reference to a local String
is still invalid. Week 3 borrows typed in-memory fields; week 4's new parser
borrows validated input and stores owned aggregation keys in summaries.

## Rust: values, errors, and polymorphism

Enums carry alternatives, including data: `Option<T>` is Some(T) or None;
`Result<T, E>` is Ok(T) or Err(E). Exhaustive `match` makes missing cases
visible at compile time. `?` propagates an error from a compatible Result
function; it does not throw a C++ exception. `unwrap()` and indexing can panic,
so they are inappropriate substitutes for handling ordinary bad input.
Using `unwrap()` in a test setup whose success is being asserted is different.

Traits describe behavior. A small `T: Trait` bound permits static dispatch;
`dyn Trait` permits runtime dispatch with restrictions on supported methods.
Neither requires a class hierarchy. Do not build a trait abstraction until
there is a real boundary to test.

Week 5 uses `R: Read` for a real boundary. Implementing `From<io::Error>` for
your error type wraps an existing error without failure, enabling conversion
with Into and propagation by `?`. This does not make arbitrary numeric casts
lossless: use TryFrom/TryInto for fallible conversions, not a narrowing `as`.
Default for Summary means the empty report state, not a fallback for bad input.
Rust closures infer capture modes from use. `move` changes how captures enter
the closure; whether it implements Fn/FnMut/FnOnce also depends on its body.

Iterators are lazy; adapters such as map and filter do nothing until consumed.
Prefer a clear loop over an unreadable iterator chain. Integer overflow is not
a portable error policy: use checked arithmetic for the contract, regardless
of debug/release configuration.

Rust Drop is useful for deterministic cleanup, but panics are not an error
reporting protocol. Safe Rust prevents many memory errors, not deadlocks,
logic errors, excessive allocations, or incorrect parsing.

## Go: copied values can share storage

```go
a := []int{1, 2, 3}
b := a
b[0] = 9 // a[0] is now 9: both descriptors reference the same array.
c := append([]int(nil), a...)
c[0] = 4 // independent backing array for these elements.
```

A slice value contains a reference to an array, a length, and a capacity.
Copying the descriptor does not copy its elements. append may reuse storage
or allocate: never infer independence from one observed append. A pointer
receiver can mutate the pointed-to struct; a value receiver gets a copy,
but reference-containing fields may still alias. Go has no Rust-style borrow
checker or general const-qualified reference type. C++ const on a containing
object does not recursively make its pointer targets const either.

Copying a map value shares its underlying map. Reads of a missing key yield
the value type's zero value; use `v, ok := m[k]` to distinguish absence.
A nil map can be read but not assigned into. Nil slices can be appended to.
Go strings are immutable byte sequences, not necessarily valid UTF-8.
`range` over a string decodes runes; indexing gets bytes. Use byte rules for
this course's ASCII contract so Rust and Go agree.

## Go: errors and resources are explicit

```go
f, err := os.Open(path)
if err != nil {
    return err
}
defer f.Close() // Illustrative read-only resource lifetime.
```

Register defer only after a successful open. Deferred calls execute when the
function returns, not at the end of the current loop iteration. Opening many
files in one loop with defer can retain all handles; use a helper per file.
For writes, handle Close errors as well as Write errors. Read-only Close
policy is specified in week 9. os.Exit skips deferred functions, so keep it
at the command boundary after your work function has returned.

Errors are values; wrap with `%w` when preserving error identity matters.
panic is not a normal validation response. An interface is satisfied by a
type's method set without an explicit declaration. Prefer a tiny consumer
interface or a standard one such as io.Reader, not an inheritance tree.
A typed nil pointer inside an interface is not necessarily a nil interface.

## Concurrency without magic

Rust scoped threads may borrow data because the scope joins them before the
borrowed data disappears. Move transfers ownership into a thread; Arc shares
ownership, not automatic synchronization. This course needs scoped threads
and per-worker summaries, not Arc<Mutex<...>> everywhere.

Go goroutines need explicit lifecycle management. Channels communicate values;
they do not make arbitrary shared objects safe. Assign ownership of channel
closure, bound worker counts, and wait for all workers before returning.
Context cancellation is cooperative: check it at specified boundaries.
Do not use sleep as a substitute for synchronization in tests.

Go's race detector observes executed paths, not every possible execution.
Clean race tests complement a design that avoids concurrent map writes.
In both languages, deterministic output requires explicit sorting and a
well-defined reduction, not predictable scheduling.

## Questions worth revisiting

When does this operation copy data rather than a descriptor? Who owns the
resource? Who may mutate it? What ends the lifetime? What is the failure
channel? Where is ordering guaranteed? Answer these with actual code in each
weekly reflection; familiar C++ names alone can conceal different semantics.

## Optional post-course roadmap

Outside the 72-hour graded plan, select a topic when you have a concrete need.
These verified Microsoft chapter pages are reference paths, not new required
assignments or promises that their dependencies belong in this project:

- [ch14 unsafe and FFI](https://microsoft.github.io/RustTraining/c-cpp-book/ch14-unsafe-rust-and-ffi.html):
  raw pointers and safe-wrapper obligations; no unsafe implementation here.
- [ch15 no_std](https://microsoft.github.io/RustTraining/c-cpp-book/ch15-no_std-rust-without-the-standard-library.html):
  core, alloc, and embedded constraints; not a week 6 prerequisite.
- [ch16 C++ translation case studies](https://microsoft.github.io/RustTraining/c-cpp-book/ch16-case-studies.html):
  compare enum dispatch and arena-index ownership with your design.
- [ch17 best practices](https://microsoft.github.io/RustTraining/c-cpp-book/ch17-best-practices.html):
  a review checklist, not a mandate to add external crates.
- [ch18 semantic deep dives](https://microsoft.github.io/RustTraining/c-cpp-book/ch18-cpp-rust-semantic-deep-dives.html):
  revisit casts, preprocessing, and C++ concepts as needed.
- [ch19 macros](https://microsoft.github.io/RustTraining/c-cpp-book/ch19-macros.html):
  macro authoring is deferred; using assert_eq/println/derive earlier is fine.

Rc/RefCell graphs, advanced iterator combinators, type-state builders, and
async are also deferred. The required bounded scoped-thread analyzer does not
need shared mutable graphs, a framework, or a new capstone.
