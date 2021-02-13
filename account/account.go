package account

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/vivian-tangle/vivian-client/config"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/ssh/terminal"
)

const newSeedCommand = "cat /dev/urandom |tr -dc A-Z9|head -c${1:-81}"
const seedFileName = "seed"

// Account is the structure for storing the account info
type Account struct {
	Seed string
}

// GetSeed will try to get seed from files in seed path
func (ac *Account) GetSeed(c config.Config) bool {
	_, err := os.Stat(c.SeedPath)
	if os.IsNotExist(err) {
		fmt.Printf("Seed path does not exist! Creating %s...\n", c.SeedPath)
		os.MkdirAll(c.SeedPath, os.ModePerm)

	}

	seedFilePath := filepath.Join(c.SeedPath, seedFileName)
	_, err = os.Stat(seedFilePath)
	if os.IsNotExist(err) {
		fmt.Println("Seed file not exist! Generating new seed...")
		reader := bufio.NewReader(os.Stdin)
		var seedContent []byte
		for {
			seedContent, err = exec.Command("bash", "-c", newSeedCommand).Output()
			handleErr(err)
			fmt.Printf("Your new seed is: %s\n", seedContent)
			fmt.Println("Press Y to confirm, press the other key to generate again.")
			key, _ := reader.ReadString('\n')
			if key == "Y\n" || key == "y\n" {
				break
			}
		}

		var bytePassword []byte
		for {
			fmt.Println("Use password to encrypt your seed:")
			bytePassword, err = terminal.ReadPassword(int(syscall.Stdin))
			handleErr(err)
			fmt.Println("Confirm you password:")
			confirmBytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
			handleErr(err)
			res := bytes.Compare(bytePassword, confirmBytePassword)
			if res == 0 {
				break
			}
			fmt.Println("Please enter your password again!")
		}

		seedFile, err := EncryptSeed(bytePassword, seedContent)
		handleErr(err)
		err = ioutil.WriteFile(seedFilePath, seedFile, 0777)
		handleErr(err)
		ac.Seed = string(seedContent)
		fmt.Printf("File %s saved successfully", seedFilePath)
	} else {
		data, err := ioutil.ReadFile(seedFilePath)
		handleErr(err)
		fmt.Println("Enter password to decrypt your seed:")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		seedContent, err := DecryptSeed(bytePassword, data)
		ac.Seed = string(seedContent)
		fmt.Println("Get seed successfully.")
	}

	return true
}

// EncryptSeed encrypts seed with password
func EncryptSeed(key, data []byte) ([]byte, error) {
	key, salt, err := DeriveKey(key, nil)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	ciphertext = append(ciphertext, salt...)

	return ciphertext, nil
}

// DecryptSeed decrypts the file for retriving the seed
func DecryptSeed(key, data []byte) ([]byte, error) {
	salt, data := data[len(data)-32:], data[:len(data)-32]

	key, _, err := DeriveKey(key, salt)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// DeriveKey converts the password to suitable keys for AES algorithm
func DeriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
