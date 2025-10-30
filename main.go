package main

import "weird/db/engine/cli"

func main() {
	cli := cli.NewCLI("http://localhost:8081")
	cli.Run()

}
