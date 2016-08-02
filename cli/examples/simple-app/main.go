package main

import (
	"fmt"
	"github.com/subchen/gstack/cli"
)

func main() {
	app := cli.NewApp("frep", "Transform template file using environment, arguments, json/yaml files")

	app.Flag("-e, --env", "set variable name=value, can be passed multiple times").Multiple()
	app.Flag("--json", "load variables from json object").Placeholder("string")
	app.Flag("--load", "load variables from json/yaml files").Multiple()
	app.Flag("--overwrite", "overwrite if destination file exists").Bool()
	app.Flag("--testing", "test mode, output transform result to console").Bool()
	app.Flag("--delims", `template tag delimiters`).Default("{{:}}")

	app.Version = "1.0.0"

	app.Usage = func() {
		fmt.Println("Usage: frep [OPTIONS] input-file:[output-file] ...")
		fmt.Println("   or: frep [ --version | --help ]")
	}

	app.MoreHelp = func() {
		fmt.Println("Examples:")
		fmt.Println("  frep nginx.conf.in -e webroot=/usr/share/nginx/html -e port=8080")
		fmt.Println("  frep nginx.conf.in:/etc/nginx.conf -e webroot=/usr/share/nginx/html -e port=8080")
		fmt.Println("  frep nginx.conf.in --json '{\"webroot\": \"/usr/share/nginx/html\", \"port\": 8080}'")
		fmt.Println("  frep nginx.conf.in --load context.json --overwrite")
	}

	app.AllowArgumentCount(1, -1)

	app.Execute = func(ctx *cli.Context) {
		fmt.Printf("-e: %v\n", ctx.StringList("-e"))
		fmt.Printf("--env: %v\n", ctx.StringList("--env"))
		fmt.Printf("--json: %v\n", ctx.String("--json"))
		fmt.Printf("--load: %v\n", ctx.StringList("--load"))
		fmt.Printf("--overwrite: %v\n", ctx.Bool("--overwrite"))
		fmt.Printf("--testing: %v\n", ctx.Bool("--testing"))
		fmt.Printf("--delims: %v\n", ctx.String("--delims"))
		fmt.Printf("Args: %v\n", ctx.Args())
		fmt.Printf("Raw Args: %v\n", ctx.RawArgs())
	}

	app.Run()
}
