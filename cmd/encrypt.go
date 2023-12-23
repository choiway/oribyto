/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
    "log"
    "os/exec"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("encrypt called")
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			panic(err)
		}

		if file == "" {
			fmt.Println("file is required")
			return
		}

		recipient := cmd.Flag("recipient").Value.String()
		fmt.Println("file:", file)
		fmt.Println("recipient:", recipient)

        fmt.Println("Encrypting file...")
		err = exec.Command("gpg", "--encrypt", "--recipient", recipient, "--recipient", "waynechoi@gmail.com", file).Run()
		if err != nil {
			log.Fatal("Error occurred while encrypting file")
			log.Fatal(err)
		}

        // shred orginal file
        fmt.Println("Shredding file...")
        err = exec.Command("shred", "--remove", file).Run()
        if err != nil {
            log.Fatal("Error occurred while shredding file")
            log.Fatal(err)
        }
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().String("file", "", "Path to file to encrypt. Required")
	encryptCmd.Flags().String("recipient", "", "recipent to encrypt")
	encryptCmd.MarkFlagsRequiredTogether("file", "recipient")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
