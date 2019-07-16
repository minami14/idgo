package main

import (
	"log"
	"math"
	"net"
	"os"

	"github.com/minami14/idgo/idgo"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":4000")
	if err != nil {
		log.Fatal(err)
	}

	s, err := idgo.NewServer(math.MaxInt16, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	s.SetLogger(log.New(os.Stdout, "", log.LstdFlags))

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
