package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"goCkup/utils"
	"io"
	"log"
	"os"
	"strconv"
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
		download(args)
	},
}

func download(args []string) {
	service := utils.GetService()
	filesList, err := service.Files.List().Do()
	if err != nil {
		log.Fatalf("Unable to read from google cloud. %v", err)
	}
	index := 1
	fmt.Println("You have some files:")
	fmt.Println("# \t\t name")
	for _, file := range filesList.Files {
		fmt.Printf("%d \t\t %s\n", index, file.Name)
		index++
	}
	cliReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Which file to download? Chose from 1 to %d\n", index-1)
		cliInput, _ := cliReader.ReadString('\n')
		// convert CRLF to LF
		cliInput = strings.Replace(cliInput, "\n", "", -1)
		indexToDownload, err := strconv.Atoi(cliInput)
		if err != nil {
			fmt.Print("It is not number")
		} else if indexToDownload < 1 || indexToDownload >= index {
			fmt.Println("It is out of the range.")
		} else {
			index = 1
			for _, file := range filesList.Files {
				if index == indexToDownload {
					resp, err := service.Files.Get(file.Id).Download()
					if err != nil {
						log.Fatalf("Error while getting file. %v", err)
					}
					fileToDownload, _ := os.OpenFile("download", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
					defer fileToDownload.Close()
					if _, err := io.Copy(fileToDownload, resp.Body); err != nil {
						log.Fatalf("Error during encription. %v", err)
					}
					break
				}
				index++
			}
			break
		}
	}
}
