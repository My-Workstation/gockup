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

// goCkup encrypt _fileName_ _fileNameOutput_ <- random generated key. Should be logged out at the end.
// goCkup encrypt --key='32bitkey' _fileNameInput_ _fileNameOutput_
// goCkup encrypt --keyFile='32bitkey' _fileName_ _fileNameOutput_
// todo suggest for bash
var CmdEncrypt = &cobra.Command{
	Use:   "encrypt _filename_",
	Short: "Take _filename_ ane encrypt it",
	Long:  `Take _filename_ ane encrypt it`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		encrypt(args, encryptionKeyFlag, encryptionKeyFileFlag)
	},
}
var encryptionKeyFlag string
var encryptionKeyFileFlag string

func init() {
	CmdEncrypt.Flags().StringVar(&encryptionKeyFlag, "key", "", "key for encryption")
	CmdEncrypt.Flags().StringVar(&encryptionKeyFileFlag, "keyFile", "", "key for encryption.!!!! Should be 32 Byte!!!")
}

func encrypt(args []string, encryptionKeyFlag string, encryptionKeyFileFlag string) {
	encryptionKey := make([]byte, 32)
	if len(encryptionKeyFlag) == 0 {
		if len(encryptionKeyFileFlag) != 0 {
			// try to open file
			encryptionKeyFile, err := os.Open(encryptionKeyFileFlag)
			if err != nil {
				log.Fatalf("Can not open key file. %v", err)
			}
			defer encryptionKeyFile.Close()
			tmp := make([]byte, 33)
			count, err := encryptionKeyFile.Read(tmp)
			if err != nil {
				log.Fatalf("Can not read key file. %v", err)
			}
			if count < 32 {
				log.Fatalf("Key should be 32 bytes")
			} else if count > 32 {
				log.Print("Key should be 32 bytes. Only the first 32 bytes will be used.")
			}
			encryptionKey = tmp[:32]
		}
	} else {
		if len(encryptionKeyFileFlag) != 0 {
			log.Print("The keyFile does not matter. The key has priority over the keyFile.")
		}
		if len(encryptionKeyFlag) < 32 {
			log.Fatalf("Key should be 32 bytes")
		} else if len(encryptionKeyFlag) > 32 {
			log.Print("Key should be 32 bytes. Only the first 32 bytes will be used.")
		}
		encryptionKey = []byte(encryptionKeyFlag[:32])
	}
	if len(encryptionKey) == 0 {
		encryptionKey = utils.MakeRandom(32)
		log.Print("The key wasn't handed over. It will be generated automatically. Key:", encryptionKey)
	}

	fileToEncode, _ := os.Open(args[0])
	defer fileToEncode.Close()
	encodedFileName := "./encoded"
	if len(args) > 1 && len(args[1]) != 0 {
		encodedFileName = args[1]
	}
	encodedFile, err := os.OpenFile(encodedFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Can not open encrypted file. %v", err)
	}
	defer encodedFile.Close()

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		log.Fatalf("Unable to create Chiper. %v", err)
	}
	iv := utils.MakeRandom(aes.BlockSize)
	stream := cipher.NewOFB(block, iv[:])
	_, err = encodedFile.Write(iv)
	if err != nil {
		log.Fatalf("Can not write encrypted file. %v", err)
	}
	writer := &cipher.StreamWriter{S: stream, W: encodedFile}
	if _, err := io.Copy(writer, fileToEncode); err != nil {
		log.Fatalf("Error during encription. %v", err)
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
