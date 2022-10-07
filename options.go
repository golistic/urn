// Copyright (c) 2022, Geert JM Vanderkelen

package urn

import "github.com/geertjanvdk/xkit/xutil"

// Option is the functional option type used with New().
type Option func(*urnOptions)

type urnOptions struct {
	resolution      *string
	fragment        *string
	query           *string
	notLowerCaseNSS bool
}

// WithFragment is a functional option setting the f-component of
// the URN; the part introduced by the number sign `#` (same syntax
// as the URI fragment component).
func WithFragment(c string) Option {
	return func(o *urnOptions) {
		o.fragment = xutil.StringPtr(c)
	}
}

// WithQuery is a functional option setting the q-component of
// the URN; the part introduced by `?=` (not `?` alone as is
// the case for a URI).
func WithQuery(c string) Option {
	return func(o *urnOptions) {
		o.query = xutil.StringPtr(c)
	}
}

// WithResolution is a functional option setting the r-component of
// the URN; the part after introduced by `?+`.
func WithResolution(c string) Option {
	return func(o *urnOptions) {
		o.resolution = xutil.StringPtr(c)
	}
}

// WithNotLowerCaseNSS will not lower-case the NSS, but leave it
// as specified to New() or Parse().
// Note that when check whether URNs are equivalent, the NSS is
// considered case-insensitive.
func WithNotLowerCaseNSS() Option {
	return func(o *urnOptions) {
		o.notLowerCaseNSS = true
	}
}
