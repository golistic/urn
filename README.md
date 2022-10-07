urn - Package offering URN creation and validation
==================================================

Copyright (c) 2022, Geert JM Vanderkelen

The Go urn package creates and validates Uniform Resource Names (URNs)
based on [RFC8141][1].

Overview
--------

Package urn helps to create valid URNs and offers functionality to parse or
validate strings as URN based on [RFC8141][1].

Quote from the RFC:

> A Uniform Resource Name (URN) is a Uniform Resource Identifier (URI)
> that is assigned under the "urn" URI scheme and a particular URN
> namespace, with the intent that the URN will be a persistent,
> location-independent resource identifier.

This is a simple implementation and could probably use some performance tricks.
We focus first to be compliant with the RFC.

Example URNs:

* `urn:isbn:978-0135800911`, with namespace `isbn` and specific string a
  book's ISBN13 `978-0135800911`
* `urn:ietf:rfc:8141`, the IETF (Internet Engineering Task Force) using URN 
  namespace 'ietf' to reference the RFC (Request For Comment) on which this 
  package is based on.

Quick Start
-----------

Examples can be found in the `examples_test.go` file.

### Generating a URN

One can simply generate a URN using string concatenation. However, using the
`urn.New()` function, the NID and NSS are validated.

```go
package main

import (
  "fmt"
  "log"

  "github.com/golistic/urn"
)

func main() {
  u, err := urn.New("ietf", "rfc:8141")
  if err != nil {
    log.Fatalln(err)
  }
  fmt.Printf("%s", u)
}
```

### Validating a URN

Using the `urn.Validates` you can simply check whether a URN is valid. The
same can be achieved using `urn.Parse` and checking for errors.

```go
package main

import (
  "fmt"

  "github.com/golistic/urn"
)

func main() {
  if urn.Validates("urn:ietf:rfc:8141#section-3") {
    fmt.Println("valid!")
  }

  if !urn.Validates("urn:ie+tf:rfc:8141#section-3") {
    fmt.Println("not valid!")
  }
}
```

License
-------

Distributed under the MIT license. See LICENSE.txt for more information.


[1]: https://www.rfc-editor.org/rfc/rfc8141
[2]: https://www.ietf.org