package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/koron/go-ssdp"
	"github.com/koron/go-ssdp/internal/multicast"
)

func main() {
	t := flag.String("t", ssdp.All, "search type")
	w := flag.Int("w", 10, "wait time")
	l := flag.String("l", "", "local address to listen")
	ttl := flag.Int("ttl", 0, "TTL for outgoing multicast packets")
	sysIf := flag.Bool("sysif", false, "use system assigned multicast interface")
	v := flag.Bool("v", false, "verbose mode")
	h := flag.Bool("h", false, "show help")
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	if *v {
		ssdp.Logger = log.New(os.Stderr, "[SSDP] ", log.LstdFlags)
	}

	var opts []ssdp.Option
	if *ttl > 0 {
		opts = append(opts, ssdp.TTL(*ttl))
	}
	if *sysIf {
		opts = append(opts, ssdp.OnlySystemInterface())
	}

	// Make connection for listening the search response.
	cfg, err := ssdp.Opts2config(opts)
	if err != nil {
		panic(err)
	}

	// dial multicast UDP packet.
	conn, err := multicast.Listen(&multicast.AddrResolver{Addr: *l}, cfg.MulticastConfig.Options()...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	list, err := ssdp.SearchUntil(*t, *w, conn, 1)
	if err != nil {
		log.Fatal(err)
	}
	for i, srv := range list {
		//fmt.Printf("%d: %#v\n", i, srv)
		fmt.Printf("%d: %s %s\n", i, srv.Type, srv.Location)
	}
}
