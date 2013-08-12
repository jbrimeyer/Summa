package main

import (
	"flag"
	"log"
	// "os"
	"summa"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "f", "server.conf", "The server configuration file")
}

func main() {
	flag.Parse()

	err := summa.Init(configFile)
	if err != nil {
		log.Fatalf("Could not initialize Summa: %s", err)
	}

	summa.StartHttp()
}
