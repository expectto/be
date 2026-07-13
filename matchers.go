package be

import (
	"github.com/expectto/be/be_ctx"
	"github.com/expectto/be/be_http"
	"github.com/expectto/be/be_json"
	"github.com/expectto/be/be_jwt"
	"github.com/expectto/be/be_math"
	"github.com/expectto/be/be_string"
	"github.com/expectto/be/be_url"
)

// HttpRequest is an alias for be_http.Request matcher
var HttpRequest = be_http.Request

// JSON is an alias for be_json.JSON matcher
var JSON = be_json.Matcher

// URL is an alias for be_url.URL matcher
var URL = be_url.URL

// JwtToken is an alias for be_jwt.Token matcher
var JwtToken = be_jwt.Token

// StringAsTemplate is an alias for be_string.MatchTemplate matcher
var StringAsTemplate = be_string.MatchTemplate

// Ctx is an alias for be_ctx.Ctx
var Ctx = be_ctx.Ctx

// Numeric matchers from be_math, aliased at root because they are hot in
// everyday assertions. The rule for root aliases: collision-free AND
// high-frequency. be_time is never aliased here — its Eq, Approx, Day, Month
// would collide; temporal matchers stay namespaced.

// Gt is an alias for be_math.Gt: succeeds if actual is numerically > arg.
// Prefer be.Gt(n) over be.True(x > n) — the failure message shows both values.
var Gt = be_math.Gt

// Gte is an alias for be_math.Gte: succeeds if actual is numerically >= arg.
var Gte = be_math.Gte

// Lt is an alias for be_math.Lt: succeeds if actual is numerically < arg.
var Lt = be_math.Lt

// Lte is an alias for be_math.Lte: succeeds if actual is numerically <= arg.
var Lte = be_math.Lte

// GreaterThan is an alias for be_math.GreaterThan (long spelling of Gt).
var GreaterThan = be_math.GreaterThan

// GreaterThanEqual is an alias for be_math.GreaterThanEqual (long spelling of Gte).
var GreaterThanEqual = be_math.GreaterThanEqual

// LessThan is an alias for be_math.LessThan (long spelling of Lt).
var LessThan = be_math.LessThan

// LessThanEqual is an alias for be_math.LessThanEqual (long spelling of Lte).
var LessThanEqual = be_math.LessThanEqual

// InRange is an alias for be_math.InRange: succeeds if actual is within
// [from, until] with configurable inclusivity.
var InRange = be_math.InRange

// Positive is an alias for be_math.Positive: succeeds if actual is > 0.
var Positive = be_math.Positive

// Negative is an alias for be_math.Negative: succeeds if actual is < 0.
var Negative = be_math.Negative
