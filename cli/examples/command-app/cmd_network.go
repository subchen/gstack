package main

import (
	"github.com/subchen/gstack/cli"
)

func networkCommand() *cli.Command {
	var cmd = cli.NewCommand("network", "Manage Docker networks")

	//cmd.Usage = "Usage: docker network [OPTIONS] COMMAND [OPTIONS]c"

	cmd.SubCommands(
		networkLsCommand(),
		networkCreateCommand(),
	)
	cmd.SubCommandRequired()

	return cmd
}
