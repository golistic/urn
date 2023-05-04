// Copyright (c) 2022, Geert JM Vanderkelen

package urn

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const (
	regexNID       = `[0-9a-z][0-9a-z-]{0,30}[0-9a-z]`
	regexNSS       = `[0-9a-z-._~*+=%$&@'()!,:;/]+`
	regexComponent = `[0-9a-z-._~*+=%$&@'()!,:;/]*`
)

var reNID = regexp.MustCompile(`^(?i)` + regexNID + `$`)
var reNSS = regexp.MustCompile(`^(?i)` + regexNSS + `$`)
var reComponent = regexp.MustCompile(`^(?i)` + regexComponent + `$`)
var reURN = regexp.MustCompile(`^(?i)urn:(` + regexNID + `):(` + regexNSS + `)` +
	`(?:(\?+` + regexComponent + `))?` +
	`(?:(\?=` + regexComponent + `))?` +
	`(?:(#` + regexComponent + `))?$`)
var reNormPerEnc = regexp.MustCompile(`(%[0-9a-f]{2})`)

// URN is the representation of a URN as defined by RFC 8141.
//
// The syntax of a URN is: `urn:<NID>:<NSS>#fragment`.
// NID stands for Namespace Identifier, and NSS for Namespace Specific String.
//
// See RFC 8141 https://tools.ietf.org/html/rfc8141
type URN struct {
	NID        string
	NSS        string
	rComponent string
	qComponent string
	fComponent string

	// Original is only set when parsing a URN from a string using Parse.
	Original string
}

// New returns a new instance of URN with nid as Namespace Identifier and nss
// as Namespace Specific String. The r-, q-, and f-components can be set through
// their respective functional options WithResolution, WithQuery, and WithFragment.
// The NSS is lower cased according the RFC. If this is not wanted, use the
// option WithNotLowerCaseNSS.
func New(nid, nss string, options ...Option) (*URN, error) {
	if !reNID.MatchString(nid) {
		return nil, fmt.Errorf("invalid namespace identifier (NID)")
	}

	if !reNSS.MatchString(nss) {
		return nil, fmt.Errorf("invalid namespace specific string (NSS)")
	}

	var opts urnOptions
	for _, o := range options {
		o(&opts)
	}

	urn := &URN{
		NSS: nss,
	}

	if opts.notLowerCaseNSS {
		urn.NID = nid
	} else {
		urn.NID = strings.ToLower(nid)
	}

	if opts.resolution != nil {
		if err := urn.SetRComponent(*opts.resolution); err != nil {
			return nil, err
		}
	}

	if opts.query != nil {
		if err := urn.SetQComponent(*opts.query); err != nil {
			return nil, err
		}
	}

	if opts.fragment != nil {
		if err := urn.SetFComponent(*opts.fragment); err != nil {
			return nil, err
		}
	}

	return urn, nil
}

// MarshalJSON returns the JSON encoding of u.
func (u *URN) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *URN) SetFComponent(c string) error {
	if !IsComponent(c) {
		return fmt.Errorf("invalid f-component")
	}
	u.fComponent = c
	return nil
}

func (u *URN) FComponent() string {
	return u.fComponent
}

func (u *URN) SetQComponent(c string) error {
	if !IsComponent(c) {
		return fmt.Errorf("invalid q-component")
	}
	u.qComponent = c
	return nil
}

func (u *URN) QComponent() string {
	return u.qComponent
}

func (u *URN) SetRComponent(c string) error {
	if !IsComponent(c) {
		return fmt.Errorf("invalid r-component")
	}
	u.rComponent = c
	return nil
}

func (u *URN) RComponent() string {
	return u.rComponent
}

// UnmarshalJSON parses the JSON-encoded data and stores the result in u.
//
// Note that the NID and NSS are being lower-cased; optional components not.
// The Original field of u will contain the original.
func (u *URN) UnmarshalJSON(data []byte) error {
	if data == nil {
		// nothing to do
		return nil
	}

	urn, err := Parse(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}

	if urn.IsZero() {
		// we got nothing; nothing to assign
		return nil
	}

	*u = *urn
	return nil
}

// Equal reports whether o and u represent the same URN.
func (u *URN) Equal(o *URN) bool {
	if u == nil || o == nil {
		panic("cannot check equivalent URN objects when both or either are nil")
	}

	uNSS := u.NSS
	oNSS := o.NSS

	// percent-encoded characters are considered case-insensitive, where the rest of the NSS is not
	if strings.ContainsRune(uNSS, '%') {
		uNSS = reNormPerEnc.ReplaceAllStringFunc(uNSS, strings.ToUpper)
	}

	if strings.ContainsRune(oNSS, '%') {
		oNSS = reNormPerEnc.ReplaceAllStringFunc(oNSS, strings.ToUpper)
	}

	return strings.ToLower(u.NID) == strings.ToLower(o.NID) && uNSS == oNSS
}

// IsZero reports whether u has NSS and NID set.
func (u *URN) IsZero() bool {
	return u == nil || !(len(u.NID) > 0 && len(u.NSS) > 0)
}

// Parse tries to parse s into a URN object.
func Parse(s string, options ...Option) (*URN, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}

	m := reURN.FindStringSubmatch(s)
	if m == nil || len(m) < 2 {
		return nil, fmt.Errorf("invalid")
	}

	var opts urnOptions
	for _, o := range options {
		o(&opts)
	}

	if opts.query != nil || opts.resolution != nil || opts.fragment != nil {
		panic("cannot use WithQuery/WithResolution/WithFragment with Parse")
	}

	u, err := New(m[1], m[2], options...)
	if err != nil {
		return nil, err
	}

	for i := 3; i < len(m); i++ {
		if len(m[i]) < 2 {
			continue
		}

		switch {
		case m[i][0:2] == "?+":
			u.rComponent = m[i][2:len(m[i])]
		case m[i][0:2] == "?=":
			u.qComponent = m[i][2:len(m[i])]
		case m[i][0] == '#':
			u.fComponent = m[i][1:len(m[i])]
		}
	}

	return u, nil
}

func mustParseURN(s string, options ...Option) *URN {
	u, err := Parse(s, options...)
	if err != nil {
		panic(err)
	}
	return u
}

// String returns the string representation of u.
func (u *URN) String() string {
	if u.NID == "" || u.NSS == "" {
		return ""
	}

	n := "urn:" + u.NID + ":" + u.NSS

	if u.rComponent != "" {
		n += "?+" + u.rComponent
	}

	if u.qComponent != "" {
		n += "?=" + u.qComponent
	}

	if u.fComponent != "" {
		n += "#" + u.fComponent
	}

	return n
}

// Validates verifies whether s can be parsed as URN.
func Validates(s string) bool {
	u, err := Parse(s)
	if err != nil {
		return false
	}

	return !u.IsZero()
}

// IsComponent verifies whether s is a valid URN fragment.
func IsComponent(s string) bool {
	return reComponent.MatchString(s)
}
