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

	// 20 worker goroutines pulling from the worklist
	wg := new(sync.WaitGroup)
	for x := 0; x < 20; x++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case p := <-worklist:
					dec, err := decrypt(p, cipherText)
					if err != nil {
						errors <- err
					}

					var e interface{}
					err = json.Unmarshal([]byte(dec), &e)
					if err != nil {
						progress <- fmt.Sprintf("Not json: %s", p)
						continue
					}
					close(done)
					fmt.Printf("Decoded: %s\n", string(dec))
					fmt.Println("Done")
					fmt.Println("************")
					fmt.Println("Success: " + p)

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

	iv := []byte{}
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
