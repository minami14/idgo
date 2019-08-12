# idgo
[![CircleCI](https://circleci.com/gh/minami14/idgo.svg?style=shield)](https://circleci.com/gh/minami14/idgo)

idgo is a very fast id generator that generates an int id that can specify the maximum value

# Install
```bash
go get github.com/minami14/idgo
```

# Usage
```go
package main

import (
	"fmt"
	"log"
	"math"

	"github.com/minami14/idgo/idgo"
)

func main() {
	store := idgo.NewLocalStore(math.MaxInt16)
	gen, err := idgo.NewIDGenerator(store)
	if err != nil {
		log.Fatal(err)
	}

	// Generate a id.
	id, err := gen.Generate()
	if err != nil {
		log.Println(err)
	}

	// Generated id.
	fmt.Println(id)

	// Allocated id count.
	fmt.Println(gen.GetAllocatedIDCount())

	// id is allocated.
	fmt.Println(gen.IsAllocated(id))
	isAllocated, err := gen.IsAllocated(id)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(isAllocated)

	// Free the id.
	if err := gen.Free(id); err != nil {
		log.Println(err)
	}

	// Allocate the id.
	if err := gen.Allocate(id); err != nil {
		log.Println(err)
	}

	// Free all id.
	if err := gen.FreeAll(); err != nil {
		log.Println(err)
	}
}

```

## Server

### Build
```bash
go build -o idgo-server ./cmd/server/main.go
chmod +x idgo-server
```

### Show usage
```bash
./idgo-server -h
```

### Run server
```bash
./idgo-server -m [maximum value of id] -p [port number]
```

## Client
```go
package main

import (
	"fmt"
	"log"
	"net"

	"github.com/minami14/idgo/idgo"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":49152")
	if err != nil {
		log.Fatal(err)
	}
	
	client := idgo.NewClient()
	if err := client.Connect(tcpAddr); err != nil {
		log.Fatal(err)
	}
	
	id, err := client.Generate()
	if err != nil {
		log.Println(id)
	}
	
	if err := client.Free(id); err != nil {
		log.Println(err)
	}
}
```

# License
MIT License

idgo uses the following libraries

* [go-flags](https://github.com/jessevdk/go-flags/blob/master/LICENSE) Copyright (c) 2012 Jesse van den Kieboom
* [redigo](https://github.com/gomodule/redigo/blob/master/LICENSE)
