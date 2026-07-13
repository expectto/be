// Package belint is a go/analysis linter for the be matcher library. It flags
// assertions that hide a raw expression inside be.True()/be.False() — where the
// failure message degenerates to "expected true" — and be.Not(...) compositions
// that have a dedicated matcher, suggesting (and with -fix, applying) the
// matcher spelling from the anti-pattern table in MATCHERS.md.
//
// Point-of-use feedback is the goal: a diagnostic on the offending line is how
// both humans and LLM agents discover a matcher exists without reading docs.
package belint

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/printer"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

const bePkgPath = "github.com/expectto/be"

// Analyzer flags be assertion anti-patterns and suggests matcher spellings.
var Analyzer = &analysis.Analyzer{
	Name:     "belint",
	Doc:      "flags raw expressions wrapped in be.True()/be.False() and be.Not(...) compositions that have a dedicated be matcher",
	URL:      "https://github.com/expectto/be/tree/main/x/belint",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Nodes already covered by an enclosing composite fix (e.g. the
	// be.HaveLength(0) inside be.Not(be.HaveLength(0))): don't double-report.
	covered := map[ast.Node]bool{}

	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		if a, ok := asAssertion(pass, call); ok {
			checkRawActual(pass, a)
		}
		checkComposite(pass, call, covered)
	})
	return nil, nil
}

// assertion is a normalized be assertion call: be.AssertThat/RequireThat(t,
// actual, matcher, ...) or be.Expect/Require(t, actual).To/NotTo/ToNot(matcher).
type assertion struct {
	actual  ast.Expr
	matcher ast.Expr
	// negated is true for NotTo/ToNot; methodSel then points at the method
	// identifier so a fix can rewrite it to To.
	negated   bool
	methodSel *ast.Ident
}

