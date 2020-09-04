package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"goCkup/utils"
	"google.golang.org/api/drive/v3"
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

var decryptFlag bool

func init() {
	CmdDownload.Flags().StringVar(&decryptionKeyFlag, "key", "", "key for encryption!!!! Should be 32 Byte!!!")
	CmdDownload.Flags().StringVar(&decryptionKeyFileFlag, "keyFile", "", "key for encryption!!!! Should be 32 Byte!!!")
	CmdDownload.Flags().BoolVar(&decryptFlag, "decrypt", false, "Flag that indicates whether the downloaded file should be encrypted.")
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
			var gFile *drive.File
			for _, file := range filesList.Files {
				if index == indexToDownload {
					gFile = file
					break
				}
				index++
			}
			if gFile == nil {
				log.Fatal("unknown error")
			}
			resp, err := service.Files.Get(gFile.Id).Download()
			if err != nil {
				log.Fatalf("Error while getting file. %v", err)
			}
			fileNameToDownload := args[0]
			if decryptFlag {
				fileNameToDownload += ".enc"
			}
			fileToDownload, _ := os.OpenFile(fileNameToDownload, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
			defer fileToDownload.Close()
			if _, err := io.Copy(fileToDownload, resp.Body); err != nil {
				log.Fatalf("Error during encription. %v", err)
			}
			if decryptFlag {
				decrypt([]string{fileNameToDownload, args[0]}, decryptionKeyFlag, decryptionKeyFileFlag)
				defer os.Remove(fileNameToDownload)
			}
			break
		}
	}
}
