package main

import (
	"log"

	"github.com/BIwashi/xpipecd-xbar/pkg/cli"
	"github.com/BIwashi/xpipecd-xbar/pkg/pipectl"
)

func main() {
	c := cli.NewCLI(
		"pipecd",
		"Pipectl command line tool for xbar",
	)
	c.AddCommands(
		pipectl.NewCommand(),
	)
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}
