// Copyright (c) 2022, Geert JM Vanderkelen

package urn_test

import (
	"fmt"
	"log"

	"github.com/golistic/urn"
)

func ExampleNew() {
	u, err := urn.New("ietf", "rfc:8141")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s", u)

	// Output:
	// urn:ietf:rfc:8141
}

func ExampleURN_Equal() {
	// err handling left out for sake of brevity
	u, _ := urn.New("ietf", "rfc:8141")
	equalNID, _ := urn.New("IETF", "rfc:8141", urn.WithNotLowerCaseNSS())
	fmt.Printf("%s == %s : %v\n", u, equalNID, u.Equal(equalNID))

	notEqualNSS, _ := urn.New("ietf", "RFC:8141")
	fmt.Printf("%s == %s : %v\n", u, notEqualNSS, u.Equal(notEqualNSS))

	// components are ignored when determining equivalence
	withComponent, _ := urn.New("ietf", "rfc:8141", urn.WithFragment("section-3"))
	fmt.Printf("%s == %s : %v\n", u, withComponent, u.Equal(withComponent))

	// Output:
	// urn:ietf:rfc:8141 == urn:IETF:rfc:8141 : true
	// urn:ietf:rfc:8141 == urn:ietf:RFC:8141 : false
	// urn:ietf:rfc:8141 == urn:ietf:rfc:8141#section-3 : true
}

func ExampleParse() {
	u, err := urn.Parse("urn:ietf:rfc:8141#section-3")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s : %v\n", u, u.Equal(&urn.URN{NID: "ietf", NSS: "rfc:8141"}))

	// Output:
	// urn:ietf:rfc:8141#section-3 : true
}

func ExampleValidates() {
	if urn.Validates("urn:ietf:rfc:8141#section-3") {
		fmt.Println("valid!")
	}

	if !urn.Validates("urn:ie+tf:rfc:8141#section-3") {
		fmt.Println("not valid!")
	}

	// Output:
	// valid!
	// not valid!
}
