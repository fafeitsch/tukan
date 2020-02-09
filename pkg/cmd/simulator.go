package main

import (
	"fmt"
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

	flags := []cli.Flag{
		cli.Int64Flag{Name: "port", Value: 80, Usage: "The port the simulated phone will listen to"},
		cli.StringFlag{Name: "login", Value: "Admin", Usage: "The login name for the simulator"},
		cli.StringFlag{Name: "password", Value: "admin", Usage: "The password for the simulator"},
	}

	app.HideHelp = true
	app.Flags = flags
	app.Action = func(c *cli.Context) error {
		fmt.Print(cli.AppHelpTemplate)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
