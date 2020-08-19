package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// goCkup download -> get to make choose
// goCkup download --decrypt -> get to make choose
// set decryption key
// set download destination
var CmdDownload = &cobra.Command{
	Use:   "download _todo_",
	Short: "_todo_",
	Long:  `_todo_`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}
