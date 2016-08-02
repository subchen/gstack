package main

import (
	"fmt"
	"github.com/subchen/gstack/cli"
)

func buildCommand() *cli.Command {
	cmd := cli.NewCommand("build", "Build an image from a Dockerfile")

	cmd.Flag("-f, --file", "Name of the Dockerfile (Default is 'PATH/Dockerfile')")
	cmd.Flag("--no-cache", "Do not use cache when building the image").Bool()
	cmd.Flag("--label", "Set metadata for an image").Multiple()

	cmd.AllowArgumentCount(1, 1)

	cmd.Usage = "Usage: docker build [OPTIONS] PATH | URL | -"

	cmd.Execute = func(ctx *cli.Context) {
		fmt.Println("Global Options: ")
		fmt.Printf("  --config = %s\n", ctx.Global().String("--config"))
		fmt.Printf("  --log-level = %s\n", ctx.Global().String("--log-level"))
		fmt.Println()
		fmt.Println("Options: ")
		fmt.Printf("  -f = %s\n", ctx.String("-f"))
		fmt.Printf("  --label = %s\n", ctx.StringList("--label"))
		fmt.Println()
		fmt.Printf("docker building %s\n", ctx.Arg(0))
	}

	return cmd
}
