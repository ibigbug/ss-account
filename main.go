package main

import (
	"flag"
	"log"

	"github.com/ibigbug/ss-account/config"
	"github.com/ibigbug/ss-account/server"
)

var bind = flag.String("bind", "0.0.0.0:9000", "management server listening address")
var portRange = flag.String("port-range", "30000-40000", "accounting port range, e.g: 30000-40000")

func main() {
	flag.Parse()
	config.LoadFromFlags(*portRange)
	config.LoadFromEnv()

	if *bind == "" {
		*bind = "127.0.0.1:9000"
	}
	log.Fatalln(server.Start(*bind))
}
