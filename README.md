# URLVALUER

The `urlvaluer` package is a Go library that simplifies the process of parsing
struct fields into `url.Values`. This package supports various field types and
enables customization through struct tags.

## Installation

To use the `urlvaluer` package, you need to install it using the Go package
manager:

```sh
go get github.com/gaigals/urlvaluer
```

## Usage

Import the package into your code:

```go
import (
    "fmt"
    "github.com/gaigals/urlvaluer"
)
```

### Limitations

Please note that there are certain limitations to the types that are supported
for conversion to `url.Values`:

- __Maps__: Map types are not directly supported for conversion unless the
value type contains a `String() string` method. The library will skip 
unsupported map fields during marshaling.
- __Structs__: Struct types without a `String() string` method are not 
supported for conversion. Any such unsupported struct fields will be skipped
during marshaling.

### Example

Here's a comprehensive example showcasing the various field types supported 
by the package:

```go
package main

import (
    "fmt"
    "github.com/gaigals/urlvaluer"
)

type Person struct {
    Name        string      `url:"name"`
    Age         int         `url:"age"`
    Height      float32     `url:"height"`
    IsActive    bool        `url:"is_active"`
    Flags       []bool      `url:"flags"`
    Coordinates [2]int      `url:"coords"`
    Ignored     string      `url:"-"`
    NotTagged   string
    Note        *string     `url:"note"`
    Extra       any         `url:"extra"`
}

func main() {
    note := "This is a note about the person."
    extraInfo := 42 // An integer value as an example

    p := Person{
        Name:        "John Doe",
        Age:         30,
        Height:      175.5,
        IsActive:    true,
        Flags:       []bool{true, false, true},
        Coordinates: [2]int{42, 24},
        Ignored:     "This won't be included",
        NotTagged:   "This won't be included either",
        Note:        &note,
        Extra:       extraInfo,
    }

    values, err := urlvaluer.Marshal(&p)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Generated URL Values:")
    for key, val := range values {
        fmt.Printf("%q: %#v\n", key, val)
    }
}
```

In this example, the `Person` struct includes a variety of field types, 
including a pointer to a string and an `any`. The `url` tags 
are used to specify the corresponding keys in the resulting `url.Values`
map. Fields with the `-` tag are ignored. The resulting URL values are:

```
Generated URL Values:
"name":       []string{"John Doe"}
"age":        []string{"30"}
"height":     []string{"175.5"}
"is_active":  []string{"true"}
"flags":      []string{"true", "false", "true"}
"coords":     []string{"42", "24"}
"note":       []string{"This is a note about the person."}
"extra":      []string{"42"}
```

## Customization

- The `url` tag specifies the key name in the resulting URL values.
- Fields with the `-` tag are ignored and not included in the URL values.
- Pointers are included in the URL values only if they have valid values.
- Slices and arrays are supported, and their values are cast to `[]string`.

## Contributions

Contributions to this package are welcome! If you encounter any issues or have
suggestions for improvements, please create an issue on the GitHub repository.