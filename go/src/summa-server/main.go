package main

import (
	"flag"
	"log"
	// "os"
	"summa"
)

var configFile string

func auth(username, password string) (*summa.User, error) {
	var u summa.User

	u.Username = "jbrimeyer"
	u.DisplayName = "Jeremy Brimeyer"
	u.Email = "jbrimeyer@leepfrog.com"

	// TODO: Return nil if the user is not authorized

	return &u, nil
}

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
