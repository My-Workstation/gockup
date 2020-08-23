package utils

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func MakeRandom(length int) []byte {
	rnd := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, rnd); err != nil {
		panic(err)
	}
	return rnd
}

func ReadKeyFromFile(filePath string) []byte {
	// try to open file
	encryptionKeyFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Can not open key file. %v", err)
	}
	defer encryptionKeyFile.Close()
	tmp := make([]byte, 33)
	count, err := encryptionKeyFile.Read(tmp)
	if err != nil {
		log.Fatal("Can not read key file. %v", err)
	}
	if count < 32 {
		log.Fatal("Key should be 32 bytes")
	} else if count > 32 {
		log.Print("Key should be 32 bytes. Only the first 32 bytes will be used.")
	}
	return tmp[:32]
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

func GetService() *drive.Service {
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
	return service
}
