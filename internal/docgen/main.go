// Command docgen generates MATCHERS.md — the flat, single-file catalog of every
// matcher across all be packages, grouped by intent, with an "instead of"
// column for the raw idioms each matcher supersedes.
//
// One flat file (rather than per-package READMEs) is deliberate: to look inside
// be_math/README.md you must already know be_math exists. The catalog is the
// discovery surface for humans and LLM agents alike.
//
// Run from the repo root (generate-docs.sh does):
//
//	go run ./internal/docgen > MATCHERS.md
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

// packages to scan: dir -> import alias used in examples.
var packages = []string{
	".",
	"be_math",
	"be_string",
	"be_time",
	"be_struct",
	"be_reflected",
	"be_http",
	"be_url",
	"be_json",
	"be_jwt",
	"be_ctx",
}

// insteadOf is the anti-pattern table: raw idioms each matcher supersedes.
// It is the single source of truth — the condensed copies in README.md and
// doc.go, and the belint rules, follow this mapping.
var insteadOf = map[string]string{
	"be.NotNil":               "`be.Not(be.Nil())`",
	"be.Empty":                "`be.HaveLength(0)`, `be.True(len(xs) == 0)`",
	"be.NotEmpty":             "`be.Not(be.HaveLength(0))`, `be.True(len(xs) > 0)`",
	"be.Zero":                 "`be.Eq(0)` for \"unset\"",
	"be.NonZero":              "`be.Not(be.Eq(0))`, `be.Ne(0)`",
	"be.Eq":                   "`be.True(x == y)`",
	"be.Ne":                   "`be.True(x != y)`, `be.Not(be.Eq(y))`",
	"be.Gt":                   "`be.True(x > n)`",
	"be.Gte":                  "`be.True(x >= n)`",
	"be.Lt":                   "`be.True(x < n)`",
	"be.Lte":                  "`be.True(x <= n)`",
	"be.HaveLength":           "`be.True(len(xs) >= n)`",
	"be.ContainElement":       "`be.True(slices.Contains(xs, v))`",
	"be.ContainSubstring":     "`be.True(strings.Contains(s, q))`",
	"be.HaveKey":              "`_, ok := m[k]` + `be.True(ok)`",
	"be.MatchError":           "`be.True(errors.Is(err, X))`",
	"be.MatchErrorAs":         "`var v E` + `be.True(errors.As(err, &v))` — only when `v` is unused afterward (the matcher does not bind it)",
	"be.HaveField":            "`be.Eq(x.Field)` on a projected value",
	"be.NoError":              "`if err != nil { t.Fatal(err) }`",
	"be_string.HavingPrefix":  "`be.True(strings.HasPrefix(s, p))`",
	"be_string.HavingSuffix":  "`be.True(strings.HasSuffix(s, x))`",
	"be_time.SameExactSecond": "`be.True(t1.Equal(t2))`",
	"be_time.Approx":          "`be.True(d < time.Second)` on a time diff",
}

// section assigns entries to intent groups. A section claims an entry either by
// explicit qualified names or by whole package. First match wins; unclaimed
// entries land in the trailing "Other" section so nothing silently disappears.
type section struct {
	Title string
	Names []string // qualified, e.g. "be.Eq"
	Pkgs  []string // whole packages, e.g. "be_math"
}

