package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cobbinma/bromide/internal"
	"github.com/manifoldco/promptui"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/urfave/cli/v2"
)

type choice string

const (
	accept choice = "accept"
	reject choice = "reject"
	skip   choice = "skip"
)

type Snapshot struct {
	contents []byte
}

type Review struct {
	path string
	old  *Snapshot
	new  Snapshot
}

func (r Review) Path(status internal.ReviewState) string {
	return r.path + status.Extension()
}

func main() {
	app := &cli.App{
		Name:  "bromide",
		Usage: "review tool for bromide snapshot testing",
		Commands: []*cli.Command{
			{
				Name:    "review",
				Aliases: []string{"r"},
				Usage:   "review snapshots",
				Action: func(c *cli.Context) error {
					reviews := []Review{}
					if err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}

						if !info.IsDir() && filepath.Ext(path) == internal.Accepted.Extension() {
							neww, err := os.ReadFile(path)
							if err != nil {
								return err
							}

							var old *Snapshot
							accepted := strings.TrimSuffix(path, internal.Accepted.Extension()) + internal.Pending.Extension()
							existing, err := os.ReadFile(accepted)
							if err != nil {
								if !os.IsNotExist(err) {
									return err
								}
							} else {
								old = &Snapshot{
									contents: existing,
								}
							}

							reviews = append(reviews, Review{
								path: strings.TrimSuffix(path, internal.Pending.Extension()),
								old:  old,
								new:  Snapshot{contents: neww},
							})
						}

						return nil
					}); err != nil {
						fmt.Println("Error:", err)
					}

					if len(reviews) == 0 {
						return nil
					}

					for i, review := range reviews {
						accepted := string(review.new.contents)
						existing := ""
						if review.old != nil {
							existing = string(review.old.contents)
						}

						dmp := diffmatchpatch.New()
						diffs := dmp.DiffMain(existing, accepted, true)

						fmt.Printf("reviewing %v of %v\n", i+1, len(reviews))
						fmt.Println(dmp.DiffPrettyText(diffs))

						prompt := promptui.Select{
							Label: "snapshot review",
							Items: []string{
								string(accept), string(reject), string(skip),
							},
						}

						_, result, err := prompt.Run()
						if err != nil {
							fmt.Printf("prompt failed: %v\n", err)
							return nil
						}

						switch choice(result) {
						case accept:
							{
								if string(existing) != "" {
									if err := os.Remove(review.Path(internal.Accepted)); err != nil {
										return err
									}
								}

								if err := os.Rename(review.Path(internal.Accepted), review.Path(internal.Pending)); err != nil {
									return err
								}
							}
						case reject:
							{
								if err := os.Remove(review.Path(internal.Pending)); err != nil {
									return err
								}
							}
						case skip:
							{
							}
						}
					}

					fmt.Printf("reviewed %v snapshot(s) 📸\n", len(reviews))

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("error:", err)
	}
}
