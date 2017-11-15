package main

import "flag"

var (
	server string
)

func init() {
	flag.StringVar(&server, "s", "localhost:5555", "address:port of server")

	flag.Parse()
}
