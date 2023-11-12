package be

import (
	"github.com/expectto/be/be_http"
	"github.com/expectto/be/be_url"
)

// HttpRequest is a type alias for be_http.Request
// So be.HttpRequest can be used instead of be_http.Request
var HttpRequest = be_http.Request

// URL is a type alias for be_url.URL
var URL = be_url.URL
