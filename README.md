# A Golang and Command-Line Interface to Ghostarchive

This package is a command-line tool named `ghostarchive` saving webpage to [Ghostarchive](https://ghostarchive.org), it also supports imports as a Golang package for a programmatic. Please report all bugs and issues on [Github](https://github.com/wabarc/ghostarchive/issues).

## Installation

From source:

```sh
go get github.com/wabarc/ghostarchive
```

From [gobinaries.com](https://gobinaries.com):

```sh
curl -sf https://gobinaries.com/wabarc/ghostarchive/cmd/ghostarchive | sh
```

From [releases](https://github.com/wabarc/ghostarchive/releases)

## Usage

### Command-line

```sh
$ ghostarchive https://example.com

Output:
version: 0.0.1
date: unknown

https://example.com => https://ghostarchive.org/archive/eb4i3
```

### Go package interfaces

```go
package main

package ga

import (
        "fmt"

        "github.com/wabarc/ghostarchive"
)

func main() {
        wbrc := &ga.Archiver{}
        saved, _ := wbrc.Wayback(args)
        for orig, dest := range saved {
                fmt.Println(orig, "=>", dest)
        }
}

// Output:
// https://example.com => https://ghostarchive.org/archive/eb4i3
```

## License

This software is released under the terms of the MIT. See the [LICENSE](https://github.com/wabarc/ghostarchive/blob/main/LICENSE) file for details.
