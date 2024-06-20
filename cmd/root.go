package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "bsp",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

var (
	filepath string
)

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&filepath, "file", "f", "", "Input file path")
}
