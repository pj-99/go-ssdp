package main

import (
	"flag"
	"log"
	"os"

	"github.com/koron/go-ssdp"
)

func main() {
	nt := flag.String("nt", "my:device", "NT: Type")
	usn := flag.String("usn", "unique:id", "USN: ID")
	laddr := flag.String("laddr", "", "local address to listen")
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

	err := ssdp.AnnounceBye(*nt, *usn, *laddr, opts...)
	if err != nil {
		log.Fatal(err)
	}
}
