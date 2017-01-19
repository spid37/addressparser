# addressparser

Parse an australian address string in go


### Example
```go
package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/spid37/addressparser"
)

func main() {
	addressString := "1234 hello st Melbourne VIC 3000"

	addressParts, err := addressparser.NewAddress(addressString)

	if err != nil {
		panic(err)
	}
	spew.Dump(addressParts)
}
```
#### Output
<pre>
 LevelType: (string) "",
 LevelNumber: (int) 0,
 FlatType: (string) "",
 FlatNumber: (int) 0,
 FlatNumberSuffix: (string) "",
 StreetNumber: (int) 1234,
 StreetNumberEnd: (int) 0,
 StreetNumberSuffix: (string) "",
 StreetName: (string) (len=5) "HELLO",
 StreetType: (string) (len=6) "STREET",
 StreetSuffix: (string) "",
 Suburb: (string) (len=9) "MELBOURNE",
 PostCode: (int) 3000,
 State: (string) (len=3) "VIC"
 </pre>