var sections = []section{
	{
		Title: "Assertions & shortcuts",
		Names: []string{
			"be.Expect", "be.Require", "be.AssertThat", "be.RequireThat",
			"be.NoError", "be.Error", "be.ErrorIs",
		},
	},
	{
		Title: "Async",
		Names: []string{"be.Eventually", "be.Consistently"},
	},
	{
		Title: "Equality & identity",
		Names: []string{
			"be.Eq", "be.Ne", "be.Identical", "be.NotIdentical",
			"be.Zero", "be.NonZero",
		},
	},
	{
		Title: "Errors",
		Names: []string{"be.Succeed", "be.HaveOccurred", "be.MatchError", "be.MatchErrorAs"},
	},
	{
		Title: "Booleans, nil & panics",
		Names: []string{"be.True", "be.False", "be.Nil", "be.NotNil", "be.Panic", "be.NotPanic"},
	},
	{
		Title: "Collections & length",
		Names: []string{
			"be.ContainElement", "be.ContainElements", "be.HaveKey", "be.HaveKeyWithValue",
			"be.Empty", "be.NotEmpty", "be.HaveLength",
			"be.Dive", "be.DiveAny", "be.DiveFirst", "be.DiveNth",
		},
	},
	{
		Title: "Structs",
		Names: []string{"be.HaveField", "be.HaveFields"},
		Pkgs:  []string{"be_struct"},
	},
	{
		Title: "Numbers",
		Names: []string{
			"be.Gt", "be.Gte", "be.Lt", "be.Lte",
			"be.GreaterThan", "be.GreaterThanEqual", "be.LessThan", "be.LessThanEqual",
			"be.InRange", "be.Positive", "be.Negative",
		},
		Pkgs: []string{"be_math"},
	},
	{
		Title: "Strings",
		Names: []string{"be.ContainSubstring", "be.StringAsTemplate"},
		Pkgs:  []string{"be_string"},
	},
	{
		Title: "Time",
		Pkgs:  []string{"be_time"},
	},
	{
		Title: "Composition & control",
		Names: []string{"be.All", "be.Any", "be.Not", "be.Always", "be.Never", "be.Via"},
	},
	{
		Title: "Types & kinds",
		Pkgs:  []string{"be_reflected"},
	},
	{
		Title: "HTTP, URL, JSON, JWT & context",
		Names: []string{"be.HttpRequest", "be.URL", "be.JSON", "be.JwtToken", "be.Ctx"},
		Pkgs:  []string{"be_http", "be_url", "be_json", "be_jwt", "be_ctx"},
	},
}

type entry struct {
	Qualified string // be_math.Approx
	Signature string // Approx(compareTo, threshold any)
	Synopsis  string
}

func main() {
	entries, order, err := collect()
	if err != nil {
		fmt.Fprintln(os.Stderr, "docgen:", err)
		os.Exit(1)
	}

	var b strings.Builder
	b.WriteString("# be — matcher catalog\n\n")
	b.WriteString("<!-- Code generated by internal/docgen (via generate-docs.sh). DO NOT EDIT. -->\n\n")
	b.WriteString("Every assertion helper and matcher across all `be` packages, grouped by\n")
	b.WriteString("intent. The **Instead of** column lists the raw idiom or anti-pattern the\n")
	b.WriteString("matcher supersedes — prefer the matcher: its failure message shows the\n")
	b.WriteString("values involved, `be.True(<expr>)` only reports \"expected true\".\n\n")
	b.WriteString("Matcher arguments may be raw values **or other matchers** (Be/Gomega/Gomock),\n")
	b.WriteString(
		"so matchers compose: `be.HaveLength(be.Gte(3))`, `be.ContainElement(be.HaveField(\"Name\", \"X\"))`.\n\n",
	)

	claimed := map[string]bool{}
	for _, sec := range sections {
		var rows []entry
		for _, name := range sec.Names {
			if e, ok := entries[name]; ok && !claimed[name] {
				rows = append(rows, e)
				claimed[name] = true
			}
		}
		for _, pkg := range sec.Pkgs {
			for _, q := range order {
				if strings.HasPrefix(q, pkg+".") && !claimed[q] {
					rows = append(rows, entries[q])
					claimed[q] = true
				}
			}
		}
		writeSection(&b, sec.Title, rows)
	}

	// Anything unclaimed still gets listed — a new matcher must never be invisible.
	var rest []entry
	for _, q := range order {
		if !claimed[q] {
			rest = append(rest, entries[q])
		}
	}
	writeSection(&b, "Other", rest)

	if _, err := os.Stdout.WriteString(b.String()); err != nil {
		fmt.Fprintln(os.Stderr, "docgen:", err)
		os.Exit(1)
	}
}

