package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "file-search",
		Usage: "Searches for files with a specific extension",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				fmt.Println("Please provide a directory path.")
				return nil
			}

			dir := c.Args().Get(0)

			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() && filepath.Ext(path) == ".new" {
					neww, err := os.ReadFile(path)
					if err != nil {
						if !os.IsNotExist(err) {
							return err
						}
					}

					accepted := strings.TrimSuffix(path, ".new") + ".accepted"
					existing, err := os.ReadFile(accepted)
					if err != nil {
						if !os.IsNotExist(err) {
							return err
						}
					}

					dmp := diffmatchpatch.New()
					diffs := dmp.DiffMain(string(existing), string(neww), true)
					fmt.Println(dmp.DiffPrettyText(diffs))

					prompt := promptui.Select{
						Label: "snapshot review",
						Items: []string{
							"accept", "reject", "skip",
						},
					}

					_, result, err := prompt.Run()
					if err != nil {
						fmt.Printf("Prompt failed %v\n", err)
						return nil
					}

					switch result {
					case "accept":
						{
							if string(existing) != "" {
								if err := os.Remove(accepted); err != nil {
									return err
								}
							}

							if err := os.Rename(path, accepted); err != nil {
								return err
							}
						}
					case "reject":
						{
							if err := os.Remove(path); err != nil {
								return err
							}
						}
					case "skip":
						{
						}
					}
				}
				return nil
			})
			if err != nil {
				fmt.Println("Error:", err)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
