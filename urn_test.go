// Copyright (c) 2022, Geert JM Vanderkelen

package urn

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestParseURN(t *testing.T) {
	t.Run("valid URN", func(t *testing.T) {

		validCases := map[string]struct {
			urn         string
			exp         string
			expNID      string
			expNSS      string
			expFragment string
		}{
			"isbn": {
				urn:    "urn:isbn:978-0135800911",
				exp:    "urn:isbn:978-0135800911",
				expNID: "isbn",
				expNSS: "978-0135800911",
			},

			"isbn lower-cased NID": {
				urn:    "UrN:IsBn:978-0135800911",
				exp:    "urn:isbn:978-0135800911",
				expNID: "isbn",
				expNSS: "978-0135800911",
			},
			"isbn with f-component": {
				urn:         "UrN:IsBn:978-0135800911#Page5",
				exp:         "urn:isbn:978-0135800911#Page5",
				expNID:      "isbn",
				expNSS:      "978-0135800911",
				expFragment: "Page5",
			},
		}
		for cn, c := range validCases {
			t.Run(cn, func(t *testing.T) {
				u, err := Parse(c.urn)
				xt.OK(t, err)
				xt.Eq(t, c.exp, u.String())
				xt.Eq(t, "", u.Original) // only set by Parse
				xt.Eq(t, c.expNID, u.NID)
				xt.Eq(t, c.expNSS, u.NSS)
				xt.Eq(t, c.expFragment, u.FComponent())

				xt.Assert(t, Validates(c.urn))
			})
		}
	})

	t.Run("invalid URN", func(t *testing.T) {
		cases := map[string]string{
			"missing urn scheme":     "isbn:978-0135800911",
			"toolong NID":            "urn:abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ:too-long#NID",
			"missing NID":            "urn:978-0135800911",
			"missing NSS value":      "urn:isbn:",
			"missing NSS":            "urn:isbn",
			"bad underscore in NID":  "urn:under_scored:nid-part",
			"NID may not end with -": "urn:no-end-dash-:that-bad",
		}

		for cn, c := range cases {
			t.Run(cn, func(t *testing.T) {
				u, err := Parse(c)
				xt.KO(t, err, "was "+c)
				xt.Assert(t, u.IsZero())
				xt.Assert(t, !Validates(c))
			})
		}
	})
}

func TestURN_MarshalJSON(t *testing.T) {
	t.Run("marshal URN as JSON", func(t *testing.T) {
		urn, err := Parse("UrN:IsBn:978-0135800911#Chapter8")
		xt.OK(t, err)
		res, err := json.Marshal(urn)
		xt.OK(t, err)
		xt.Eq(t, []byte(`"urn:isbn:978-0135800911#Chapter8"`), res)
	})

	t.Run("marshal struct containing URN", func(t *testing.T) {
		var data = testDataURN{
			U: &URN{
				NID:        "example",
				NSS:        "json",
				fComponent: "struct",
			},
		}

		res, err := json.Marshal(data)
		xt.OK(t, err)
		fmt.Println(string(res))

		xt.Eq(t, []byte(`{"urn":"urn:example:json#struct"}`), res)
	})
}

type testDataURN struct {
	U *URN `json:"urn"`
}

func TestURN_UnmarshalJSON(t *testing.T) {
	t.Run("unmarshal JSON into URN", func(t *testing.T) {
		var urn URN
		xt.OK(t, urn.UnmarshalJSON([]byte(`"UrN:IsBn:978-0135800911#chapter1"`)))

		xt.Eq(t, "isbn", urn.NID)
		xt.Eq(t, "978-0135800911", urn.NSS)
		xt.Eq(t, "chapter1", urn.FComponent())
	})

	t.Run("unmarshal JSON containing invalid URN", func(t *testing.T) {
		var urn URN
		xt.KO(t, urn.UnmarshalJSON([]byte(`"UrN:spaced:[with spaces]"`)))
	})

	t.Run("unmarshal JSON object which has attribute URN", func(t *testing.T) {
		exp := URN{
			NID:        "example",
			NSS:        "json",
			fComponent: "struct",
		}
		data := testDataURN{}
		xt.OK(t, json.Unmarshal([]byte(`{"urn":"urn:example:json#struct"}`), &data))
		xt.Assert(t, exp.Equal(data.U))

		data = testDataURN{}
		xt.OK(t, json.Unmarshal([]byte(`{"urn":""}`), &data))
		xt.Eq(t, "", data.U.String())
	})
}

