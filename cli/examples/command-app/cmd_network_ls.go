package main

import (
	"fmt"
	"github.com/subchen/gstack/cli"
)

func networkLsCommand() *cli.Command {
	var cmd = cli.NewCommand("ls", "Lists networks")

	cmd.Flag("-q, --quiet", "Only display numeric IDs")
	cmd.Flag("--no-trunc", "Do not truncate the output").Bool()

	cmd.Execute = func(ctx *cli.Context) {
		fmt.Printf("docker network ls\n")
	}

	return cmd
}
