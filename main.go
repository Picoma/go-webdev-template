package main

import (
	"context"
	"fmt"
	"os"

	"idp/cmd"
	"idp/internal/config"
)

//nolint:gochecknoglobals // build vars
var (
	name        = "idp"
	description = "Runs a web server serving HTMX"
	version     = "v0.0.1-a1fa420"
	commit      = "123456"
	env         = "development"
)

func main() {
	service := config.Service{
		Name:        name,
		Description: description,
		Version:     version,
		Env:         env,
		HashCommit:  commit,
	}
	cfg := config.Defaults(service)
	cmd := cmd.NewRootCmd(cfg)

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
