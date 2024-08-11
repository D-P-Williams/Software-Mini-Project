package main

import (
	"fmt"
	clihandler "work-mini-project/pkg/cliHandler"
	commandhandler "work-mini-project/pkg/commandHandler"
	"work-mini-project/pkg/configuration"
	crmhandler "work-mini-project/pkg/crmHandler"
	customerhandler "work-mini-project/pkg/customerHandler"
	transporthandler "work-mini-project/pkg/transportHandler"
)

func StartApp(commandHandler *commandhandler.CommandHandler) {
	for {
		err := commandHandler.Handle()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	config, err := configuration.LoadConfig()
	if err != nil {
		panic(err)
	}

	cliHandler := clihandler.New()

	crmHandler, err := crmhandler.New(config, cliHandler)
	if err != nil {
		panic(err)
	}

	customerHandler, err := customerhandler.New(config)
	if err != nil {
		panic(err)
	}

	transportHandler := transporthandler.New(config)

	commandHandler := commandhandler.New(config, cliHandler, crmHandler, customerHandler, transportHandler)

	cliHandler.ClearTerminal()

	StartApp(commandHandler)
}
