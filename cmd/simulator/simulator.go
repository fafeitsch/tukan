package main

import (
	"fmt"
	"github.com/fafeitsch/Tukan/tukan/mock"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Author = "Fabian Feitsch"
	app.Name = "Tukan HTTP/REST Phone Simulator"
	app.Usage = "This application simulates a REST server having endpoints which may be found on some IP/SIP telephones. This " +
		"application should only be used to test to Tukan Phone Configurator."
	app.UsageText = "Call simulator with appropriate global options (see below)."

	var port int
	var login string
	var password string
	var fnKeysFile string
	flags := []cli.Flag{
		cli.IntFlag{Name: "port", Value: 80, Usage: "The port the simulated phone will listen to", Destination: &port},
		cli.StringFlag{Name: "login", Value: "Admin", Usage: "The login name for the simulator", Destination: &login},
		cli.StringFlag{Name: "password", Value: "admin", Usage: "The password for the simulator", Destination: &password},
		cli.StringFlag{Name: "functionKeys", Value: "", Usage: "CSV file of function key definitions", Destination: &fnKeysFile},
	}

	app.HideHelp = true
	app.Flags = flags
	app.Action = func(c *cli.Context) error {
		handler, phone := mock.CreatePhone(login, password)
		if len(fnKeysFile) != 0 {
			csv, err := mock.ParseFunctionKeysCsv(fnKeysFile)
			if err != nil {
				log.Fatal(err)
			}
			phone.Parameters.FunctionKeys = csv
		}
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
