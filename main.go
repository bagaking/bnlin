// Package main provides a command-line tool for generating commit comments based on diff information.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bagaking/botheater/driver/coze"
	"github.com/urfave/cli/v2"

	"github.com/bagaking/botheater/bot"
	"github.com/bagaking/easycmd"
)

const (
	CMDNameRun = "run"
)

// defaultConf is the default configuration for the bot.
var defaultConf = bot.Config{
	Endpoint:   "",
	PrefabName: "bnlin",
	Prompt: &bot.Prompt{
		Content: "",
	},
}

func runApp() error {
	app := easycmd.New("bnlin").Set.Custom(func(command *cli.Command) {
		command.Usage = `bnlin is an AI-powered command-line tool that automatically generates and executes
bash scripts based on natural language input. It interprets user requests and 
translates them into precise bash commands to achieve the desired results.

Main Usage: 
   bnlin run <the proposed commands in natural language>
`
	}).End

	app.Child(CMDNameRun).Flags(
		&cli.StringFlag{Name: "access_key", Usage: fmt.Sprintf("Access key for the API (alternative to %s)", coze.EnvKeyVOLCAccessKey), Aliases: []string{"ak"}, Required: false},
		&cli.StringFlag{Name: "secret_key", Usage: fmt.Sprintf("Secret key for the API (alternative to %s)", coze.EnvKeyVOLCAccessKey), Aliases: []string{"sk"}, Required: false},
		&cli.StringFlag{Name: "endpoint", Usage: fmt.Sprintf("Endpoint for generating the comment (alternative to  %s)", coze.EnvKeyDoubaoEndpoint), Aliases: []string{"e"}, Required: false},
		&cli.StringFlag{Name: "prompt", Usage: "Custom prompt for generating the comment", Aliases: []string{"p"}, Required: false},
	).Set.Alias("r").Usage(`run <commands in natural language>

Example:
   bnlin run find all uncommitted files and list their line counts
   bnlin run -ak your_access_key -sk your_secret_key -e your_endpoint find all uncommitted files and list their line counts

Environment Variables:
   VOLC_ACCESS_KEY      Access key for the API (alternative to -ak)
   VOLC_SECRET_KEY      Secret key for the API (alternative to -sk)
   DOU_BAO_ENDPOINT     Endpoint for the API (alternative to -e)`).End.Action(func(c *cli.Context) error {
		eg := ExecutionGroup{
			ak: c.String("access_key"),
			sk: c.String("secret_key"),
			ep: c.String("endpoint"),
			pp: c.String("prompt"),
		}
		if c.Args().Len() == 0 {
			return cli.ShowCommandHelp(c, CMDNameRun)
		}

		cmdInNL := strings.Join(c.Args().Slice(), " ")
		fmt.Println("the proposed commands in natural language is", cmdInNL)
		return autoComment(c.Context, cmdInNL, eg)
	})

	return app.RunBaseAsApp()
}

func main() {
	if err := runApp(); err != nil {
		fmt.Printf("=== EXECUTION FAILED===\n\n%s\n===\n", err)
		os.Exit(1)
	}
}
