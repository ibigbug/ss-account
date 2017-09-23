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
func LoadFromFlags(portRange string) {
	c.setPortRange(portRange)
}

// LoadFromEnv ...
func LoadFromEnv() {
	if p := os.Getenv("PORT_RANGE"); p != "" {
		c.setPortRange(p)
	}
}

// GetRandomPort ...
func GetRandomPort() int {
	p := rand.Intn(c.PortEnd - c.PortStart)
	return p + c.PortStart
}

var c = &Config{}
