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
	portRange string) {

	c.setPortRange(portRange)
	c.RedisHost = redisHost
	c.RedisPort = redisPort
	c.RedisPass = redisPass
	c.RedisDB = redisDB
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
}

// GetRandomPort ...
func GetRandomPort() int {
	p := rand.Intn(c.PortEnd - c.PortStart)
	return p + c.PortStart
}

func GetRedisOptions() (host, port, pass string, db int) {
	return c.RedisHost, c.RedisPort, c.RedisPass, c.RedisDB
}

var c = &Config{}
