package main

import "flag"

var (
	address        string
	perStreamLimit int
)

func init() {
	flag.StringVar(&address, "a", ":5555", "address:port for requests")
	flag.IntVar(&perStreamLimit, "limit", 0, "max number of messages stream can handle")

	flag.Parse()
}