func writeSection(b *strings.Builder, title string, rows []entry) {
	if len(rows) == 0 {
		return
	}
	fmt.Fprintf(b, "## %s\n\n", title)
	b.WriteString("| Matcher | What it does | Instead of |\n|---|---|---|\n")
	for _, e := range rows {
		fmt.Fprintf(b, "| `%s` | %s | %s |\n", e.Signature, sanitizeCell(e.Synopsis), insteadOf[e.Qualified])
	}
	b.WriteString("\n")
}

// sanitizeCell keeps a doc synopsis table-safe.
func sanitizeCell(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.TrimSpace(s)
	return strings.TrimSuffix(s, ":") // doc comments often lead into an example with ":"
}

func collect() (map[string]entry, []string, error) {
	entries := map[string]entry{}
	var order []string

	fset := token.NewFileSet()
	for _, dir := range packages {
		pkgName := "be"
		if dir != "." {
			pkgName = dir
		}

		files, err := parsePackageFiles(fset, dir)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", dir, err)
		}
		d, err := doc.NewFromFiles(fset, files, "github.com/expectto/be/"+dir)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", dir, err)
		}

		for _, f := range d.Funcs {
			if !ast.IsExported(f.Name) || !isCatalogFunc(f.Decl.Type, fset) {
				continue
			}
			q := pkgName + "." + f.Name
			entries[q] = entry{
				Qualified: q,
				Signature: renderSignature(pkgName, f.Decl, fset),
				Synopsis:  d.Synopsis(f.Doc),
			}
			order = append(order, q)
		}

		// Matcher aliases declared as vars (e.g. root's `var Gt = be_math.Gt`).
		for _, v := range d.Vars {
			for _, name := range v.Names {
				if !ast.IsExported(name) || name == "V" {
					continue
				}
				q := pkgName + "." + name
				entries[q] = entry{
					Qualified: q,
					Signature: pkgName + "." + name + "(...)",
					Synopsis:  d.Synopsis(v.Doc),
				}
				order = append(order, q)
			}
		}
	}
	return entries, order, nil
}

func parsePackageFiles(fset *token.FileSet, dir string) ([]*ast.File, error) {
	des, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []*ast.File
	for _, de := range des {
		name := de.Name()
		if de.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, dir+"/"+name, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

// isCatalogFunc reports whether a function belongs in the catalog: it returns a
// matcher, or it is an assertion helper driven by TestingT.
func isCatalogFunc(ft *ast.FuncType, fset *token.FileSet) bool {
	if ft.Results != nil && len(ft.Results.List) > 0 {
		if strings.Contains(typeString(ft.Results.List[len(ft.Results.List)-1].Type, fset), "BeMatcher") {
			return true
		}
	}
	if ft.Params != nil && len(ft.Params.List) > 0 {
		if strings.Contains(typeString(ft.Params.List[0].Type, fset), "TestingT") {
			return true
		}
	}
	return false
}

func renderSignature(pkgName string, decl *ast.FuncDecl, fset *token.FileSet) string {
	var b strings.Builder
	b.WriteString(pkgName + "." + decl.Name.Name)
	if decl.Type.TypeParams != nil {
		b.WriteString("[")
		for i, p := range decl.Type.TypeParams.List {
			if i > 0 {
				b.WriteString(", ")
			}
			for j, n := range p.Names {
				if j > 0 {
					b.WriteString(", ")
				}
				b.WriteString(n.Name)
			}
			b.WriteString(" " + typeString(p.Type, fset))
		}
		b.WriteString("]")
	}
	b.WriteString("(")
	if decl.Type.Params != nil {
		for i, p := range decl.Type.Params.List {
			if i > 0 {
				b.WriteString(", ")
			}
			for j, n := range p.Names {
				if j > 0 {
					b.WriteString(", ")
				}
				b.WriteString(n.Name)
			}
			if len(p.Names) > 0 {
				b.WriteString(" ")
			}
			b.WriteString(typeString(p.Type, fset))
		}
	}
	b.WriteString(")")
	return b.String()
}

func typeString(t ast.Expr, fset *token.FileSet) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, t); err != nil {
		return "?"
	}
	return buf.String()
}
