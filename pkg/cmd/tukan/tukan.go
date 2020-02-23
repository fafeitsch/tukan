package main

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/domain"
	http2 "github.com/fafeitsch/Tukan/pkg/http"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Author = "Fabian Feitsch"
	app.Name = "Elmeg ip620/630 HTTP Configurator"
	app.Usage = "This application configures some parts of Elmeg ip620/630 telephones"

	loginFlag := cli.StringFlag{Name: "login", Value: "Admin", Usage: "The login to be used"}
	passwordFlag := cli.StringFlag{Name: "password", Value: "admin", Usage: "The password to be used"}
	portFlag := cli.IntFlag{Name: "port", Value: 80, Usage: "The port to be used to connect to the telephones"}
	ipFlag := cli.StringFlag{Name: "ip", Required: true, Usage: "The IP of the first phone to interact with"}
	numberFlag := cli.IntFlag{Name: "number", Value: 1, Usage: "The number of phones to contact, including IP"}

	scanCommand := cli.Command{
		Name:  "scan",
		Usage: "Scans an IP range for elmeg ip620/630 and tries to log into them",
		Flags: []cli.Flag{
			ipFlag,
			numberFlag,
			loginFlag,
			passwordFlag,
			portFlag,
		},
		Action: func(c *cli.Context) error {
			phoneClient := http2.BuildPhoneClient(c.Int("port"), c.String("login"), c.String("password"))
			err := phoneClient.Scan(c.String("ip"), c.Int(numberFlag.Name))
			return err
		},
	}

	phonebookUploadCommand := cli.Command{
		Name:  "pb-up",
		Usage: "Uploads a phone book to a set of elmeg ip 620/630 phones",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "ip", Required: true, Usage: "The IP of the first phone to upload to phone book to"},
			cli.StringFlag{Name: "number", Value: "1", Usage: "Number of phones to configure, starting at IP"},
			cli.StringFlag{Name: "file", Required: true, Usage: "The phone book file to upload", TakesFile: true},
			cli.StringFlag{Name: "delimiter", Usage: "A string that is not contained in the phone book. Needed for the upload. Must consist of at most 70 bytes of ASCII printable-characters", Value: "XXXXX"},
			loginFlag,
			passwordFlag,
			portFlag,
		},
		Action: func(c *cli.Context) error {
			phoneClient := http2.BuildPhoneClient(c.Int("port"), c.String("login"), c.String("password"))
			payload, err := domain.LoadAndEmbedPhonebook(c.String("file"), c.String("delimiter"))
			if err != nil {
				return fmt.Errorf("could not prepare payload for sending to phones: %v", err)
			}
			phoneClient.UploadPhoneBook(c.String("ip"), c.Int("number"), *payload, c.String("delimiter"))
			return nil
		},
	}

	phonebookDownloadCommand := cli.Command{
		Name:  "pb-down",
		Usage: "Downloads a phone book from a elmeg ip 620/630 phone",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "ip", Required: true, Usage: "The IP of the phone to download to phonebook from."},
			loginFlag,
			passwordFlag,
			portFlag,
		},
		Action: func(c *cli.Context) error {
			phoneClient := http2.BuildPhoneClient(c.Int("port"), c.String("login"), c.String("password"))
			down := phoneClient.DownloadPhoneBook(c.String("ip"))
			fmt.Printf(down)
			return nil
		},
	}

	app.Commands = []cli.Command{scanCommand, phonebookUploadCommand, phonebookDownloadCommand}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
