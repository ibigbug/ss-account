package config

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// Config ...
type Config struct {
	PortStart int
	PortEnd   int

	RedisHost string
	RedisPort string
	RedisPass string
	RedisDB   int

	StripeKey    string
	StripeSecret string
}

func (c *Config) setPortRange(portRange string) {
	r := strings.Split(portRange, "-")
	if len(r) != 2 {
		log.Panicln("Wrong port range")
	}
	s, err := strconv.Atoi(r[0])
	if err != nil {
		log.Panicln("Wrong port start")
	}
	e, err := strconv.Atoi(r[1])
	if err != nil {
		log.Panicln("Wrong port end")
	}
	if s > e || s < 1 || e > 65535 {
		log.Panicln("Port range must between 1-65535")
	}
	c.PortStart = s
	c.PortEnd = e
}

// LoadFromFlags ...
func LoadFromFlags(
	redisHost,
	redisPort,
	redisPass string,
	redisDB int,
	portRange string,
	stripeKey,
	stripeSecret string) {

	c.setPortRange(portRange)
	c.RedisHost = redisHost
	c.RedisPort = redisPort
	c.RedisPass = redisPass
	c.RedisDB = redisDB
	c.StripeKey = stripeKey
	c.StripeSecret = stripeSecret
}

// LoadFromEnv ...
func LoadFromEnv() {
	if p := os.Getenv("PORT_RANGE"); p != "" {
		c.setPortRange(p)
	}
	if p := os.Getenv("REDIS_HOST"); p != "" {
		c.RedisHost = p
	}
	if p := os.Getenv("REDIS_PORT"); p != "" {
		c.RedisPort = p
	}
	if p := os.Getenv("REDIS_PASS"); p != "" {
		c.RedisPass = p
	}
	if p := os.Getenv("REDISS_DB"); p != "" {
		c.RedisPass = p
	}

	if p := os.Getenv("STRIPE_KEY"); p != "" {
		c.StripeKey = p
	}
	if p := os.Getenv("STRIPE_SECRET"); p != "" {
		c.StripeSecret = p
	}
}

// GetRandomPort ...
func GetRandomPort() string {
	p := rand.Intn(c.PortEnd - c.PortStart)
	return strconv.Itoa(p + c.PortStart)
}

// GetRedisOptions return redis options
func GetRedisOptions() (host, port, pass string, db int) {
	return c.RedisHost, c.RedisPort, c.RedisPass, c.RedisDB
}

func GetStripeKey() string {
	return c.StripeKey
}

var c = &Config{}
