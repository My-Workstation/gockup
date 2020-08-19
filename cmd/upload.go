package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// goCkup upload --encrypt  _fileName_
// goCkup upload --encrypt --provider=google _fileName_
// goCkup upload --encrypt --saveLocal --provider=google _fileName_
// set encryption key
// generate new encryption key
// SUPER FEATURE: if it is a directory make tar
var CmdUpload = &cobra.Command{
	Use:   "upload _todo_",
	Short: "_todo_",
	Long:  `_todo_`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}

//cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

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
	for _, i := range filesList.Items {
		fmt.Printf("%s (%s) deleted(%v)\n", i.OriginalFilename, i.Id, i.TrashedDate)
	}
}
