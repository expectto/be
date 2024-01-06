package be

import (
	"github.com/expectto/be/be_http"
	"github.com/expectto/be/be_jwt"
	"github.com/expectto/be/be_strings"
	"github.com/expectto/be/be_url"
)

// HttpRequest is an alias for be_http.Request matcher
var HttpRequest = be_http.Request

// URL is an alias for be_url.URL matcher
var URL = be_url.URL

// JwtToken is an alias for be_jwt.Token matcher
var JwtToken = be_jwt.Token

// StringAsTemplate is an alias for be_strings.MatchTemplate matcher
var StringAsTemplate = be_strings.MatchTemplate
