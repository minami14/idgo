package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/minami14/idgo/idgo"
)

type options struct {
	MaxSize int    `short:"m" long:"max" description:"Maximum value of ID to be generated" default:"2147483647"`
	Port    uint16 `short:"p" long:"port" description:"Port number" default:"49152"`
	Redis   string `short:"r" long:"redis" description:"Redis server hostname" default:""`
	Key     string `short:"k" long:"key" description:"Redis key" default:"idgo"`
}

func main() {
	var opt options
	if _, err := flags.Parse(&opt); err != nil {
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%d", opt.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	var store idgo.AllocatedIDStore
	if opt.Redis == "" {
		store, err = idgo.NewLocalStore(opt.MaxSize)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		store, err = idgo.NewRedisStore(opt.Redis, opt.Key, opt.MaxSize)
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
