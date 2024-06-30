package cmd

import (
	"bsp/pkg/bd_img"
	"github.com/spf13/cobra"
)

// shetouCmd represents the shetou command
var shetouCmd = &cobra.Command{
	Use:   "shetou",
	Short: "baidu image shetou",
	Run: func(cmd *cobra.Command, args []string) {
		bd, err := bd_img.NewBdImg(filepath)
		if err != nil {
			panic(err)
		}
		err = bd.Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(shetouCmd)
}
