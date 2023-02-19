package cmd

import (
	"learn-im/cmd/business"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "start",
	Short: "this is a cmd for server start",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	Command.AddCommand(business.BusinessCmd)
}
