package main

import (
	"github.com/Kifen/go-blockchain/cli"
	"os"
)

func main(){
	defer os.Exit(0)
	cli := cli.CLI{}
	cli.Run()
}

