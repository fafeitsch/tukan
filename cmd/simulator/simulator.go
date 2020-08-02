package main

import (
	"fmt"
	"github.com/fafeitsch/Tukan/tukan/mock"
	"github.com/goccy/go-yaml"
	"github.com/urfave/cli"
	"io/ioutil"
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
	var parametersFile string
	flags := []cli.Flag{
		cli.IntFlag{Name: "port", Value: 80, Usage: "The port the simulated phone will listen to", Destination: &port},
		cli.StringFlag{Name: "login", Value: "Admin", Usage: "The login name for the simulator", Destination: &login},
		cli.StringFlag{Name: "password", Value: "admin", Usage: "The password for the simulator", Destination: &password},
		cli.StringFlag{Name: "parameters", Value: "", Usage: "yaml or json file containing the parameters", Destination: &parametersFile},
	}

	app.HideHelp = true
	app.Flags = flags
	app.Action = func(c *cli.Context) error {
		handler, phone := mock.CreatePhone(login, password)
		if len(parametersFile) != 0 {
			data, err := ioutil.ReadFile(parametersFile)
			if err != nil {
				log.Fatal(err)
			}
			err = yaml.Unmarshal(data, &phone.Parameters)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
