package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/spf13/cobra"
	"goCkup/utils"
	"io"
	"log"
	"os"
	"strings"
)

// goCkup encrypt _fileName_
var CmdEncrypt = &cobra.Command{
	Use:   "encrypt _todo_",
	Short: "_todo_",
	Long:  `_todo_`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}

func main_() {

	key := []byte("AES256Key-32Characters1234567890")

	fileToEncode, _ := os.Open("2")
	defer fileToEncode.Close()
	encodedFile, _ := os.OpenFile("./encoded", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer encodedFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("Unable to create Chiper. %v", err)
	}
	iv := utils.MakeRandom(aes.BlockSize)
	stream := cipher.NewOFB(block, iv[:])
	encodedFile.Write(iv)
	writer := &cipher.StreamWriter{S: stream, W: encodedFile}
	if _, err := io.Copy(writer, fileToEncode); err != nil {
		panic(err)
	}
}

// block cipher
//key := []byte("AES256Key-32Characters1234567890")
//block, err := aes.NewCipher(key)
//if err != nil {
//	log.Fatalf("Unable to create Chiper. %v", err)
//}
//aesgcm, err := cipher.NewGCM(block)
//if err != nil {
//	log.Fatalf("Unable to create GCM. %v", err)
//}
//nonce := make([]byte, aesgcm.NonceSize())
//if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
//	log.Fatalf("Unable to create nonce. %v", err)
//}
//
//fileToEncode, _ := os.Open("1")
//defer fileToEncode.Close()
//chunk := make([]byte, 32)
//encodedFile, _ := os.OpenFile("./encoded", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
//defer encodedFile.Close()
//for {
//	readCount, err := fileToEncode.Read(chunk)
//	if err == io.EOF {
//		break
//	}
//	encodedFile.Write(aesgcm.Seal(nil, nonce, chunk[:readCount], nil))
//}
//
//fileToDecode, _ := os.Open("encoded")
//defer fileToDecode.Close()
//decodedFile, _ := os.OpenFile("./decoded", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
//defer decodedFile.Close()
//chunk = make([]byte, 48)
//for {
//	readCount, err := fileToDecode.Read(chunk)
//	if err == io.EOF {
//		break
//	}
//	dec, _ := aesgcm.Open(nil, nonce, chunk[:readCount], nil)
//	decodedFile.Write(dec)
//}
//return
