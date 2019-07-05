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
	logger := log.New(os.Stdout, "", log.LstdFlags)
	s, err := idgo.NewServer(math.MaxInt16, tcpAddr, logger)
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
