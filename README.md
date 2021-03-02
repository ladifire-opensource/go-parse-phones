# go-parse-phones
Simple library to detect phone numbers in a string.

## Example
```go
package main

import (
	"fmt"
	"github.com/ladifire-opensource/go-parse-phones"
)

func main() {
	message := "This is phone number: 0912525555. Other number 0943 311 366, and other number +84 968 552 221"
	phoneNumbers := goparsephone.FindInText(message, goparsephone.TypeAll)
	fmt.Println(phoneNumbers)
  // output: [{0912525555 +84912525555 Vinaphone 22 32} {0943 311 366 +84943311366 Vinaphone 47 59} {+84 968 552 221 +84968552221 Viettel 78 93}]
}
```

## Notes
This package only suport Vietnamese phone numbers yet. Other countries will be support in the next version
