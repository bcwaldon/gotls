# gotls

Because TLS is hard.

## library

Use the TLSInfo object to establish the required configuration for TLS servers and clients.

For example, use TLSInfo to configure a command-line HTTP client like this:

```Go
var info TLSInfo
flag.StringVar(&info.KeyFile, "keyfile", "", "Location of TLS key file")
flag.StringVar(&info.CertFile, "certfile", "", "Location of TLS cert file")
flag.StringVar(&info.CAFile, "cafile", "", "Location of TLS CA file")
flag.Parse()

client := &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: info.ClientConfig(),
	},
}
```

## tools

Build the term/ pkg for a TLS-terminating TCP proxy.
Use this to test your clients and servers without actually implementing TLS server-side.

```Shell
$ go build -o bin/gotls-term github.com/bcwaldon/gotls/term
$ ./bin/gotls-term --bind 127.0.0.1:34819 --proxy 127.0.0.1:43019 --key-file server.key.insecure --cert-file server.crt
2014/11/19 19:05:02 Established proxy on 127.0.0.1:34819, waiting for connections...
2014/11/19 19:05:05 127.0.0.1:58210 <-> 127.0.0.1:43019: established proxy
2014/11/19 19:05:10 127.0.0.1:58210 <-> 127.0.0.1:43019: closing proxy
2014/11/19 19:05:18 127.0.0.1:58212 <-> 127.0.0.1:43019: established proxy
2014/11/19 19:05:19 127.0.0.1:58214 <-> 127.0.0.1:43019: established proxy
2014/11/19 19:05:19 127.0.0.1:58216 <-> 127.0.0.1:43019: established proxy
2014/11/19 19:05:22 127.0.0.1:58218 <-> 127.0.0.1:43019: established proxy
2014/11/19 19:05:23 127.0.0.1:58212 <-> 127.0.0.1:43019: closing proxy
2014/11/19 19:05:24 127.0.0.1:58214 <-> 127.0.0.1:43019: closing proxy
2014/11/19 19:05:24 127.0.0.1:58216 <-> 127.0.0.1:43019: closing proxy
2014/11/19 19:05:27 127.0.0.1:58218 <-> 127.0.0.1:43019: closing proxy
2014/11/19 19:05:29 127.0.0.1:58220 <-> 127.0.0.1:43019: established proxy
2014/11/19 19:05:34 127.0.0.1:58220 <-> 127.0.0.1:43019: closing proxy
```
