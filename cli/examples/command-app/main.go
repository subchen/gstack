package main

import (
	"fmt"
	"github.com/subchen/gstack/cli"
)

func main() {
	app := cli.NewApp("docker", "A sample cli application simulate docker")

	app.Flag("--config", "Location of client config files").Default("~/.docker")
	app.Flag("-D, --debug", "Enable debug mode").Bool()
	app.Flag("-H, --host", "Daemon socket(s) to connect to").Multiple()
	app.Flag("-l, --log-level", "Set the logging level").Default("info")
	app.Flag("-tls", "Use TLS; implied by --tlsverify").Bool()

	app.Commands(
		buildCommand(),
		runCommand(),
		networkCommand(),
	)
	app.CommandRequired()

	app.Version = "1.10.1"

	app.Usage = func() {
		fmt.Println("Usage: docker [OPTIONS] COMMAND [arg...]")
		fmt.Println("       docker daemon [ --help | ... ]")
		fmt.Println("       docker [ --version | --help ]")
	}

	app.Run()
}