func asAssertion(pass *analysis.Pass, call *ast.CallExpr) (assertion, bool) {
	if name, ok := beFunc(pass, call.Fun); ok {
		if (name == "AssertThat" || name == "RequireThat") && len(call.Args) >= 3 {
			return assertion{actual: call.Args[1], matcher: call.Args[2]}, true
		}
		return assertion{}, false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || len(call.Args) < 1 {
		return assertion{}, false
	}
	var negated bool
	switch sel.Sel.Name {
	case "To":
	case "NotTo", "ToNot":
		negated = true
	default:
		return assertion{}, false
	}
	inner, ok := ast.Unparen(sel.X).(*ast.CallExpr)
	if !ok {
		return assertion{}, false
	}
	if name, ok := beFunc(pass, inner.Fun); ok && (name == "Expect" || name == "Require") && len(inner.Args) >= 2 {
		return assertion{actual: inner.Args[1], matcher: call.Args[0], negated: negated, methodSel: sel.Sel}, true
	}
	return assertion{}, false
}

// beFunc reports whether expr is a selector into the be package (under any
// import alias) and returns the selected function name.
func beFunc(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	sel, ok := ast.Unparen(expr).(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return "", false
	}
	pkgName, ok := pass.TypesInfo.Uses[id].(*types.PkgName)
	if !ok || pkgName.Imported().Path() != bePkgPath {
		return "", false
	}
	return sel.Sel.Name, true
}

// beQual returns the identifier the be package is referenced by at this call
// site (usually "be", but honor aliases in fixes).
func beQual(expr ast.Expr) string {
	if sel, ok := ast.Unparen(expr).(*ast.SelectorExpr); ok {
		if id, ok := sel.X.(*ast.Ident); ok {
			return id.Name
		}
	}
	return "be"
}

// ---- raw expression in actual position -------------------------------------

func checkRawActual(pass *analysis.Pass, a assertion) {
	matcherCall, ok := ast.Unparen(a.matcher).(*ast.CallExpr)
	if !ok {
		return
	}
	name, ok := beFunc(pass, matcherCall.Fun)
	if !ok || (name != "True" && name != "False") || len(matcherCall.Args) != 0 {
		return
	}
	want := name == "True"
	if a.negated {
		want = !want
	}

	qual := beQual(matcherCall.Fun)
	r, ok := rewriteExpr(pass, a.actual, want, qual)
	if !ok {
		return
	}

	msg := fmt.Sprintf("prefer %s over wrapping the raw expression in %s.%s()", r.matcher, qual, name)
	diag := analysis.Diagnostic{Pos: a.actual.Pos(), End: a.actual.End(), Message: msg}
	if r.fixable {
		edits := []analysis.TextEdit{
			{Pos: a.actual.Pos(), End: a.actual.End(), NewText: []byte(r.actual)},
			{Pos: a.matcher.Pos(), End: a.matcher.End(), NewText: []byte(r.matcher)},
		}
		if a.negated {
			// the rewritten matcher already carries the negation
			edits = append(edits, analysis.TextEdit{Pos: a.methodSel.Pos(), End: a.methodSel.End(), NewText: []byte("To")})
		}
		diag.SuggestedFixes = []analysis.SuggestedFix{{
			Message:   fmt.Sprintf("replace with %s", r.matcher),
			TextEdits: edits,
		}}
	}
	pass.Report(diag)
}

// rewrite is the matcher spelling of a raw boolean expression: assert `actual`
// against `matcher`. fixable is false when applying it would need a new import
// (e.g. be_string) or a type name (MatchErrorAs) — those are report-only.
type rewrite struct {
	actual  string
	matcher string
	fixable bool
}

func rewriteExpr(pass *analysis.Pass, expr ast.Expr, want bool, qual string) (rewrite, bool) {
	expr = ast.Unparen(expr)
	switch e := expr.(type) {
	case *ast.UnaryExpr:
		if e.Op == token.NOT {
			return rewriteExpr(pass, e.X, !want, qual)
		}
	case *ast.BinaryExpr:
		return rewriteBinary(pass, e, want, qual)
	case *ast.CallExpr:
		return rewriteCall(pass, e, want, qual)
	}
	return rewrite{}, false
}

func rewriteBinary(pass *analysis.Pass, e *ast.BinaryExpr, want bool, qual string) (rewrite, bool) {
	op := e.Op
	switch op {
	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
	default:
		return rewrite{}, false
	}
	if !want {
		op = negateOp(op)
	}

	lhs, rhs := ast.Unparen(e.X), ast.Unparen(e.Y)
	// Keep the non-constant side as the assertion subject: `0 < len(xs)`
	// becomes `len(xs) > 0`, `nil != err` becomes `err != nil`.
	if (isConstExpr(pass, lhs) || isNilIdent(pass, lhs)) && !isConstExpr(pass, rhs) && !isNilIdent(pass, rhs) {
		lhs, rhs = rhs, lhs
		op = mirrorOp(op)
	}

	// len(x) OP n
	if lenArg, ok := asLenCall(pass, lhs); ok {
		return rewriteLen(pass, lenArg, op, rhs, qual)
	}

	// x == nil / x != nil
	if isNilIdent(pass, rhs) {
		switch op {
		case token.EQL:
			return rewrite{actual: render(pass, lhs), matcher: qual + ".Nil()", fixable: true}, true
		case token.NEQ:
			return rewrite{actual: render(pass, lhs), matcher: qual + ".NotNil()", fixable: true}, true
		}
		return rewrite{}, false
	}

	switch op {
	case token.EQL:
		return rewrite{actual: render(pass, lhs), matcher: fmt.Sprintf("%s.Eq(%s)", qual, render(pass, rhs)), fixable: true}, true
	case token.NEQ:
		return rewrite{actual: render(pass, lhs), matcher: fmt.Sprintf("%s.Ne(%s)", qual, render(pass, rhs)), fixable: true}, true
	default:
		// ordered comparisons: only rewrite numerics (be.Gt is numeric-only)
		if !isNumeric(pass, lhs) && !isNumeric(pass, rhs) {
			return rewrite{}, false
		}
		return rewrite{actual: render(pass, lhs), matcher: fmt.Sprintf("%s.%s(%s)", qual, orderedMatcher(op), render(pass, rhs)), fixable: true}, true
	}
}

func rewriteLen(pass *analysis.Pass, lenArg ast.Expr, op token.Token, rhs ast.Expr, qual string) (rewrite, bool) {
	actual := render(pass, lenArg)
	zero := isZeroConst(pass, rhs)
	switch {
	case op == token.EQL && zero:
		return rewrite{actual: actual, matcher: qual + ".Empty()", fixable: true}, true
	case (op == token.NEQ || op == token.GTR) && zero:
		return rewrite{actual: actual, matcher: qual + ".NotEmpty()", fixable: true}, true
	case op == token.EQL:
		return rewrite{actual: actual, matcher: fmt.Sprintf("%s.HaveLength(%s)", qual, render(pass, rhs)), fixable: true}, true
	case op == token.NEQ:
		return rewrite{actual: actual, matcher: fmt.Sprintf("%s.Not(%s.HaveLength(%s))", qual, qual, render(pass, rhs)), fixable: true}, true
	default:
		return rewrite{actual: actual, matcher: fmt.Sprintf("%s.HaveLength(%s.%s(%s))", qual, qual, orderedMatcher(op), render(pass, rhs)), fixable: true}, true
	}
}

func rewriteCall(pass *analysis.Pass, e *ast.CallExpr, want bool, qual string) (rewrite, bool) {
	fn := typeutil.Callee(pass.TypesInfo, e)
	if fn == nil || fn.Pkg() == nil {
		return rewrite{}, false
	}

	// maybeNot wraps the matcher in be.Not(...) when the raw expression was
	// asserted false.
	maybeNot := func(matcher string) string {
		if want {
			return matcher
		}
		return fmt.Sprintf("%s.Not(%s)", qual, matcher)
	}

	switch fn.Pkg().Path() + "." + fn.Name() {
	case "slices.Contains":
		if len(e.Args) != 2 {
			return rewrite{}, false
		}
		return rewrite{
			actual:  render(pass, e.Args[0]),
			matcher: maybeNot(fmt.Sprintf("%s.ContainElement(%s)", qual, render(pass, e.Args[1]))),
			fixable: true,
		}, true
	case "strings.Contains":
		if len(e.Args) != 2 {
			return rewrite{}, false
		}
		return rewrite{
			actual:  render(pass, e.Args[0]),
			matcher: maybeNot(fmt.Sprintf("%s.ContainSubstring(%s)", qual, render(pass, e.Args[1]))),
			fixable: true,
		}, true
	case "errors.Is":
		if len(e.Args) != 2 {
			return rewrite{}, false
		}
		return rewrite{
			actual:  render(pass, e.Args[0]),
			matcher: maybeNot(fmt.Sprintf("%s.MatchError(%s)", qual, render(pass, e.Args[1]))),
			fixable: true,
		}, true
	case "errors.As":
		if len(e.Args) != 2 {
			return rewrite{}, false
		}
		target := "T"
		if ptr, ok := pass.TypesInfo.TypeOf(e.Args[1]).(*types.Pointer); ok {
			target = types.TypeString(ptr.Elem(), types.RelativeTo(pass.Pkg))
		}
		// report-only: naming the type may require adding an import
		return rewrite{
			actual:  render(pass, e.Args[0]),
			matcher: maybeNot(fmt.Sprintf("%s.MatchErrorAs[%s]()", qual, target)),
		}, true
	case "strings.HasPrefix":
		if len(e.Args) != 2 {
			return rewrite{}, false
		}
		// report-only: be_string may not be imported yet
		return rewrite{
			actual:  render(pass, e.Args[0]),
			matcher: maybeNot(fmt.Sprintf("be_string.HavingPrefix(%s)", render(pass, e.Args[1]))),
		}, true
	case "strings.HasSuffix":
		if len(e.Args) != 2 {
			return rewrite{}, false
		}
		return rewrite{
			actual:  render(pass, e.Args[0]),
			matcher: maybeNot(fmt.Sprintf("be_string.HavingSuffix(%s)", render(pass, e.Args[1]))),
		}, true
	}
	return rewrite{}, false
}

// ---- composites in matcher position ----------------------------------------

func checkComposite(pass *analysis.Pass, call *ast.CallExpr, covered map[ast.Node]bool) {
	if covered[call] {
		return
	}
	name, ok := beFunc(pass, call.Fun)
	if !ok {
		return
	}
	qual := beQual(call.Fun)

	suggest := func(repl string) {
		pass.Report(analysis.Diagnostic{
			Pos:     call.Pos(),
			End:     call.End(),
			Message: fmt.Sprintf("prefer %s over %s", repl, render(pass, call)),
			SuggestedFixes: []analysis.SuggestedFix{{
				Message:   fmt.Sprintf("replace with %s", repl),
				TextEdits: []analysis.TextEdit{{Pos: call.Pos(), End: call.End(), NewText: []byte(repl)}},
			}},
		})
	}

	switch name {
	case "HaveLength":
		if len(call.Args) == 1 && isZeroConst(pass, call.Args[0]) {
			suggest(qual + ".Empty()")
		}
	case "Not":
		if len(call.Args) != 1 {
			return
		}
		inner, ok := ast.Unparen(call.Args[0]).(*ast.CallExpr)
		if !ok {
			return
		}
		innerName, ok := beFunc(pass, inner.Fun)
		if !ok {
			return
		}
		switch innerName {
		case "Nil":
			if len(inner.Args) == 0 {
				suggest(qual + ".NotNil()")
			}
		case "Empty":
			if len(inner.Args) == 0 {
				suggest(qual + ".NotEmpty()")
			}
		case "HaveLength":
			if len(inner.Args) == 1 && isZeroConst(pass, inner.Args[0]) {
				covered[inner] = true
				suggest(qual + ".NotEmpty()")
			}
		case "Eq":
			if len(inner.Args) == 1 && isZeroConst(pass, inner.Args[0]) {
				covered[inner] = true
				suggest(qual + ".NonZero()")
			}
		}
	}
}

// ---- small helpers -----------------------------------------------------------

func negateOp(op token.Token) token.Token {
	switch op {
	case token.EQL:
		return token.NEQ
	case token.NEQ:
		return token.EQL
	case token.LSS:
		return token.GEQ
	case token.LEQ:
		return token.GTR
	case token.GTR:
		return token.LEQ
	case token.GEQ:
		return token.LSS
	}
	return op
}

func mirrorOp(op token.Token) token.Token {
	switch op {
	case token.LSS:
		return token.GTR
	case token.LEQ:
		return token.GEQ
	case token.GTR:
		return token.LSS
	case token.GEQ:
		return token.LEQ
	}
	return op
}

func orderedMatcher(op token.Token) string {
	switch op {
	case token.GTR:
		return "Gt"
	case token.GEQ:
		return "Gte"
	case token.LSS:
		return "Lt"
	case token.LEQ:
		return "Lte"
	}
	return ""
}

func asLenCall(pass *analysis.Pass, expr ast.Expr) (ast.Expr, bool) {
	call, ok := ast.Unparen(expr).(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return nil, false
	}
	id, ok := ast.Unparen(call.Fun).(*ast.Ident)
	if !ok || id.Name != "len" {
		return nil, false
	}
	if b, ok := pass.TypesInfo.Uses[id].(*types.Builtin); !ok || b.Name() != "len" {
		return nil, false
	}
	return call.Args[0], true
}

func isConstExpr(pass *analysis.Pass, expr ast.Expr) bool {
	tv, ok := pass.TypesInfo.Types[expr]
	return ok && tv.Value != nil
}

func isZeroConst(pass *analysis.Pass, expr ast.Expr) bool {
	tv, ok := pass.TypesInfo.Types[ast.Unparen(expr)]
	if !ok || tv.Value == nil {
		return false
	}
	switch tv.Value.Kind() {
	case constant.Int, constant.Float:
		return constant.Sign(tv.Value) == 0
	}
	return false
}

func isNilIdent(pass *analysis.Pass, expr ast.Expr) bool {
	tv, ok := pass.TypesInfo.Types[ast.Unparen(expr)]
	return ok && tv.IsNil()
}

func isNumeric(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}
	b, ok := t.Underlying().(*types.Basic)
	return ok && b.Info()&types.IsNumeric != 0
}

func render(pass *analysis.Pass, node ast.Node) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, pass.Fset, node); err != nil {
		return "<expr>"
	}
	return buf.String()
}
