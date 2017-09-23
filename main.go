package main

import (
	"flag"

	raven "github.com/getsentry/raven-go"
	"github.com/ibigbug/ss-account/config"
	"github.com/ibigbug/ss-account/server"
)

func init() {
	raven.SetDSN("https://3913f327a7ef45f5af635b8bcbc2e0a4:981cefa1e6a940eeba684abb446609d1@sentry.io/163096")

}

var bind = flag.String("bind", "0.0.0.0:9000", "management server listening address")
var portRange = flag.String("port-range", "30000-40000", "accounting port range, e.g: 30000-40000")

func main() {
	flag.Parse()
	config.LoadFromFlags(*portRange)
	config.LoadFromEnv()

	if *bind == "" {
		*bind = "127.0.0.1:9000"
	}
	raven.CapturePanic(server.Start(*bind), nil)
}
