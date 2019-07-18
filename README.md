# idgo
idgo is fast id generator

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
	gen, err := idgo.NewIDGenerator(math.MaxInt16)
	if err != nil {
		log.Fatal(err)
	}

	// Generate a id.
	id, err := gen.Generate()
	if err != nil {
		log.Println(id)
	}

	// Generated id.
	fmt.Println(id)
	
	// Allocated id count.
	fmt.Println(gen.GetAllocatedIDCount())
	
	// id is allocated.
	fmt.Println(gen.IsAllocated(id))
	
	// Free the id.
	gen.Free(id)
	
	// Allocate the id.
	if err := gen.Allocate(id); err != nil {
		log.Println(err)
	}
	
	// Free all id.
	gen.FreeAll()
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
