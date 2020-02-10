package main

import (
	http2 "github.com/fafeitsch/Tukan/pkg/http"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Author = "Fabian Feitsch"
	app.Name = "Elmeg ip620/630 HTTP Configurator"
	app.Usage = "This application configures some parts of Elmeg ip620/630 telephones"

	scanCommand := cli.Command{
		Name:  "scan",
		Usage: "Scans an IP range for elmeg ip620/630 and tries to log into them",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "cidr", Value: "192.168.2.0/24", Usage: "The IP range to scan"},
			cli.IntFlag{Name: "port", Value: 80, Usage: "The port to be used to connect to the telephones"},
			cli.StringFlag{Name: "login", Value: "Admin", Usage: "The login to be used"},
			cli.StringFlag{Name: "password", Value: "admin", Usage: "The password to be used"},
		},
		Action: func(c *cli.Context) error {
			client := &http.Client{
				Timeout: 20 * time.Second,
			}
			phoneClient := http2.PhoneClient{Client: client}
			err := phoneClient.Scan(c.String("cidr"), c.Int("port"), c.String("login"), c.String("password"))
			return err
		},
	}

	app.Commands = []cli.Command{scanCommand}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
