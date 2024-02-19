package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cobbinma/bromide/internal"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	Accept     = "Accept"
	Reject     = "Reject"
	Skip       = "Skip"
	listHeight = 14
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).MarginTop(1)
	subTitle          = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("12"))
	listStyle         = lipgloss.NewStyle().MarginLeft(2)
	diffStyle         = lipgloss.NewStyle().MarginLeft(2).MarginTop(1)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("12"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
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

func (r Review) Diff() string {
	old := ""
	if r.old != nil {
		old = string(r.old.contents)
	}

	return internal.Diff(old, string(r.new.contents))
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

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	reviews  []Review
	index    int
	list     list.Model
	quitting bool
	err      error
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.reviews) == 0 {
		m.quitting = true
		return m, tea.Quit
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				review := m.reviews[m.index]
				switch i {
				case Accept:
					{
						if review.old != nil {
							if err := os.Remove(review.Path(internal.Accepted)); err != nil {
								m.err = err
								return m, tea.Quit
							}
						}

						if err := os.Rename(review.Path(internal.Pending), review.Path(internal.Accepted)); err != nil {
							m.err = err
							return m, tea.Quit
						}
					}
				case Reject:
					{
						if err := os.Remove(review.Path(internal.Pending)); err != nil {
							m.err = err
							return m, tea.Quit
						}
					}
				case Skip:
					{
					}
				}
			}

			m.index++

			if m.index >= len(m.reviews) {
				m.quitting = true
				return m, tea.Quit
			}

			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if len(m.reviews) == 0 || m.quitting {
		return quitTextStyle.Render(fmt.Sprintf("Reviewed %v of %v 📸", m.index, len(m.reviews)))
	}

	if m.err != nil {
		return quitTextStyle.Render(m.err.Error())
	}

	review := m.reviews[m.index]

	out := titleStyle.Render(review.Title()) + "\n" + subTitle.Render(review.TestName())

	diff := review.Diff()

	out = out + "\n" + diffStyle.Render(diff)

	m.list.Title = ""

	return out + m.list.View()
}

func main() {
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
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	items := []list.Item{
		item(Accept),
		item(Reject),
		item(Skip),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = listStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l, reviews: reviews}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
