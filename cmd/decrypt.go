package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/spf13/cobra"
	"goCkup/utils"
	"io"
	"log"
	"os"
)

// goCkup decrypt --key='32bitkey _fileName_
// goCkup decrypt --key='32bitkey _fileName_ _fileNameOutput_
// goCkup decrypt --keyFile='32bitkey' _fileName_ _fileNameOutput_
var CmdDecrypt = &cobra.Command{
	Use:   "decrypt _filename_",
	Short: "Take _filename_ and decrypt it",
	Long:  `Take _filename_ and decrypt it`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		decrypt(args, decryptionKeyFlag, decryptionKeyFileFlag)
	},
}

var decryptionKeyFlag string
var decryptionKeyFileFlag string

func init() {
	CmdDecrypt.Flags().StringVar(&decryptionKeyFlag, "key", "", "key for encryption!!!! Should be 32 Byte!!!")
	CmdDecrypt.Flags().StringVar(&decryptionKeyFileFlag, "keyFile", "", "key for encryption!!!! Should be 32 Byte!!!")
}

func decrypt(args []string, decryptionKeyFlag string, decryptionKeyFileFlag string) {
	if len(decryptionKeyFlag) == 0 && len(decryptionKeyFileFlag) == 0 {
		log.Fatal("Neither the key nor the key file have been transferred.")
	}
	if len(decryptionKeyFlag) != 0 && len(decryptionKeyFileFlag) != 0 {
		log.Print("The keyFile does not matter. The key has priority over the keyFile.")
	}
	encryptionKey := make([]byte, 32)
	if len(decryptionKeyFlag) != 0 {
		if len(decryptionKeyFlag) < 32 {
			log.Fatal("Key should be 32 bytes")
		}
		if len(decryptionKeyFlag) > 32 {
			log.Print("Key should be 32 bytes. Only the first 32 bytes will be used.")
		}
		encryptionKey = []byte(decryptionKeyFlag[:32])
	} else {
		encryptionKey = utils.ReadKeyFromFile(decryptionKeyFileFlag)
	}

	fileToDecode, _ := os.Open(args[0])
	defer fileToDecode.Close()
	decodedFileName := "./decrypted"
	if len(args) > 1 && len(args[1]) != 0 {
		decodedFileName = args[1]
	}
	decodedFile, _ := os.OpenFile(decodedFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer decodedFile.Close()
	block, _ := aes.NewCipher(encryptionKey)
	iv := make([]byte, aes.BlockSize)
	count, _ := fileToDecode.Read(iv)
	if count != len(iv) {
		log.Fatal("no iv found")
	}
	stream := cipher.NewOFB(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: fileToDecode}
	if _, err := io.Copy(decodedFile, reader); err != nil {
		log.Fatalf("Error during copy. %v", err)
	}
}
