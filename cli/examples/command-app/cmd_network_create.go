package main

import (
	"fmt"
	"github.com/subchen/gstack/cli"
)

func networkCreateCommand() *cli.Command {
	var cmd = cli.NewCommand("create", "Creates a new network with a name specified by the user")

	cmd.Flag("--subnet", "Only display numeric IDs").Multiple()
	cmd.Flag("--gateway", "ipv4 or ipv6 Gateway for the master subnet").Multiple()
	cmd.Flag("-d, --driver", "Driver to manage the Network").Default("bridge")
	cmd.Flag("---ipv6", "enable IPv6 networking").Bool()

	cmd.AllowArgumentCount(1, 1)

	cmd.Usage = "Usage: docker network create [OPTIONS] NETWORK-NAME"

	cmd.Execute = func(ctx *cli.Context) {
		fmt.Printf("docker network create %s\n", ctx.Arg(0))
	}

	return cmd
}
