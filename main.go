package main

import (
	"flag"

	"github.com/trilopin/godinary/server"
)

func main() {
	port := flag.Int("port", 3001, "Port to listen to")
	flag.Parse()
	server.StartServer(*port)
}
