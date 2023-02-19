package provider

import "errors"

var (
	ErrorExternalApi    = errors.New("external api error")
	ErrorNoSuchProvider = errors.New("not such provider")
)
