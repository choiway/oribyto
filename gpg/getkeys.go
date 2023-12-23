package gpg

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type GpgKey struct {
	UID         string
	Email       string
	Fingerprint string
}

// GetSecretKeys returns a list of secret keys
func GetSecretKeys() []GpgKey {
	b, err := exec.Command("gpg", "--list-secret-keys", "--with-colons").Output()
	if err != nil {
		fmt.Println("The following error ocurred while trying to get list secret keys: ", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(b))
	secretGpgKeys := []GpgKey{}
	var newGpgKey GpgKey
	pk := ""

	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ":")

		k := words[0]

		switch k {
		case "sec":
			// Initialize new GpgKey
			newGpgKey = GpgKey{}
		case "uid":
			newGpgKey.UID = words[9]
			newGpgKey.Email = strings.TrimSuffix(strings.SplitAfter(words[9], "<")[1], ">")
		case "fpr":
			if pk == "sec" {
				newGpgKey.Fingerprint = words[9]
			}
		case "ssb":
			// append to secretGpgKeys
			secretGpgKeys = append(secretGpgKeys, newGpgKey)
		}

		pk = k
	}


	if err := scanner.Err(); err != nil {
		fmt.Println("Error ocurred while scanning keys", err)
		os.Exit(1)
	}

    return secretGpgKeys
}
