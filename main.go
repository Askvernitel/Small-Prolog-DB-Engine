package main

import (
	"weird/db/engine/client"
	"weird/db/engine/executor"
	"weird/db/engine/gui"
)

const (
	URL = "http://localhost:8080"
)

func main() {
	c := client.NewClient(URL)
	e := executor.NewExecutor(c)
	g := gui.New(e)

	g.Start()
	//	newCli := cli.NewCLI("http://localhost:8081")
	//	newCli.Run()
}
