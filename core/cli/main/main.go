package main

import (
	"fmt"
	"github.com/LeeroyLin/goengine/core/cli"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cli.NewRootCmd("welcome")

	helloCmd := rootCmd.NewSubCmd("hello", "xxx", "xxx", func(cmd *cobra.Command) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			name = "Guest"
		}
		fmt.Printf("Hello, %s! 👋\n", name)
	})

	helloCmd.AddStringFlag("name", "n", "", "xxx")

	rootCmd.Run()
}
