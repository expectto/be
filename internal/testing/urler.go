package testing

import "net/url"

type Urler interface {
	SetUrl(url *url.URL)
}
