package main

import (
	"fmt"
	"github.com/subchen/gstack/cli"
)

func runCommand() *cli.Command {
	var cmd = cli.NewCommand("run", "Run a command in a new container")

	cmd.Flag("--link", "Add link to another container").Multiple()
	cmd.Flag("-P, --publish-all", "Publish all exposed ports to random ports").Bool()
	cmd.Flag("--rm", "Automatically remove the container when it exits").Bool()
	cmd.Flag("-v, --volume", "Bind mount a volume").Multiple()

	cmd.AllowArgumentCount(1, 1)

	cmd.Usage = "Usage: docker run [OPTIONS] IMAGE"

	cmd.Execute = func(ctx *cli.Context) {
		fmt.Printf("link: %v\n", ctx.StringList("--link"))
		fmt.Printf("volume: %v\n", ctx.StringList("-v"))
		fmt.Printf("docker running %s\n", ctx.Arg(0))
	}

	return cmd
}
