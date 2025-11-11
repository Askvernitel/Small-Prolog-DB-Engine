package main

import (
	"weird/db/engine/gui"
	"weird/db/engine/stub"
)

const (
	URL = "http://localhost:8080"
)

func main() {
	//c := client.NewClient(URL)
	//e := executor.NewExecutor(c)
	es := &stub.StubDbExecutor{}
	g := gui.New(es)

	g.Start()
	//	newCli := cli.NewCLI("http://localhost:8081")
	//	newCli.Run()
}
