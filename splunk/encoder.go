package splunk

import (
	"net/url"
	"github.com/gorilla/schema"
)

var (
	encoder = schema.NewEncoder()
)

func encode(i interface{}) (r url.Values, e error) {
	r = url.Values{}
	e = encoder.Encode(i, r)
	return
}