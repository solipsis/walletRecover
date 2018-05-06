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

	//
	worklist := make(chan string)
	results := make(chan string)
	go func(sc *bufio.Scanner) {
		defer close(worklist)
		for sc.Scan() {
			worklist <- sc.Text()
		}
	}(sc)

	wg := new(sync.WaitGroup)

	for x := 0; x < 20; x++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range worklist {
				//cipherKey := pbkdf2.Key([]byte(p), []byte{}, 10, 32, sha1.New)
				dec, err := decrypt(p, cipherText)
				if err != nil {
					log.Fatal(err)
				}
				var e interface{}
				err = json.Unmarshal([]byte(dec), &e)
				if err != nil {
					results <- fmt.Sprintf("Not json: %s", p)
					continue
				}
				fmt.Printf("Decoded: %s\n", string(dec))
				fmt.Println("Done")
				fmt.Println("************")
				fmt.Println("Success: " + p)
				results <- fmt.Sprintf("Success: %s", p)
				panic("Yay")
			}
		}()
	}

	go func() {
		x := 0
		for r := range results {
			if x%10000 == 0 {
				fmt.Printf("%s %d\n", r, x)
			}
			x++
		}
	}()

	wg.Wait()
}

/*
func main() {
	js := `{"Id":4, "Potato":"Tomato"}`
	cipherKey := pbkdf2.Key([]byte("stack901"), []byte{}, 10, 32, sha1.New)
	enc, err := encrypt([]byte(cipherKey), js)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encrypted: %s\n", string(enc))
	dec, err := decrypt(cipherKey, enc)
	if err != nil {
		log.Fatal(err)
	}
	var e interface{}
	err = json.Unmarshal([]byte(dec), &e)
	if err != nil {
		fmt.Println("NOT JSON", err)
	}
	fmt.Printf("Decoded: %s\n", string(dec))
	fmt.Println("Done")

}
*/

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

/*
// no salt
func decrypt(plainKey string, encrypted string) (string, error) {

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("Ciphertext block size is too short!")
	}

	// iv and salt are the first 16 bytes
	iv := cipherText[:aes.BlockSize]
	salt := []byte{}
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
*/

func unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

/*
func encrypt(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	//stream := cipher.NewCFBEncrypter(block, iv)
	//stream := cipher.NewCBCEncrypter(block, iv)
	//stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmess = base64.StdEncoding.EncodeToString(cipherText)
	return
}
*/
