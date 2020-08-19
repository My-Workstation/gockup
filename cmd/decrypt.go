package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

// goCkup decrypt _fileName_
var CmdDecrypt = &cobra.Command{
	Use:   "encrypt _todo_",
	Short: "_todo_",
	Long:  `_todo_`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}

func main__() {
	key := []byte("AES256Key-32Characters1234567890")
	fileToDecode, _ := os.Open("encoded")
	defer fileToDecode.Close()
	decodedFile, _ := os.OpenFile("./decoded", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer decodedFile.Close()
	block, _ := aes.NewCipher(key)
	block, _ = aes.NewCipher(key)
	iv := make([]byte, aes.BlockSize)
	count, _ := fileToDecode.Read(iv)
	if count != len(iv) {
		panic("no iv")
	}
	stream := cipher.NewOFB(block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: fileToDecode}
	if _, err := io.Copy(decodedFile, reader); err != nil {
		panic(err)
	}
}
