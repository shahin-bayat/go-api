// cmd/api-server/config/config.go
package config

import (
	"flag"
)

// Config contains configuration options parsed from command line flags.
type Config struct {
	Seed bool
	// Add more configuration options as needed
}

// ParseFlags parses command line flags and returns a Config.
func ParseFlags() Config {
	var config Config

	flag.BoolVar(&config.Seed, "seed", false, "seed the db")
	// Add more flags as needed

	flag.Parse()

	return config
}
