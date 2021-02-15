package account

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/vivian-tangle/vivian-client/tools"
	"golang.org/x/crypto/ssh/terminal"
)

const newSeedCommand = "cat /dev/urandom |tr -dc A-Z9|head -c${1:-81}"
const seedFileName = "seed"

// GetSeed will try to get seed from files in seed path
func (ac *Account) GetSeed() bool {
	seedPath := ac.Config.SeedPath
	_, err := os.Stat(seedPath)
	if os.IsNotExist(err) {
		fmt.Printf("Seed path does not exist! Creating %s...\n", seedPath)
		os.MkdirAll(seedPath, os.ModePerm)

	}

	seedFilePath := filepath.Join(seedPath, seedFileName)
	_, err = os.Stat(seedFilePath)
	if os.IsNotExist(err) {
		fmt.Println("Seed file not exist! Generating new seed...")
		reader := bufio.NewReader(os.Stdin)
		var seedContent []byte
		for {
			seedContent, err = exec.Command("bash", "-c", newSeedCommand).Output()
			tools.HandleErr(err)
			fmt.Printf("Your new seed is: %s\n", seedContent)
			fmt.Print("Press Y to confirm, press the other key to generate again.")
			key, _ := reader.ReadString('\n')
			fmt.Println()
			if key == "Y\n" || key == "y\n" {
				break
			}
		}

		var bytePassword []byte
		for {
			fmt.Print("Use password to encrypt your seed:")
			bytePassword, err = terminal.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			tools.HandleErr(err)
			fmt.Print("Confirm you password:")
			confirmBytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
			tools.HandleErr(err)
			fmt.Println()
			res := bytes.Compare(bytePassword, confirmBytePassword)
			if res == 0 {
				break
			}
			fmt.Println("Please enter your password again!")
		}

		seedFile, err := tools.EncryptSeed(bytePassword, seedContent)
		tools.HandleErr(err)
		err = ioutil.WriteFile(seedFilePath, seedFile, 0777)
		tools.HandleErr(err)
		ac.Seed = string(seedContent)
		fmt.Printf("File %s saved successfully\n", seedFilePath)
	} else {
		data, err := ioutil.ReadFile(seedFilePath)
		tools.HandleErr(err)
		fmt.Print("Enter password to decrypt your seed:")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		seedContent, err := tools.DecryptSeed(bytePassword, data)
		ac.Seed = string(seedContent)
		fmt.Println("Get seed successfully.")
	}

	return true
}
