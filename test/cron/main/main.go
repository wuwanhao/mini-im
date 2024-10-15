package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {

}

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "hello world",
	Long:  "hello world",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello world")
	},
}

func init() {
	if err := helloCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
