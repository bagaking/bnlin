package main

import (
	"context"

	"github.com/bagaking/botheater/bot"
	"github.com/bagaking/botheater/driver"
	"github.com/bagaking/botheater/driver/coze"
	"github.com/bagaking/botheater/driver/ollama"
	"github.com/khicago/got/util/typer"
	"github.com/khicago/irr"
)

type ExecutionGroup struct {
	driver string
	ak, sk string
	ep     string
	pp     string
}

const (
	DriverDoubao = "doubao"
	DriverOllama = "ollama"
)

func (eg ExecutionGroup) Use(prompt string) ExecutionGroup {
	return ExecutionGroup{
		driver: typer.Or(eg.driver, DriverDoubao),
		ak:     typer.Or(eg.ak, coze.EnvKeyVOLCAccessKey.Read()),
		sk:     typer.Or(eg.sk, coze.EnvKeyVOLCSecretKey.Read()),
		ep:     typer.Or(eg.ep, coze.EnvKeyDoubaoEndpoint.Read()),
		pp:     typer.Or(eg.pp, prompt),
	}
}

func (eg ExecutionGroup) Bot(ctx context.Context) *bot.Bot {
	var d driver.Driver
	switch eg.driver {
	case DriverOllama:
		d = ollama.New(ollama.NewClient(ctx), eg.ep)
	case DriverDoubao:
		fallthrough
	default:
		d = coze.New(coze.NewClient(ctx), eg.ep)
	}

	conf := defaultConf
	conf.Prompt.Content = eg.pp
	return bot.New(conf, d, nil)
}

func (eg ExecutionGroup) Assert() error {
	if eg.driver != DriverDoubao && eg.driver != DriverOllama {
		return irr.Error("Invalid driver")
	}

	if eg.driver == DriverOllama && eg.ep == "" {
		return irr.Error("Please provide the endpoint for the Ollama driver")
	}

	if eg.driver == DriverDoubao {
		// Check if the access key and secret key are set
		if eg.ak == "" || eg.sk == "" {
			return irr.Error("Please provide the access key and secret key using flags or environment variables")
		}
		if eg.ep == "" {
			return irr.Error("Please provide the endpoint using flags or environment variables")
		}
	}
	return nil
}
