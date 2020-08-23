package cmd

import (
	"github.com/spf13/cobra"
	"goCkup/utils"
	"google.golang.org/api/drive/v3"
	"log"
	"os"
	"time"
)

// goCkup upload _fileName_
// goCkup upload --encrypt  _fileName_
// goCkup upload --encrypt --provider=google _fileName_
// goCkup upload --encrypt --saveLocal --provider=google _fileName_
// set encryption key
// generate new encryption key
// SUPER FEATURE: if it is a directory make tar
var CmdUpload = &cobra.Command{
	Use:   "upload _fileName_",
	Short: "Take _filename_ and upload it",
	Long:  `Take _filename_ and upload it`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		upload(args)
	},
}

var encryptFlag bool
var keepLocalFlag bool

func init() {
	CmdUpload.Flags().BoolVar(&encryptFlag, "encrypt", false, "Flag that indicates whether or not to encrypt when downloading.")
	CmdUpload.Flags().StringVar(&encryptionKeyFlag, "key", "", "key for encryption")
	CmdUpload.Flags().StringVar(&encryptionKeyFileFlag, "keyFile", "", "key for encryption!!!! Should be 32 Byte!!!")
	CmdUpload.Flags().BoolVar(&keepLocalFlag, "keepLocal", false, "Flag that indicates whether the encrypted file should be kept local.")
}

func upload(args []string) {
	service := utils.GetService()

	toUploadFileName := args[0]
	if encryptFlag {
		toUploadFileName = encrypt(args[:1], encryptionKeyFlag, encryptionKeyFileFlag)
		if !keepLocalFlag {
			defer os.Remove(toUploadFileName)
		}
	}
	toUpload, _ := os.Open(toUploadFileName)
	defer toUpload.Close()
	_, err := service.Files.Create(&drive.File{Name: args[0], ModifiedTime: time.Now().Format(time.RFC3339), CreatedTime: time.Now().Format(time.RFC3339)}).Media(toUpload).Do()
	if err != nil {
		log.Printf("Error during upload. %v", err)
	}
}
