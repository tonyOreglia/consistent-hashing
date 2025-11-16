package main

import (
	"context"
	"log"
	"os"

	"consistent_hash"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "server",
		Usage: "Start the web server",
		Action: func(context.Context, *cli.Command) error {
			s := consistent_hash.NewServer()
			return s.Run(":8080")
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
