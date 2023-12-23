/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	// "database/sql"
	_ "github.com/mattn/go-sqlite3"

	"bufio"
	"bytes"
	"encoding/gob"
	"github.com/spf13/cobra"
	"os/exec"
	"strings"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		// output, err := exec.Command("gpg", "--list-secret-keys").Output()
		//
		// if err != nil {
		//     fmt.Println("Error: ", err)
		//     os.Exit(1)
		// }
		//
		// fmt.Println(string(output))
		GetSecretKeys()
		settings := Settings{DefaultKeyEmail: "waynechoi@gmail.com"}
		path := "settings.gob"
		// write
		WriteSettings(settings, path)

		// read
		s := ReadSettings(path)
		fmt.Printf("%q", s)
	},
}

type Settings struct {
	DefaultKeyEmail string
}

// WriteSettings writes settings to a file
func WriteSettings(settings Settings, filename string) {
	out, err := os.Create(filename)
	if err != nil {
		fmt.Printf("File write error: %v\n", err)
		os.Exit(1)
	}
	w := bufio.NewWriter(out)
	enc := gob.NewEncoder(w)
	enc.Encode(settings)
	w.Flush()
	out.Close()
}

type GpgKey struct {
	UID         string
	Email       string
	Fingerprint string
}

func GetSecretKeys() {
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

	// print secretGpgKeys
	for _, gpgKey := range secretGpgKeys {
		fmt.Println(gpgKey)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error ocurred while scanning keys", err)
		os.Exit(1)
	}
}

// ReadSettings reads settings from a file
func ReadSettings(filename string) Settings {
	var b Settings

	in, err := os.Open(filename)
	if err != nil {
		fmt.Printf("File read error: %v\n", err)
		os.Exit(1)
	}
	r := bufio.NewReader(in)
	dec := gob.NewDecoder(r)
	dec.Decode(&b)
	return b
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// create new sqlite db
	//  os.Create("oribyte.db")
}
