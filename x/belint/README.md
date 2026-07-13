# belint

A `go/analysis` linter for the [be](https://github.com/expectto/be) matcher
library. It flags assertions that hide a raw expression inside
`be.True()`/`be.False()` — where the failure message degenerates to
"expected true, got false" — plus `be.Not(...)` compositions that have a
dedicated matcher. Most findings come with a suggested fix, so `-fix` rewrites
them mechanically.

| Flags | Suggests |
|---|---|
| `be.True(x == y)` / `be.False(x == y)` | `be.Eq(y)` / `be.Ne(y)` |
| `be.True(x >= n)` (numeric) | `be.Gte(n)` (and `Gt`/`Lt`/`Lte`) |
| `be.True(err == nil)` / `be.True(err != nil)` | `be.Nil()` / `be.NotNil()` |
| `be.True(len(xs) == 0)` / `be.True(len(xs) > 0)` | `be.Empty()` / `be.NotEmpty()` |
| `be.True(len(xs) >= n)` | `be.HaveLength(be.Gte(n))` |
| `be.True(slices.Contains(xs, v))` | `be.ContainElement(v)` |
| `be.True(strings.Contains(s, q))` | `be.ContainSubstring(q)` |
| `be.True(strings.HasPrefix(s, p))` | `be_string.HavingPrefix(p)` (report-only) |
| `be.True(errors.Is(err, X))` | `be.MatchError(X)` |
| `be.True(errors.As(err, &v))` | `be.MatchErrorAs[V]()` (report-only) |
| `be.Not(be.Nil())` | `be.NotNil()` |
| `be.HaveLength(0)` / `be.Not(be.HaveLength(0))` | `be.Empty()` / `be.NotEmpty()` |
| `be.Not(be.Empty())` | `be.NotEmpty()` |
| `be.Not(be.Eq(0))` | `be.NonZero()` |

All assertion spellings are recognized: `be.AssertThat` / `be.RequireThat` and
`be.Expect(...)` / `be.Require(...)` with `.To` / `.NotTo` / `.ToNot` (a fixed
`.NotTo(be.True())` is rewritten to `.To(...)` with the negation folded into
the matcher). "Report-only" rows have no autofix because applying one would
need a new import or a type name.

## golangci-lint (module plugin)

belint registers with golangci-lint's
[module plugin system](https://golangci-lint.run/plugins/module-plugins/).
Build a custom binary once:

```yaml
# .custom-gcl.yml
version: v2.1.0 # your golangci-lint version
plugins:
  - module: 'github.com/expectto/be/x/belint'
    import: 'github.com/expectto/be/x/belint'
    version: latest
```

```sh
golangci-lint custom # produces ./custom-gcl
```

Then enable it:

```yaml
# .golangci.yml
linters:
  enable:
    - belint
  settings:
    custom:
      belint:
        type: 'module'
        description: 'be matcher anti-patterns'
```

## Standalone (vettool, with -fix)

```sh
go run github.com/expectto/be/x/belint/cmd/belint@latest ./...
go run github.com/expectto/be/x/belint/cmd/belint@latest -fix ./...
```

or via `go vet`:

```sh
go install github.com/expectto/be/x/belint/cmd/belint@latest
go vet -vettool=$(which belint) ./...
```
