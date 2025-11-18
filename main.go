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
	es := executor.NewExecutor(c)
	//es := &stub.StubDbExecutor{}
	g := gui.New(es)

	g.Start()
	//	newCli := cli.NewCLI("http://localhost:8081")
	//	newCli.Run()
}
