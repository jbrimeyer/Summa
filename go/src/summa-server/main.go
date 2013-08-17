package main

import (
	"flag"
	"log"
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

	summa.SetAuthProvider(auth)
	summa.StartHttp()
}

func auth(username, password string) (*summa.User, error) {
	var u summa.User

	u.Username = "anonymous"
	u.DisplayName = "Anonymous"
	u.Email = "anon@anonymous.com"

	return &u, nil
}
