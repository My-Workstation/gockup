package main

import (
	"github.com/spf13/cobra"
	"goCkup/cmd"
)

func main() {
	var cmdGoCkup = &cobra.Command{Use: "goCkup"}
	cmdGoCkup.AddCommand(cmd.CmdUpload)
	cmdGoCkup.AddCommand(cmd.CmdDownload)
	cmdGoCkup.AddCommand(cmd.CmdEncrypt)
	cmdGoCkup.AddCommand(cmd.CmdDecrypt)
	cmdGoCkup.Execute()
}
