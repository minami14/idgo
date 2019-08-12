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
	MaxSize int    `short:"m" long:"max" description:"Maximum value of ID to be generated" default:"2147483647"`
	Port    uint16 `short:"p" long:"port" description:"Port number" default:"49152"`
	Redis   string `short:"r" long:"redis" description:"Redis server hostname" default:""`
	Key     string `short:"k" long:"key" description:"Redis key" default:"idgo"`
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

	var store idgo.AllocatedIDStore
	if options.Redis == "" {
		store, err = idgo.NewLocalStore(options.MaxSize)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		store, err = idgo.NewRedisStore(options.Redis, options.Key, options.MaxSize)
		if err != nil {
			log.Fatal(err)
		}
	}

	s, err := idgo.NewServer(store, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	s.SetLogger(log.New(os.Stdout, "idgo ", log.LstdFlags))

	fmt.Println("idgo server started")
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
