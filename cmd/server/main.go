// Command server is the entrypoint for the Static Asset Management Service.
package main

import (
	"flag"
	"fmt"
	"os"

	"m-sfs/internal/config"
	"m-sfs/internal/logging"
)

func main() {
	configPath := flag.String("config", "config.json", "path to config.json")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := logging.New(cfg.Logging.Level)
	logger.Info("starting server",
		"config_path", *configPath,
		"log_level", cfg.Logging.Level,
		"database_driver", cfg.Database.Driver,
	)

	// Later phases wire up the database, auth middleware, and HTTP router
	// here before blocking on ListenAndServe.

	logger.Info("shutting down")
}
