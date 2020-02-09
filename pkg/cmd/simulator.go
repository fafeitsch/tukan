package main

import (
	"github.com/fafeitsch/Tukan/pkg/mock"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Author = "Fabian Feitsch"
	app.Name = "Elmeg ip620/630 HTTP Simulator"
	app.Usage = "This application simulates a limited set of HTTP endpoints of the Elmeg IP 620/630 phones."
	app.UsageText = "Call simulator with appropriate global options (see below)"

	var port int
	var login string
	var password string
	flags := []cli.Flag{
		cli.IntFlag{Name: "port", Value: 80, Usage: "The port the simulated phone will listen to", Destination: &port},
		cli.StringFlag{Name: "login", Value: "Admin", Usage: "The login name for the simulator", Destination: &login},
		cli.StringFlag{Name: "password", Value: "admin", Usage: "The password for the simulator", Destination: &password},
	}

	app.HideHelp = true
	app.Flags = flags
	app.Action = func(c *cli.Context) error {
		mock.StartHandler(port, login, password)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
