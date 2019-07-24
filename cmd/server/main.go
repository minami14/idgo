package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/minami14/idgo/idgo"
)

type Options struct {
	MinSize int    `short:"min" long:"minimum" description:"Minimum value of ID to be generated" default:"1"`
	MaxSize int    `short:"max" long:"maximum" description:"Maximum value of ID to be generated" default:"2147483647"`
	Port    uint16 `short:"p" long:"port" description:"Port number" default:"49152"`
}

func main() {
	var options Options
	if _, err := flags.Parse(&options); err != nil {
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%d", options.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	s, err := idgo.NewServer(options.MinSize, options.MaxSize, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	s.SetLogger(log.New(os.Stdout, "idgo ", log.LstdFlags))

	fmt.Println("idgo server started")
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
