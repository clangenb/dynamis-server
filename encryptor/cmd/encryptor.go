package main

import (
	"encryptor/crypto"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Encryptor",
		Usage: "Encrypt and decrypt audio files using AES-256-GCM",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "in",
				Usage: "Input file path",
			},
			&cli.StringFlag{
				Name:  "out",
				Usage: "Output file path",
			},
		},
		Action: func(c *cli.Context) error {
			input := c.String("in")
			output := c.String("out")
			if input == "" || output == "" {
				return fmt.Errorf("both input and output paths are required")
			}

			masterKey, err := crypto.LoadMasterKey()
			if err != nil {
				return err
			}

			if err := crypto.EncryptFile(input, output, masterKey); err != nil {
				return fmt.Errorf("encryption failed: %v", err)
			}
			fmt.Println("File encrypted successfully!")
			return nil
		},
	}

	// Running the application
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
