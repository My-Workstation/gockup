package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
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
	Short: "_todo_",
	Long:  `_todo_`,
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
	bytes, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials.json. %v", err)
	}
	config, err := google.ConfigFromJSON(bytes, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse credentials.json. %v", err)
	}
	token := getToken(config)
	service, err := drive.NewService(context.Background(), option.WithTokenSource(config.TokenSource(context.Background(), token)))
	if err != nil {
		log.Fatalf("Unable to create service. %v", err)
	}

	toUploadFileName := args[0]
	if encryptFlag {
		toUploadFileName = encrypt(args[:1], encryptionKeyFlag, encryptionKeyFileFlag)
		if !keepLocalFlag {
			defer os.Remove(toUploadFileName)
		}
	}
	toUpload, _ := os.Open(toUploadFileName)
	defer toUpload.Close()
	_, err = service.Files.Create(&drive.File{Name: args[0]}).Media(toUpload).Do()
	if err != nil {
		log.Printf("Error during upload. %v", err)
	}
}

func getToken(config *oauth2.Config) *oauth2.Token {
	tokenFileName := "token.json"
	tokenFile, err := os.Open(tokenFileName)
	token := &oauth2.Token{}
	if err != nil {
		// there is no token
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)
		var authCode string
		if _, err := fmt.Scan(&authCode); err != nil {
			log.Fatalf("Unable to read code. %v", err)
		}
		token, err = config.Exchange(context.TODO(), authCode)
		if err != nil {
			log.Fatalf("unable to get token. %v", err)
		}
		// create new file
		tokenFile, err = os.OpenFile(tokenFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatalf("unable to create token file. %v", err)
		}
		err = json.NewEncoder(tokenFile).Encode(token)
		if err != nil {
			log.Fatalf("unable to save token. %v", err)
		}
	} else {
		// there is a token
		err = json.NewDecoder(tokenFile).Decode(token)
		if err != nil {
			log.Fatalf("Unable to parse token. %v", err)
		}
	}
	tokenFile.Close()
	return token
}

func main___() {
	bytes, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials.json. %v", err)
	}
	config, err := google.ConfigFromJSON(bytes, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse credentials.json. %v", err)
	}
	token := getToken(config)
	service, err := drive.NewService(context.Background(), option.WithTokenSource(config.TokenSource(context.Background(), token)))
	if err != nil {
		log.Fatalf("Unable to create service. %v", err)
	}

	//toUpload, _ := os.Open("1")
	//defer toUpload.Close()
	//service.Files.Insert(&drive.File{Title: "1"}).Media(toUpload).Do()

	filesList, err := service.Files.List().Do()
	if err != nil {
		log.Fatalf("Unable to read from google cloud. %v", err)
	}
	for _, i := range filesList.Files {
		fmt.Printf("%s (%s) deleted(%v)\n", i.OriginalFilename, i.Id, i.TrashedTime)
	}
}
