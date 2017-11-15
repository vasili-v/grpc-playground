package main

import "flag"

var (
	address string
)

func init() {
	flag.StringVar(&address, "a", ":5555", "address:port for requests")

	flag.Parse()
}
