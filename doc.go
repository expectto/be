// Package be provides fluent, composable test matchers for Go with a native,
// dependency-free assertion runner. Assertions read as an English sentence —
// "expect ... to be ..." — and every matcher argument can be a raw value or
// another matcher (Be/Gomega/Gomock), so matchers compose:
//
//	be.Expect(t, resp.Items).To(be.HaveLength(be.Gte(3)))
//
// # Assertion spellings
//
//	be.Expect(t, actual).To(matcher)           // soft: Errorf, test continues
//	be.Require(t, actual).To(matcher)          // hard: Fatalf, test stops
//	be.AssertThat(t, actual, matcher)          // flat soft spelling
//	be.RequireThat(t, actual, matcher)         // flat hard spelling
//	be.NoError(t, err)                         // hard error shortcuts:
//	be.Error(t, err)                           //   the testify require trio
//	be.ErrorIs(t, err, target)
//	be.Eventually(t, poll, matcher)            // async: poll until it matches
//
// All matchers also work inside gomega (Expect(x).To(be.Eq(y))) and as gomock
// argument matchers.
//
// # Say it with a matcher, not with be.True
//
// Wrapping a raw expression in be.True hides everything from the failure
// message. There is a matcher for almost every idiom:
//
//	be.True(x == y)                  -> be.Eq(y)
//	be.True(x >= n)                  -> be.Gte(n)
//	be.Not(be.Nil())                 -> be.NotNil()
//	be.HaveLength(0)                 -> be.Empty()
//	be.Not(be.HaveLength(0))         -> be.NotEmpty()
//	be.Not(be.Eq(0)), be.Ne(0)       -> be.NonZero()
//	be.True(len(xs) >= n)            -> be.HaveLength(be.Gte(n))
//	be.True(slices.Contains(xs, v))  -> be.ContainElement(v)
//	be.True(strings.Contains(s, q))  -> be.ContainSubstring(q)
//	be.True(strings.HasPrefix(s, p)) -> be_string.HavingPrefix(p)
//	_, ok := m[k]; be.True(ok)       -> be.HaveKey(k)
//	be.True(errors.Is(err, target))  -> be.MatchError(target)
//	var e E; be.True(errors.As(err, &e)) -> be.MatchErrorAs[E]() (if e unused after)
//	be.True(t1.Equal(t2))            -> be_time.SameExactSecond(t2)
//
// The full catalog of matchers across all packages lives in MATCHERS.md at the
// repository root: https://github.com/expectto/be/blob/main/MATCHERS.md
//
// # Subpackages
//
// The root package covers everyday matchers (equality, nil, errors, booleans,
// collections, lengths, structs) plus root aliases for hot numeric matchers.
// Specialized matchers live in subpackages:
//
//   - be_math: numbers (Approx, Odd, Even, DivisibleBy, ...)
//   - be_string: strings (HavingPrefix, MatchTemplate, MatchWildcard, ...)
//   - be_time: time.Time (SameExactSecond, Approx, LaterThan, ...)
//   - be_struct: typed struct fields (HavingField[T])
//   - be_reflected: kind/type assertions (AsNumericString, AsKind, ...)
//   - be_http, be_url, be_json, be_jwt, be_ctx: HTTP requests, URLs, JSON,
//     JWT tokens and contexts
//
// Temporal matchers are never aliased at root (their Eq, Approx, Day would
// collide) — always reach for be_time explicitly.
package be
