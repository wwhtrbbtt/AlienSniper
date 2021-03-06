package main

import (
	host "Alien/host"
	node "Alien/node"
	types "Alien/types"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/urfave/cli/v2"
)

const version = "0.0.1"

var clear map[string]func() //create a map for storing clear funcs
var prefix = map[string]string{
	"windows": "",
	"darwin":  "./",
	"linux":   "./",
}

// https://stackoverflow.com/a/22896706
func init() {
	clear = make(map[string]func()) //Initialize it
	unix := func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["linux"] = unix
	clear["darwin"] = unix

	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func ClearScreen() {
	// Clear screen on supported platforms

	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! Supported: Darwin (MacOS), Linux, Windows")
	}
}

func main() {
	ClearScreen()
	PrintLogo()
	Handle()
	fmt.Print("\n\n")
}

func PrintLogo() {
	text := "               _,--=--._\n"
	text += "             ,'    _    `.\n"
	text += "            -    _(_)_o   - \n"
	text += "       ____'    /_  _/]    `____\n"
	text += "-=====::(+):::::::::::::::::(+)::=====-\n"
	text += `         (+).""""""""""""",(+)` + "\n"
	text += "             .           ,\n"
	text += "               `  -=-  '\n"
	fmt.Println(text)
}

func Handle() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "help",
				Action: func(c *cli.Context) error {
					usage := fmt.Sprintf("Usage: %salien [options]\n\n", prefix[runtime.GOOS])
					usage += "Options:\n"
					usage += "    help: Print this help message\n"
					usage += "    version: Print the version number\n"
					usage += "    configure: Configure the application\n"
					usage += "    start: Start the CLI\n"
					usage += "    host: Start as the host in the background\n"
					usage += "    node: Start as a node in the background\n"

					fmt.Println(usage)
					return nil
				},
			},
			{
				Name: "version",
				Action: func(c *cli.Context) error {
					fmt.Println("Version: " + version)
					return nil
				},
			},
			{
				Name: "configure",
				Action: func(ca *cli.Context) error {
					c := types.Configure()
					c.SaveToFile()
					os.Exit(0)
					return nil
				},
			},
			{
				Name: "start",
				Action: func(c *cli.Context) error {
					fmt.Println("Starting...")
					os.Exit(0)
					return nil
				},
			},
			{
				Name: "node",
				Action: func(c *cli.Context) error {
					fmt.Println("Starting as node...")
					node.Start(c.String("ip"), c.String("key"))
					os.Exit(0)
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "ip",
						Usage: "ip of the host",
					},
					&cli.StringFlag{
						Name:  "key",
						Usage: "key of the host",
					},
				},
			},
			{
				Name: "host",
				Action: func(c *cli.Context) error {
					fmt.Println("Starting as host...")
					host.Start()
					os.Exit(0)
					return nil
				},
			},
		},

		Name:    "Alien",
		Usage:   "Faster than every cowboy",
		Version: version,
	}

	app.Run(os.Args)
}
