package cmd

import (
	"bsp/pkg/techan"
	"github.com/spf13/cobra"
)

var teChanCmd = &cobra.Command{
	Use:   "techan",
	Short: "https://www.zhtechan.cn",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var regionCmd = &cobra.Command{
	Use:   "region",
	Short: "all region",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := techan.NewFullRegion(filepath).Run()
		if err != nil {
			panic(err)
		}
	},
}

var singleCmd = &cobra.Command{
	Use:   "single",
	Short: "process all single items",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := techan.NewSingle(filepath).Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(teChanCmd)

	teChanCmd.AddCommand(regionCmd)
	teChanCmd.AddCommand(singleCmd)
}
