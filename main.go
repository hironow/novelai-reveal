package main

import (
	"log"
	"novelai-reveal/novelai"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "novelai-reveal",
		Version: "v0.0.1",
		Usage:   "reveal NovelAI images!",
		Action:  novelai.CheckDirectory,
		Suggest: true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
