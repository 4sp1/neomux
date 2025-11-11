package main

import (
	"github.com/4sp1/neomux/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.New())
}
