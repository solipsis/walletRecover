package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"golang.org/x/crypto/pbkdf2"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: walletRecover [pathToEncryptedWallet] [pathToPasswordDictionary]")
		return
	}

	// read the encrypted file
	encodedPayload, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Decode base64 payload
	cipherText, err := base64.StdEncoding.DecodeString(string(encodedPayload))
	if err != nil {
		log.Fatal("Payload does not appear to be base64", err)
	}
	if len(cipherText) < aes.BlockSize {
		log.Fatal("Ciphertext block size is too short!")
	}

	// Load the dictionary
	f, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	worklist := make(chan string)
	progress := make(chan string)
	errors := make(chan error)
	done := make(chan string)

	// seed the worklist using the passwords from the dictionary file
	go func(sc *bufio.Scanner) {
		defer close(worklist)
		for sc.Scan() {
			worklist <- sc.Text()
		}
	}(sc)

	// array of methods for decrypting different blockchain.info legacy formats
	formats := []func(string, []byte) (string, error){decrypt, decryptLegacy1, decryptLegacy2}

	// 20 worker goroutines pulling from the worklist
	wg := new(sync.WaitGroup)
	for x := 0; x < 20; x++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case p := <-worklist:

					// Try decrypting in each legacy format
					for _, decryptFormat := range formats {
						dec, err := decryptFormat(p, cipherText)
						if err != nil {
							errors <- err
						}

						// try to marshal into json
						if err = attemptJSON(dec); err != nil {
							continue
						}

						// We decrypted successfully
						close(done)
						fmt.Println("Wallet decoded successfully")
						fmt.Printf("Decoded: %s\n", string(dec))
						return
					}
					progress <- fmt.Sprintf("Not json: %s", p)

				case <-done:
					return
				}
			}
		}()
	}

	// Occasionally print intermediate results
	go func() {
		x := 0
		for update := range progress {
			if x%10000 == 0 {
				fmt.Printf("Health Check: %s | %d passwords tried so far\n", update, x)
			}
			x++
		}
	}()

	// Gracefully shutdown if an error occurs
	go func() {
		for err := range errors {
			fmt.Println("Error: ", err)
			close(done)
			close(errors)
		}
	}()

	wg.Wait()
	close(progress)
}

// attempt to marshal into json
func attemptJSON(s string) error {
	var f interface{}
	return json.Unmarshal([]byte(s), &f)
}

// most recent legacy wallet format
func decrypt(plainKey string, cipherText []byte) (string, error) {

	// iv and salt are the first 16 bytes
	iv := cipherText[:aes.BlockSize]
	salt := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	cipherKey := pbkdf2.Key([]byte(plainKey), salt, 10, 32, sha1.New)
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	ecb := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(cipherText))
	ecb.CryptBlocks(decrypted, cipherText)

	return string(unpad(decrypted)), nil
}

// Only 1 iteration for pbkdf2
func decryptLegacy1(plainKey string, cipherText []byte) (string, error) {

	// iv and salt are the first 16 bytes
	iv := cipherText[:aes.BlockSize]
	salt := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	cipherKey := pbkdf2.Key([]byte(plainKey), salt, 1, 32, sha1.New)
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	ecb := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(cipherText))
	ecb.CryptBlocks(decrypted, cipherText)

	return string(unpad(decrypted)), nil
}

// No Salt or IV
func decryptLegacy2(plainKey string, cipherText []byte) (string, error) {

	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	salt := []byte{}

	cipherKey := pbkdf2.Key([]byte(plainKey), salt, 10, 32, sha1.New)
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	ecb := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(cipherText))
	ecb.CryptBlocks(decrypted, cipherText)

	return string(unpad(decrypted)), nil
}

// Remove ISO 10126 padding
func unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
