package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/bcwaldon/gotls"
)

var (
	info gotls.TLSInfo
)

func main() {
	fs := flag.NewFlagSet("gotls", flag.ExitOnError)
	bind := fs.String("bind", "0.0.0.0:39281", "Listen on this TCP address")
	remote := fs.String("proxy", "", "Proxy connections to this TCP address")
	fs.StringVar(&info.CAFile, "ca-file", "", "Location of TLS CA file")
	fs.StringVar(&info.CertFile, "cert-file", "", "Location of TLS cert file")
	fs.StringVar(&info.KeyFile, "key-file", "", "Location of TLS key file")

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatalf("Failed parsing flags: %v", err)
	}

	l, err := newListener(*bind, info)
	if err != nil {
		log.Fatalf("Failed creating listener: %v", err)
	}

	raddr, err := net.ResolveTCPAddr("tcp", *remote)
	if err != nil {
		log.Fatalf("Failed resolving remote address: %v", err)
	}

	log.Printf("Established proxy on %v, waiting for connections...", *bind)
	proxy(l, raddr)
}

func newListener(addr string, info gotls.TLSInfo) (net.Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	cfg, err := info.ServerConfig()
	if err != nil {
		return nil, err
	}

	return tls.NewListener(l, cfg), nil
}

func proxy(l net.Listener, raddr *net.TCPAddr) {
	for {
		ic, err := l.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection: %v", err)
			continue
		}
		ic.SetDeadline(time.Now().Add(5 * time.Second))
		go func() {
			oc, err := net.DialTCP("tcp", nil, raddr)
			if err != nil {
				log.Printf("Failed dialing remote address: %v", err)
				return
			}
			defer oc.Close()

			oc.SetDeadline(time.Now().Add(5 * time.Second))

			log.Printf("%v <-> %v: established proxy", oc.LocalAddr(), oc.RemoteAddr())

			go io.Copy(oc, ic)
			io.Copy(ic, oc)

			log.Printf("%v <-> %v: closing proxy", oc.LocalAddr(), oc.RemoteAddr())
		}()
	}
}
