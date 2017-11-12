package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	raven "github.com/getsentry/raven-go"
	"github.com/ibigbug/ss-account/config"
	"github.com/ibigbug/ss-account/server"
	stripe "github.com/stripe/stripe-go"
)

var (
	bind      = flag.String("bind", "0.0.0.0:9000", "management server listening address")
	portRange = flag.String("port-range", "30000-40000", "accounting port range, e.g: 30000-40000")

	redisHost = flag.String("redis-host", "localhost", "redis host")
	redisPort = flag.String("redis-port", "6379", "redis port")
	redisPass = flag.String("redis-pass", "", "redis password")
	redisDB   = flag.Int("redis-db", 0, "redis database")

	raveDSN = flag.String("sentry-dsn", "", "sentry DSN")

	stripeKey    = flag.String("stripe-key", "", "Stripe API key")
	stripeSecret = flag.String("stripe-secret", "", "Stripe secret")
)

func waitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Printf("Got sig: %v, exiting\n", sig)
}

func main() {
	flag.Parse()
	config.LoadFromFlags(
		*redisHost,
		*redisPort,
		*redisPass,
		*redisDB,
		*portRange,
		*stripeKey,
		*stripeSecret,
	)
	config.LoadFromEnv()

	if *bind == "" {
		*bind = "127.0.0.1:9000"
	}

	stripe.Key = *stripeSecret

	raven.SetDSN(*raveDSN)
	go func() {
		raven.CaptureError(server.Start(*bind), nil)
	}()

	waitForSignal()
}
