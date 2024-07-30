package be

import (
	"github.com/expectto/be/be_ctx"
	"github.com/expectto/be/be_http"
	"github.com/expectto/be/be_json"
	"github.com/expectto/be/be_jwt"
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
