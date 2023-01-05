package main

import (
	"log"
	"os"

	server "github.com/ssbostan/prometheus-multi-tenant-proxy-server/internal/prometheus-multi-tenant-proxy-server"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "prometheus-multi-tenant-proxy-server",
		Usage: "Multi-tenant reverse proxy for Prometheus server",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Saeid Bostandoust",
				Email: "ssbostan@yahoo.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "Configuration file",
				Value:   "config.yaml",
				Aliases: []string{"c"},
			},
		},
		Action: server.Run,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