func TestNewURN(t *testing.T) {
	t.Run("valid parts; no optional components", func(t *testing.T) {
		u, err := New("IsBn", "978-0135800911")
		xt.OK(t, err)
		xt.Eq(t, "isbn", u.NID)
		xt.Eq(t, "978-0135800911", u.NSS)
		xt.Eq(t, "", u.RComponent())
		xt.Eq(t, "", u.QComponent())
		xt.Eq(t, "", u.FComponent())
	})

	t.Run("valid parts with r-component", func(t *testing.T) {
		u, err := New("ISBN", "978-0135800911", WithResolution("good:resolutions"))
		xt.OK(t, err)
		xt.Eq(t, "isbn", u.NID)
		xt.Eq(t, "978-0135800911", u.NSS)
		xt.Eq(t, "good:resolutions", u.RComponent())
	})

	t.Run("valid parts with q-component", func(t *testing.T) {
		u, err := New("isbn", "978-0135800911", WithQuery("callback=https://something.example.com"))
		xt.OK(t, err)
		xt.Eq(t, "isbn", u.NID)
		xt.Eq(t, "978-0135800911", u.NSS)
		xt.Eq(t, "callback=https://something.example.com", u.QComponent())
	})

	t.Run("valid parts with f-component", func(t *testing.T) {
		u, err := New("isbn", "978-0135800911", WithFragment("Chapter1"))
		xt.OK(t, err)
		xt.Eq(t, "isbn", u.NID)
		xt.Eq(t, "978-0135800911", u.NSS)
		xt.Eq(t, "Chapter1", u.FComponent())
	})

	t.Run("specify all optional components", func(t *testing.T) {
		u, err := New("isbn", "978-0135800911",
			WithResolution("resolution"), WithQuery("query=123"), WithFragment("fragment"))
		xt.OK(t, err)
		xt.Eq(t, "resolution", u.RComponent())
		xt.Eq(t, "query=123", u.QComponent())
		xt.Eq(t, "fragment", u.FComponent())

		up, err := Parse(u.String())
		xt.Eq(t, "resolution", up.RComponent())
		xt.Eq(t, "query=123", up.QComponent())
		xt.Eq(t, "fragment", up.FComponent())
	})

	t.Run("cases in NSS are retained when option provided", func(t *testing.T) {
		u, err := New("IsBn", "978-0135800911", WithNotLowerCaseNSS())
		xt.OK(t, err)
		xt.Eq(t, "IsBn", u.NID)
	})

	t.Run("invalid namespace identifier", func(t *testing.T) {
		u, err := New("no spaces allowed", "valid", WithFragment("ancher2"))
		xt.KO(t, err)
		xt.Assert(t, u == nil, "expected result of New nil")
		xt.Eq(t, "invalid namespace identifier (NID)", err.Error())
	})

	t.Run("invalid namespace specific string", func(t *testing.T) {
		u, err := New("valid-namespace", "no spaces allowed", WithFragment("ancher2"))
		xt.KO(t, err)
		xt.Assert(t, u == nil, "expected result of New nil")
		xt.Eq(t, "invalid namespace specific string (NSS)", err.Error())
	})

	t.Run("invalid component", func(t *testing.T) {
		var cases = []struct {
			opts   []Option
			expErr string
		}{
			{
				opts:   []Option{WithResolution("no spaces in component")},
				expErr: "invalid r-component",
			},
			{
				opts:   []Option{WithQuery("no spaces in component")},
				expErr: "invalid q-component",
			},
			{
				opts:   []Option{WithFragment("no spaces in component")},
				expErr: "invalid f-component",
			},
		}

		for _, c := range cases {
			t.Run(c.expErr, func(t *testing.T) {
				_, err := New("valid-namespace", "valid", c.opts...)
				xt.KO(t, err)
				xt.Eq(t, c.expErr, err.Error())
			})
		}
	})
}

func TestURN_Equal(t *testing.T) {
	t.Run("cannot compare nil", func(t *testing.T) {
		xt.Panics(t, func() {
			var u *URN
			u.Equal(nil)
		})
	})

	t.Run("equivalence according to RFC8141", func(t *testing.T) {
		var cases = map[string]struct {
			base  *URN
			eq    []string
			notEq []string
		}{
			"NID case insensitive": {
				base: mustParseURN("urn:example:a123,z456"),
				eq:   []string{"URN:example:a123,z456", "urn:EXAMPLE:a123,z456"},
			},
			"Percent-Encoded case insensitive": {
				base:  mustParseURN("urn:example:a123%2Cz456"),
				eq:    []string{"urn:example:a123%2cz456"},
				notEq: []string{"urn:example:a123,z456"},
			},
		}

		for cn, c := range cases {
			t.Run(cn, func(t *testing.T) {
				for _, s := range c.eq {
					xt.Assert(t, c.base.Equal(mustParseURN(s)), "supposed to be equivalent")
				}
				for _, s := range c.notEq {
					xt.Assert(t, !c.base.Equal(mustParseURN(s)), "supposed to be not equivalent")
				}
			})
		}
	})
}
