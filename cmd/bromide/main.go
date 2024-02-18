package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cobbinma/bromide/internal"
	"github.com/manifoldco/promptui"
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

func (r Review) TestName() string {
	splits := strings.SplitAfter(r.path, "/")
	if len(splits) == 0 {
		return ""
	}

	return splits[len(splits)-1]
}

func (r Review) Path(status internal.ReviewState) string {
	return r.path + status.Extension()
}

func (r Review) Title() string {
	title := ""
	switch {
	case r.old == nil:
		title = "New Snapshot"
	default:
		title = "Mismatched Snapshot"
	}

	return title
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

						if !info.IsDir() && filepath.Ext(path) == internal.Pending.Extension() {
							neww, err := os.ReadFile(path)
							if err != nil {
								return err
							}

							var old *Snapshot
							accepted := strings.TrimSuffix(path, internal.Pending.Extension()) + internal.Accepted.Extension()
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
						fmt.Println("error: ", err)
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

						fmt.Printf("reviewing %v of %v\n\n", i+1, len(reviews))
						fmt.Printf("%s\n", review.TestName())
						fmt.Println(internal.Diff(existing, accepted))

						prompt := promptui.Select{
							Label: review.Title(),
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

								if err := os.Rename(review.Path(internal.Pending), review.Path(internal.Accepted)); err != nil {
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

					fmt.Printf("\nreviewed %v snapshot(s) 📸\n", len(reviews))

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("error: ", err)
	}
}
